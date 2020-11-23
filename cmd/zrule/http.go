package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eveisesi/zrule/internal/universe"

	"golang.org/x/oauth2"

	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/esi"
	"github.com/eveisesi/zrule/internal/http"
	"github.com/eveisesi/zrule/internal/mdb"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/eveisesi/zrule/internal/token"
	"github.com/eveisesi/zrule/internal/user"

	"github.com/urfave/cli"
)

func httpCommand(c *cli.Context) {

	basics := basics("http")

	// Initialize Repositories
	actionRepo, err := mdb.NewActionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize actionRepo")
	}

	basics.logger.Info("actionRepo initialized")

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

	policyRepo, err := mdb.NewPolicyRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize policyRepo")
	}

	basics.logger.Info("policyRepo initialized")

	regionRepo, err := mdb.NewRegionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize regionRepo")
	}

	basics.logger.Info("regionRepo initialized")

	solarSystemRepo, err := mdb.NewSolarSystemRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize solarSystemRepo")
	}

	basics.logger.Info("solarSystemRepo initialized")

	userRepo, err := mdb.NewUserRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize userRepo")
	}

	basics.logger.Info("userRepo initialized")

	tokenServ := token.NewService(
		basics.client,
		&oauth2.Config{
			ClientID:     basics.cfg.Auth.ClientID,
			ClientSecret: basics.cfg.Auth.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   basics.cfg.Auth.AuthorizationURL,
				TokenURL:  basics.cfg.Auth.TokenURL,
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
		basics.logger,
		basics.redis,
		basics.cfg.Auth.JWKSURL,
	)

	esiServ := esi.NewService(basics.redis, "zrule v0.1.0")
	actionServ := action.NewService(actionRepo)
	universeSer := universe.NewService(
		basics.redis, basics.newrelic, esiServ,
		allianceRepo, corporationRepo, charactersRepo,
		regionRepo, constellationRepo, solarSystemRepo,
		itemRepo, itemGroupRepo,
	)
	userServ := user.NewService(basics.logger, basics.redis, tokenServ, universeSer, userRepo)
	policyServ := policy.NewService(policyRepo)

	server := http.NewServer(
		basics.cfg.Server.Port,
		basics.db,
		basics.logger,
		basics.redis,
		basics.newrelic,
		tokenServ,
		userServ,
		actionServ,
		policyServ,
		universeSer,
	)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.Run()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		basics.logger.WithError(err).Fatal("server encountered an unexpected error and had to quit")
	case sig := <-osSignals:
		basics.logger.WithField("sig", sig).Info("interrupt signal received, starting server shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = server.GracefullyShutdown(ctx)
		if err != nil {
			basics.logger.WithError(err).Fatal("failed to shutdown server")
		}

		basics.logger.Info("server gracefully shutdown successfully")
	}

}
