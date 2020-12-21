package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/discord"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/eveisesi/zrule/internal/rest"
	"github.com/eveisesi/zrule/internal/slack"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run() error
	SendTestMessage(ctx context.Context, action *zrule.Action, message string) error
}

type service struct {
	redis    *redis.Client
	logger   *logrus.Logger
	newrelic *newrelic.Application
	client   *http.Client
	policy   policy.Service
	action   action.Service
}

func NewService(redis *redis.Client, logger *logrus.Logger, newrelic *newrelic.Application, client *http.Client, policy policy.Service, action action.Service) Service {

	return &service{
		redis:    redis,
		logger:   logger,
		newrelic: newrelic,
		client:   client,
		policy:   policy,
		action:   action,
	}

}

func (s *service) Run() error {

	for {
		var ctx = context.Background()
		txn := s.newrelic.StartTransaction("dispatch queue check")
		ctx = newrelic.NewContext(ctx, txn)

		stop, err := s.redis.Get(ctx, zrule.QUEUE_STOP).Int64()
		if err != nil {
			s.logger.WithError(err).Error("stop flag is missing. attempting to create with default value of 0")
			_, err := s.redis.Set(ctx, zrule.QUEUE_STOP, 0, 0).Result()
			if err != nil {
				txn.NoticeError(err)
				s.logger.WithError(err).Fatal("error encountered attempting to create stop flag with default value")
			}
			txn.End()
			continue
		}

		if stop == 1 {
			s.logger.Info("stop signal set, sleeping for 5 seconds")
			time.Sleep(time.Second * 5)
			txn.End()
			continue
		}

		count, err := s.redis.ZCount(ctx, zrule.QUEUES_KILLMAIL_MATCHED, "-inf", "+inf").Result()
		if err != nil {
			txn.NoticeError(err)
			s.logger.WithError(err).Error("unable to determine count of message queue")
			txn.End()
			time.Sleep(time.Second * 2)
			continue
		}

		if count == 0 {
			txn.End()
			time.Sleep(time.Second * 2)
			continue
		}

		results, err := s.redis.ZPopMax(ctx, zrule.QUEUES_KILLMAIL_MATCHED, 1).Result()
		if err != nil {
			txn.NoticeError(err)
			s.logger.WithError(err).Fatal("unable to retrieve hashes from queue")
			txn.End()
			continue
		}

		for _, result := range results {
			s.logger.Info("handling message")
			data := result.Member.(string)
			s.handleMessage(ctx, []byte(data), 0)
		}

		txn.End()

	}

}

func (s *service) handleMessage(ctx context.Context, data []byte, sleep int) {

	var message = new(zrule.Dispatchable)
	err := json.Unmarshal(data, message)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.logger.WithError(err).WithField("data", string(data)).Error("failed to unmarsahl data onto dispatchable struct")
		return
	}

	policy, err := s.policy.Policy(ctx, message.PolicyID)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.logger.WithError(err).WithField("policyID", message.PolicyID).Error("failed to look up policy")
		return
	}

	newrelic.FromContext(ctx).AddAttribute("policyID", message.PolicyID.Hex())

	for _, actionID := range policy.Actions {
		entry := s.logger.WithField("policyID", message.PolicyID.Hex()).WithField("actionID", actionID)
		action, err := s.action.Action(ctx, actionID)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			entry.WithError(err).Error("failed to lookup action")
			continue
		}

		entry = entry.WithField("platform", action.Platform.String())

		platform, err := s.serviceForPlatform(action)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			entry.WithError(err).Error("unable to determine platform to use")
			continue
		}

		err = platform.Send(ctx, policy, message.ID, message.Hash)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			entry.WithError(err).Error("failed to send message to platform")
			continue
		}
		time.Sleep(time.Second)
	}

}

func (s *service) SendTestMessage(ctx context.Context, action *zrule.Action, message string) error {

	platform, err := s.serviceForPlatform(action)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.logger.WithError(err).Error("failed to lookup action")
		return err
	}

	err = platform.SendTest(ctx, message)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.logger.WithError(err).Error("failed to send message")
		return err
	}

	return nil

}

func (s *service) serviceForPlatform(action *zrule.Action) (zrule.Dispatcher, error) {
	switch action.Platform {
	case zrule.PlatformDiscord:
		return discord.NewService(action, s.client)
	case zrule.PlatformSlack:
		return slack.NewService(action, s.client)
	case zrule.PlatformRest:
		return rest.NewService(action, s.client)
	default:
		return nil, fmt.Errorf("unknown platform specified")
	}

}
