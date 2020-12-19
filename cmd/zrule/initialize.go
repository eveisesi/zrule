package main

import (
	"context"
	"sync"
	"time"

	"github.com/eveisesi/zrule/internal/esi"
	"github.com/urfave/cli"
)

func initializeCommand(c *cli.Context) {
	basics := basics("dispatcher")

	ctx := context.Background()

	esiServ := esi.NewService(basics.redis, "zrule v0.1.0")

	status, m := esiServ.GetStatus(ctx)
	if m.IsErr() {
		basics.logger.WithError(m.Msg).Fatal("failed to fetch ESI Server Status")
	}

	time.Sleep(time.Second)
	basics.logger.WithField("server_version", status.ServerVersion).Info("ESI Server Status OK, proceeding with initialization")
	time.Sleep(time.Second)

	repos := initializeRepositories(basics)
	universe := newUniverseService(basics, repos)

	wg := &sync.WaitGroup{}

	for _, categoryID := range []uint{6, 7, 8, 18, 22, 32, 87} {
		category, m := esiServ.GetUniverseCategoriesCategoryID(ctx, categoryID)
		if m.IsErr() {
			basics.logger.WithError(m.Msg).Fatal("received error from esi when fetch ship category")
		}

		// Loop over the Category Groups and call ESI for each Group to get its types
		for _, groupID := range category.Groups {
			wg.Add(1)
			go func(basics *app, wg *sync.WaitGroup, groupID uint) {
				defer wg.Done()
				group, m := esiServ.GetUniverseGroupsGroupID(ctx, groupID)
				if m.IsErr() {
					basics.logger.WithError(m.Msg).WithField("groupID", groupID).Error("received error from esi when fetching group")
					return
				}

				// Now loop over each of the Group Types and call ESI for each Type
				for _, typeID := range group.Types {
					item, m := esiServ.GetUniverseTypesTypeID(ctx, typeID)
					if m.IsErr() {
						basics.logger.WithError(m.Msg).WithField("groupID", groupID).WithField("typeID", typeID).Error("received error from esi when fetching type")
						return
					}

					_, err := universe.CreateItem(ctx, item)
					if err != nil {
						basics.logger.WithError(err).WithField("groupID", groupID).WithField("typeID", typeID).Error("failed to save type to database")
						return
					}

					basics.logger.WithField("groupID", groupID).WithField("typeID", typeID).Info("type saved successfully")

				}

				_, err := universe.CreateItemGroup(ctx, group)
				if err != nil {
					basics.logger.WithError(err).WithField("groupID", groupID).Error("failed to save group to database")
					return
				}

				basics.logger.WithField("groupID", groupID).Info("group saved successfully")
			}(basics, wg, groupID)
		}

		wg.Wait()
	}
	basics.logger.Info("groups and items imported successfully")

	regions, m := esiServ.GetUniverseRegions(ctx)
	if err != nil {
		basics.logger.WithError(m.Msg).Error("failed to fetch regions from ESI")
		return
	}

	for _, regionID := range regions {
		wg.Add(1)
		go func(basics *app, wg *sync.WaitGroup, regionID uint) {
			defer wg.Done()
			region, m := esiServ.GetUniverseRegionsRegionID(ctx, regionID)
			if err != nil {
				basics.logger.WithError(m.Msg).WithField("regionID", regionID).Error("failed to fetch region from ESI")
				return
			}

			for _, constellationID := range region.Constellations {
				constellation, m := esiServ.GetUniverseConstellationsConstellationID(ctx, constellationID)
				if err != nil {
					basics.logger.WithError(m.Msg).WithField("regionID", regionID).WithField("constellationID", constellationID).Error("failed to fetch constellation from ESI")
					return
				}

				for _, systemID := range constellation.Systems {
					system, m := esiServ.GetUniverseSolarSystemsSolarSystemID(ctx, systemID)
					if err != nil {
						basics.logger.WithError(m.Msg).WithField("regionID", regionID).WithField("constellationID", constellationID).WithField("systemID", systemID).Error("failed to fetch system from ESI")
						return
					}

					_, err := universe.CreateSolarSystem(ctx, system)
					if err != nil {
						basics.logger.WithError(err).WithField("regionID", regionID).WithField("constellationID", constellationID).WithField("systemID", systemID).Error("failed to save system to database")
						return
					}

					basics.logger.WithField("regionID", regionID).WithField("constellationID", constellationID).WithField("systemID", systemID).Info("system saved successfully")

				}

				_, err := universe.CreateConstellation(ctx, constellation)
				if err != nil {
					basics.logger.WithError(err).WithField("regionID", regionID).WithField("constellationID", constellationID).Error("failed to save constellation to database")
					return
				}

				basics.logger.WithField("regionID", regionID).WithField("constellationID", constellationID).Info("constellation saved successfully")

			}

			_, err := universe.CreateRegion(ctx, region)
			if err != nil {
				basics.logger.WithError(err).WithField("regionID", regionID).Error("failed to save region to database")
				return
			}

			basics.logger.WithField("regionID", regionID).Info("region saved successfully")
		}(basics, wg, regionID)

	}

	wg.Wait()
	basics.logger.Info("All Universe Info Initialized Successfully")

}
