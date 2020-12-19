package universe

import (
	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/esi"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type SearchResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	zrule.AllianceRepository
	zrule.CorporationRepository
	zrule.CharacterRepository
	zrule.RegionRepository
	zrule.ConstellationRepository
	zrule.SolarSystemRepository

	zrule.ItemRepository
	zrule.ItemGroupRepository
}

type service struct {
	redis    *redis.Client
	newrelic *newrelic.Application

	esi esi.Service

	zrule.AllianceRepository
	zrule.CorporationRepository
	zrule.CharacterRepository
	zrule.RegionRepository
	zrule.ConstellationRepository
	zrule.SolarSystemRepository

	zrule.ItemRepository
	zrule.ItemGroupRepository
}

func NewService(
	redis *redis.Client,
	newrelic *newrelic.Application,
	esi esi.Service,

	alliance zrule.AllianceRepository,
	corporation zrule.CorporationRepository,
	character zrule.CharacterRepository,
	region zrule.RegionRepository,
	constellation zrule.ConstellationRepository,
	solarSystem zrule.SolarSystemRepository,

	item zrule.ItemRepository,
	itemGroup zrule.ItemGroupRepository,
) Service {
	return &service{
		redis:    redis,
		newrelic: newrelic,

		esi: esi,

		AllianceRepository:      alliance,
		CorporationRepository:   corporation,
		CharacterRepository:     character,
		RegionRepository:        region,
		ConstellationRepository: constellation,
		SolarSystemRepository:   solarSystem,
		ItemRepository:          item,
		ItemGroupRepository:     itemGroup,
	}
}
