package universe

import (
	"context"

	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/esi"
	"github.com/go-redis/redis/v8"
	newrelic "github.com/newrelic/go-agent"
)

type SearchResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	SearchName(ctx context.Context, category, term string) ([]*SearchResult, error)

	zrule.UniverseRepository
}

type service struct {
	redis    *redis.Client
	newrelic *newrelic.Application

	esi esi.Service

	zrule.UniverseRepository
}

func NewService(redis *redis.Client, newrelic *newrelic.Application, esi esi.Service, universe zrule.UniverseRepository) Service {
	return &service{
		redis:    redis,
		newrelic: newrelic,

		esi: esi,

		UniverseRepository: universe,
	}
}
