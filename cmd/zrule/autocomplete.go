package main

import (
	"context"

	"github.com/eveisesi/zrule/internal/search"
	"github.com/jinzhu/copier"
)

func initializeAutocompleter(basics *app, repos repositories, searchServ search.Service) {
	ctx := context.Background()

	entry := basics.logger.WithField("autocompleter", "regions")
	entry.Info("initializing autocompleter")

	regions, _ := repos.region.Regions(ctx)
	entities := make([]*search.Entity, len(regions))
	for i, region := range regions {
		entity := new(search.Entity)
		err = copier.Copy(entity, region)
		if err != nil {
			basics.logger.WithError(err).Error("failed to copy region to generic entity")
			continue
		}
		entities[i] = entity
	}
	err := searchServ.InitializeAutocompleter(search.KeyRegions, entities)
	if err != nil {
		entry.WithError(err).Fatal("failed to initialize autocompleter")
	}

	entry.Info("autocompleter initialized successfully")
	entry = basics.logger.WithField("autocompleter", "constellations")
	entry.Info("initializing autocompleter")

	constellations, _ := repos.constellation.Constellations(ctx)
	entities = make([]*search.Entity, len(constellations))
	for i, constellation := range constellations {
		entity := new(search.Entity)
		err = copier.Copy(entity, constellation)
		if err != nil {
			basics.logger.WithError(err).Error("failed to copy region to generic entity")
			continue
		}
		entities[i] = entity
	}
	err = searchServ.InitializeAutocompleter(search.KeyConstellations, entities)
	if err != nil {
		entry.WithError(err).Fatal("failed to initialize autocompleter")
	}

	entry.Info("autocompleter initialized successfully")
	entry = basics.logger.WithField("autocompleter", "systems")
	entry.Info("initializing autocompleter")

	systems, _ := repos.system.SolarSystems(ctx)
	entities = make([]*search.Entity, len(systems))
	for i, system := range systems {
		entity := new(search.Entity)
		err = copier.Copy(entity, system)
		if err != nil {
			basics.logger.WithError(err).Error("failed to copy region to generic entity")
			continue
		}
		entities[i] = entity
	}
	err = searchServ.InitializeAutocompleter(search.KeySystems, entities)
	if err != nil {
		entry.WithError(err).Fatal("failed to initialize autocompleter")
	}

	entry.Info("autocompleter initialized successfully")
	entry = basics.logger.WithField("autocompleter", "itemGroups")
	entry.Info("initializing autocompleter")

	itemGroups, _ := repos.itemGroup.ItemGroups(ctx)
	entities = make([]*search.Entity, len(itemGroups))
	for i, group := range itemGroups {
		entity := new(search.Entity)
		err = copier.Copy(entity, group)
		if err != nil {
			basics.logger.WithError(err).Error("failed to copy region to generic entity")
			continue
		}
		entities[i] = entity
	}
	err = searchServ.InitializeAutocompleter(search.KeyItemGroups, entities)
	if err != nil {
		entry.WithError(err).Fatal("failed to initialize autocompleter")
	}

	entry.Info("autocompleter initialized successfully")
	entry = basics.logger.WithField("autocompleter", "items")
	entry.Info("initializing autocompleter")

	items, _ := repos.item.Items(ctx)
	entities = make([]*search.Entity, len(items))
	for i, item := range items {
		entity := new(search.Entity)
		err = copier.Copy(entity, item)
		if err != nil {
			basics.logger.WithError(err).Error("failed to copy region to generic entity")
			continue
		}
		entities[i] = entity
	}
	err = searchServ.InitializeAutocompleter(search.KeyItems, entities)
	if err != nil {
		entry.WithError(err).Fatal("failed to initialize autocompleter")
	}

	entry.Info("autocompleter initialized successfully")

	entry = basics.logger.WithField("autocompleter", "factions")
	entry.Info("initializing autocompleter")

	factions, _ := repos.faction.Factions(ctx)
	entities = make([]*search.Entity, len(factions))
	for i, faction := range factions {
		entity := new(search.Entity)
		err = copier.Copy(entity, faction)
		if err != nil {
			basics.logger.WithError(err).Error("failed to copy region to generic entity")
			continue
		}
		entities[i] = entity
	}
	err = searchServ.InitializeAutocompleter(search.KeyFactions, entities)
	if err != nil {
		entry.WithError(err).Fatal("failed to initialize autocompleter")
	}

	entry.Info("autocompleter initialized successfully")
}
