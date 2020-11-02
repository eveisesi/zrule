package main

import (
	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/dispatcher"
	"github.com/eveisesi/zrule/internal/mdb"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/urfave/cli"
)

func dispatcherCommand(c *cli.Context) {
	basics := basics("dispatcher")

	policyRepo, err := mdb.NewPolicyRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize policyRepository")
	}

	actionRepo, err := mdb.NewActionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize actionRepo")
	}

	err = dispatcher.NewService(
		basics.redis,
		basics.logger,
		basics.newrelic,
		basics.client,
		policy.NewService(policyRepo),
		action.NewService(actionRepo),
	).Run()
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize and run dispatcher service")
	}
}
