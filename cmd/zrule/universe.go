package main

import (
	"github.com/eveisesi/zrule/internal/esi"
	"github.com/eveisesi/zrule/internal/mdb"
	"github.com/eveisesi/zrule/internal/universe"
)

func newUniverseService(basics *zrule) universe.Service {

	allianceRepo, err := mdb.NewAllianceRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize allianceRepo")
	}

	basics.logger.Info("allianceRepo initialized")

	charactersRepo, err := mdb.NewCharacterRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize charactersRepo")
	}

	basics.logger.Info("charactersRepo initialized")

	constellationRepo, err := mdb.NewConstellationRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize constellationRepo")
	}

	basics.logger.Info("constellationRepo initialized")

	corporationRepo, err := mdb.NewCorporationRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize corporationRepo")
	}

	basics.logger.Info("corporationRepo initialized")

	itemRepo, err := mdb.NewItemRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize itemRepo")
	}

	basics.logger.Info("itemRepo initialized")

	itemGroupRepo, err := mdb.NewItemGroupRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize itemGroupRepo")
	}

	basics.logger.Info("itemGroupRepo initialized")

	regionRepo, err := mdb.NewRegionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize regionRepo")
	}

	basics.logger.Info("regionRepo initialized")

	solarSystemRepo, err := mdb.NewSolarSystemRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize solarSystemRepo")
	}

	esiServ := esi.NewService(basics.redis, "zrule v0.1.0")
	return universe.NewService(
		basics.redis, basics.newrelic, esiServ,
		allianceRepo, corporationRepo, charactersRepo,
		regionRepo, constellationRepo, solarSystemRepo,
		itemRepo, itemGroupRepo,
	)

}
