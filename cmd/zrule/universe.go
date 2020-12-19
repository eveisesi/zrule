package main

import (
	"github.com/eveisesi/zrule/internal/esi"
	"github.com/eveisesi/zrule/internal/universe"
)

func newUniverseService(basics *app, repos repositories) universe.Service {

	esiServ := esi.NewService(basics.redis, "zrule v0.1.0")
	return universe.NewService(
		basics.redis, basics.newrelic, esiServ,
		repos.alliance, repos.corporation, repos.character,
		repos.region, repos.constellation, repos.system,
		repos.item, repos.itemGroup,
	)

}
