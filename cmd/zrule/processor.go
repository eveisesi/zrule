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

	err = processor.NewService(
		basics.redis,
		basics.logger,
		basics.newrelic,
		policy.NewService(newUniverseService(basics), policyRepo),
	).Run(5)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to start processor service")
	}

}
