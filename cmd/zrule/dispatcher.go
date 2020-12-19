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

	// Initialize Repositories
	actionRepo, err := mdb.NewActionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize actionRepo")
	}

	basics.logger.Info("actionRepo initialized")

	policyRepo, err := mdb.NewPolicyRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize policyRepo")
	}

	basics.logger.Info("policyRepo initialized")
	repos := initializeRepositories(basics)
	err = dispatcher.NewService(
		basics.redis,
		basics.logger,
		basics.newrelic,
		basics.client,
		policy.NewService(newUniverseService(basics, repos), policyRepo),
		action.NewService(actionRepo),
	).Run()
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize and run dispatcher service")
	}
}
