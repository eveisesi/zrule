package main

import (
	"github.com/eveisesi/zrule/internal/mdb"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/eveisesi/zrule/internal/processor"
	"github.com/urfave/cli"
)

func processorCommand(c *cli.Context) {

	basics := basics("processor")

	policyRepo, err := mdb.NewPolicyRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize policyRepository")
	}

	basics.logger.Info("policyRepo initialized")
	repos := initializeRepositories(basics)

	universeServ := newUniverseService(basics, repos)

	err = processor.NewService(
		basics.redis,
		basics.logger,
		basics.newrelic,
		policy.NewService(universeServ, policyRepo),
		universeServ,
	).Run(5)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to start processor service")
	}
}
