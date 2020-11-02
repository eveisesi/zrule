package character

import (
	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/esi"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service interface {
	// UpdateExpired(ctx context.Context)
	// CharactersByCharacterIDs(ctx context.Context, ids []uint64) ([]*zrule.Character, error)
	zrule.CharacterRespository
}

type service struct {
	redis    *redis.Client
	logger   *logrus.Logger
	newrelic *newrelic.Application
	esi      esi.Service
	// tracker  tracker.Service
	zrule.CharacterRespository
}

// tracker tracker.Service
// NewService initializes the Character Service and returns and interface expossing its functionality
func NewService(redis *redis.Client, logger *logrus.Logger, newrelic *newrelic.Application, esi esi.Service, character zrule.CharacterRespository) Service {
	return &service{
		redis,
		logger,
		newrelic,
		esi,
		// tracker,
		character,
	}
}
