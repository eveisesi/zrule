package main

import (
	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/internal/mdb"
)

type repositories struct {
	alliance      zrule.AllianceRepository
	character     zrule.CharacterRepository
	corporation   zrule.CorporationRepository
	region        zrule.RegionRepository
	constellation zrule.ConstellationRepository
	system        zrule.SolarSystemRepository
	item          zrule.ItemRepository
	itemGroup     zrule.ItemGroupRepository
}

func initializeRepositories(basics *app) repositories {

	repos := repositories{}

	repos.alliance, err = mdb.NewAllianceRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize allianceRepo")
	}

	basics.logger.Info("allianceRepo initialized")

	repos.character, err = mdb.NewCharacterRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize charactersRepo")
	}

	basics.logger.Info("charactersRepo initialized")

	repos.constellation, err = mdb.NewConstellationRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize constellationRepo")
	}

	basics.logger.Info("constellationRepo initialized")

	repos.corporation, err = mdb.NewCorporationRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize corporationRepo")
	}

	basics.logger.Info("corporationRepo initialized")

	repos.item, err = mdb.NewItemRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize itemRepo")
	}

	basics.logger.Info("itemRepo initialized")

	repos.itemGroup, err = mdb.NewItemGroupRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize itemGroupRepo")
	}

	basics.logger.Info("itemGroupRepo initialized")

	repos.region, err = mdb.NewRegionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize regionRepo")
	}

	basics.logger.Info("regionRepo initialized")

	repos.system, err = mdb.NewSolarSystemRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize solarSystemRepo")
	}

	return repos

}
