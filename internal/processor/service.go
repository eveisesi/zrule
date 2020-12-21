package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/zrule/internal/universe"

	"github.com/eveisesi/zrule/pkg/ruler"

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

	universe universe.Service
	policy   policy.Service
}

func NewService(
	redis *redis.Client,
	logger *logrus.Logger,
	newrelic *newrelic.Application,

	policy policy.Service,
	universe universe.Service,

) Service {

	s := &service{
		redis:    redis,
		logger:   logger,
		newrelic: newrelic,
		policy:   policy,
		universe: universe,
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

	// limiter := limiter.NewConcurrencyLimiter(int(limit))
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
			txn.End()
			continue
		}

		if s.trackers.expiresAt.Before(time.Now()) {
			err := s.initializeTracker(ctx)
			if err != nil {
				err = fmt.Errorf("failed to initialize trackers: %w", err)
				txn.NoticeError(err)
				txn.End()
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
			txn.End()
			continue
		}

		if restart == 1 {
			err = s.initializeTracker(ctx)
			if err != nil {
				err = fmt.Errorf("failed to initialize trackers: %w", err)
				txn.NoticeError(err)
				txn.End()
				s.logger.WithError(err).Errorln()
				return err
			}
			_, err := s.redis.Set(ctx, zrule.QUEUE_RESTART_TRACKER, 0, 0).Result()
			if err != nil {
				txn.NoticeError(err)
				s.logger.WithError(err).Fatal("error encountered attempting to create stop flag with default value")
			}
			time.Sleep(time.Second * 5)
			txn.End()
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
			txn.End()
			continue
		}

		if stop == 1 {
			s.logger.Info("stop signal set")

			s.logger.Info("sleeping for 5 seconds")
			time.Sleep(time.Second * 5)
			txn.End()
			continue
		}

		count, err := s.redis.ZCount(ctx, zrule.QUEUES_KILLMAIL_PROCESSING, "-inf", "+inf").Result()
		if err != nil {
			txn.NoticeError(err)
			txn.End()
			s.logger.WithError(err).Error("unable to determine count of message queue")
			time.Sleep(time.Second * 2)
			continue
		}

		if count == 0 {
			txn.End()
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
			// limiter.ExecuteWithTicket(func(workerID int) {
			s.handleMessage(newrelic.NewContext(ctx, s.newrelic.StartTransaction("handle message").NewGoroutine()), []byte(message))
			// })
		}

		txn.End()

	}

}

func (s *service) handleMessage(ctx context.Context, message []byte) {
	defer newrelic.FromContext(ctx).End()
	var killmail *zrule.Killmail
	err := json.Unmarshal(message, &killmail)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.logger.WithError(err).Errorf("failed to decode message onto %T", killmail)
	}

	// Hydrate the Killmail with Constellation and Region ID based on the SolarSystem

	s.hydrateKillmail(ctx, killmail)

	for _, tracker := range s.trackers.trackers {
		result := tracker.ruler.Test(killmail)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			s.logger.WithError(err).WithFields(logrus.Fields{
				"policyID": tracker.policy.ID,
			}).Error("failed to apply rules to killmail")
		}

		if !result {
			continue
		}

		entry := s.logger.WithField("policyID", tracker.policy.ID.Hex()).WithField("policyName", tracker.policy.Name)
		entry.Info("handling match")

		payload := zrule.Dispatchable{
			PolicyID: tracker.policy.ID,
			ID:       killmail.ID,
			Hash:     killmail.Hash,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			s.logger.WithError(err).Error("failed to marsahl payload for successfully match")
		}

		_, err = s.redis.ZAdd(ctx, zrule.QUEUES_KILLMAIL_MATCHED, &redis.Z{Score: float64(time.Now().UnixNano()), Member: string(data)}).Result()
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			s.logger.WithError(err).WithField("payload", string(message)).Error("unable to push payload to matched queue")
		}

	}

}

func (s *service) hydrateKillmail(ctx context.Context, killmail *zrule.Killmail) {

	entry := s.logger.WithField("killmail_id", killmail.ID)
	if killmail.Meta != nil {
		killmail.Hash = killmail.Meta.Hash
		entry = entry.WithField("killmail_hash", killmail.Hash)
	}

	system, err := s.universe.SolarSystem(ctx, killmail.SolarSystemID)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		entry.WithError(err).WithField("SolarSystemID", killmail.SolarSystemID).Debug("failed to look up solar system for solar system")
		return
	}

	constellation, err := s.universe.Constellation(ctx, system.ConstellationID)
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		entry.WithError(err).WithField("SolarSystemID", killmail.SolarSystemID).WithField("ConstellationID", system.ConstellationID).Debug("failed to look up constellation for solar system")
		return
	}

	killmail.ConstellationID = constellation.ID
	killmail.RegionID = constellation.RegionID

	if killmail.Victim != nil {
		victimShip, err := s.universe.Item(ctx, killmail.Victim.ShipTypeID)
		if err != nil {
			newrelic.FromContext(ctx).NoticeError(err)
			entry.WithError(err).WithField("Victim.ShipTypeID", killmail.Victim.ShipTypeID).Debug("failed to lookup victim ship")
		}
		killmail.Victim.ShipGroupID = victimShip.GroupID
	}

	if len(killmail.Attackers) > 0 {
		for i, attacker := range killmail.Attackers {

			if attacker.ShipTypeID != nil {
				attackerShip, err := s.universe.Item(ctx, *attacker.ShipTypeID)
				if err != nil {
					newrelic.FromContext(ctx).NoticeError(err)
					entry.WithError(err).WithField("attackerID", i).Debug("failed to lookup attacker ship")
					continue
				}

				attacker.ShipGroupID = &attackerShip.GroupID

			}

			if attacker.WeaponTypeID != nil {
				attackerWeapon, err := s.universe.Item(ctx, *attacker.WeaponTypeID)
				if err != nil {
					newrelic.FromContext(ctx).NoticeError(err)
					entry.WithError(err).WithField("attackerID", i).Debug("failed to lookup attacker ship")
					continue
				}

				attacker.WeaponGroupID = &attackerWeapon.GroupID

			}

		}
	}

}
