package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/zrule/pkg/ruler"

	"github.com/korovkin/limiter"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/policy"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run(limit int64) error
}

type tracker struct {
	policy *zrule.Policy
	ruler  *ruler.Ruler
}

type policyTracker struct {
	expiresAt time.Time
	trackers  []tracker
}

type service struct {
	redis    *redis.Client
	logger   *logrus.Logger
	newrelic *newrelic.Application
	trackers *policyTracker

	policy policy.Service
}

func NewService(
	redis *redis.Client,
	logger *logrus.Logger,
	newrelic *newrelic.Application,

	policy policy.Service,
) Service {

	s := &service{
		redis:    redis,
		logger:   logger,
		newrelic: newrelic,
		policy:   policy,
	}

	err := s.initializeTracker(context.Background())
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize tracker")
	}

	return s

}

func (s *service) initializeTracker(ctx context.Context) error {
	s.logger.Info("starting tracker initialization")
	defer s.logger.Info("finishing tracker initialization")
	policies, err := s.policy.Policies(ctx, zrule.NewEqualOperator("paused", false))
	if err != nil {
		return err
	}

	trackers := make([]tracker, len(policies))

	for i, policy := range policies {
		if len(policy.Rules) == 0 {
			continue
		}

		data, err := json.Marshal(policy.Rules)
		if err != nil {
			panic(fmt.Errorf("failed to marsahl policy rules: %w", err))
		}

		ruler := ruler.NewRuler()
		ruler.SetRulesWithJSON(data)

		trackers[i].policy = policy
		trackers[i].ruler = ruler

	}

	s.trackers = &policyTracker{
		expiresAt: time.Now().Add(time.Minute * 5),
		trackers:  trackers,
	}

	return nil

}

func (s *service) Run(limit int64) error {

	limiter := limiter.NewConcurrencyLimiter(int(limit))
	for {
		var ctx = context.Background()
		txn := s.newrelic.StartTransaction("killmail queue check")

		restart, err := s.redis.Get(ctx, zrule.QUEUE_RESTART_TRACKER).Int64()
		if err != nil {
			s.logger.WithError(err).Error("restart flag is missing. attempting to create with default value of 0")
			_, err := s.redis.Set(ctx, zrule.QUEUE_RESTART_TRACKER, 0, 0).Result()
			if err != nil {
				txn.NoticeError(err)
				s.logger.WithError(err).Fatal("error encountered attempting to create stop flag with default value")
			}
			continue
		}

		if s.trackers.expiresAt.Before(time.Now()) {
			err := s.initializeTracker(ctx)
			if err != nil {
				err = fmt.Errorf("failed to initialize trackers: %w", err)
				s.logger.WithError(err).Errorln()
				return err
			}

			if restart == 1 {
				_, err = s.redis.Set(ctx, zrule.QUEUE_RESTART_TRACKER, 0, 0).Result()
				if err != nil {
					txn.NoticeError(err)
					s.logger.WithError(err).Fatal("error encountered attempting to create stop flag with default value")
				}

			}
			time.Sleep(time.Second * 5)
			txn.Ignore()
			continue
		}

		if restart == 1 {
			err = s.initializeTracker(ctx)
			if err != nil {
				err = fmt.Errorf("failed to initialize trackers: %w", err)
				s.logger.WithError(err).Errorln()
				return err
			}
			_, err := s.redis.Set(ctx, zrule.QUEUE_RESTART_TRACKER, 0, 0).Result()
			if err != nil {
				txn.NoticeError(err)
				s.logger.WithError(err).Fatal("error encountered attempting to create stop flag with default value")
			}
			time.Sleep(time.Second * 5)
			txn.Ignore()
			continue
		}

		stop, err := s.redis.Get(ctx, zrule.QUEUE_STOP).Int64()
		if err != nil {
			s.logger.WithError(err).Error("stop flag is missing. attempting to create with default value of 0")
			_, err := s.redis.Set(ctx, zrule.QUEUE_STOP, 0, 0).Result()
			if err != nil {
				txn.NoticeError(err)
				s.logger.WithError(err).Fatal("error encountered attempting to create stop flag with default value")
			}
			continue
		}

		if stop == 1 {
			s.logger.Info("stop signal set")

			s.logger.Info("sleeping for 5 seconds")
			time.Sleep(time.Second * 5)
			txn.Ignore()
			continue
		}

		count, err := s.redis.ZCount(ctx, zrule.QUEUES_KILLMAIL_PROCESSING, "-inf", "+inf").Result()
		if err != nil {
			txn.NoticeError(err)
			s.logger.WithError(err).Error("unable to determine count of message queue")
			time.Sleep(time.Second * 2)
			continue
		}

		if count == 0 {
			txn.Ignore()
			s.logger.Info("processing queue is empty")
			time.Sleep(time.Second * 2)
			continue
		}

		results, err := s.redis.ZPopMax(ctx, zrule.QUEUES_KILLMAIL_PROCESSING, limit).Result()
		if err != nil {
			txn.NoticeError(err)
			s.logger.WithError(err).Fatal("unable to retrieve hashes from queue")
		}

		for _, result := range results {
			message := result.Member.(string)
			s.logger.Info("handling message")
			limiter.ExecuteWithTicket(func(workerID int) {
				s.handleMessage(ctx, []byte(message))
				time.Sleep(time.Millisecond * 500)
			})
		}

		txn.End()

	}

}

func (s *service) handleMessage(ctx context.Context, message []byte) {

	var killmail map[string]interface{}
	err := json.Unmarshal(message, &killmail)
	if err != nil {
		s.logger.WithError(err).Errorf("failed to decode message onto %T", killmail)
	}

	for _, tracker := range s.trackers.trackers {
		result, err := tracker.ruler.Test(killmail)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"policyID": tracker.policy.ID,
			}).Error("failed to apply rules to killmail")
		}

		if !result {
			continue
		}

		entry := s.logger.WithField("policyID", tracker.policy.ID.Hex()).WithField("policyName", tracker.policy.Name)
		entry.Info("handling match")
		var (
			ok   bool
			id   uint
			hash string
		)

		_, ok = killmail["killmail_id"]
		if !ok {
			entry.Error("cannot find killmail_id")
			continue
		}

		_, ok = killmail["killmail_id"].(float64)
		if !ok {
			entry.Error("killmail_id is not of type float64")
			continue
		}

		id = uint(killmail["killmail_id"].(float64))

		zkb, ok := killmail["zkb"]
		if !ok {
			entry.Error("cannot find zkb object")
			continue
		}
		_, ok = zkb.(map[string]interface{})
		if !ok {
			entry.Error("zkb object is not of type map[string]interface{}{}")
			continue
		}
		_, ok = zkb.(map[string]interface{})["hash"]
		if !ok {
			entry.Error("cannot find zkb.hash")
			continue
		}
		hash, ok = zkb.(map[string]interface{})["hash"].(string)
		if !ok {
			entry.Error("zkb.hash is not of type string")
			continue
		}

		if hash == "" {
			continue
		}

		payload := zrule.Dispatchable{
			PolicyID: tracker.policy.ID,
			ID:       id,
			Hash:     hash,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			s.logger.WithError(err).Error("failed to marsahl payload for successfully match")
		}

		_, err = s.redis.ZAdd(ctx, zrule.QUEUES_KILLMAIL_MATCHED, &redis.Z{Score: float64(time.Now().UnixNano()), Member: string(data)}).Result()
		if err != nil {
			s.logger.WithError(err).WithField("payload", string(message)).Error("unable to push payload to matched queue")
		}

	}

}
