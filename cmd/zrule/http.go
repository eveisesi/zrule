package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/character"
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

	userRepo, err := mdb.NewUserRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize userRepo")
	}

	characterRepo, err := mdb.NewCharacterRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize userRepo")
	}

	actionRepo, err := mdb.NewActionRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize actionRepo")
	}

	policyRepo, err := mdb.NewPolicyRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize policyRepo")
	}

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
	characterServ := character.NewService(basics.redis, basics.logger, basics.newrelic, esiServ, characterRepo)
	userServ := user.NewService(basics.logger, basics.redis, tokenServ, characterServ, userRepo)
	policyServ := policy.NewService(policyRepo)

	server := http.NewServer(
		basics.cfg.Server.Port,
		basics.logger,
		basics.redis,
		basics.newrelic,
		tokenServ,
		userServ,
		actionServ,
		policyServ,
		characterServ,
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
