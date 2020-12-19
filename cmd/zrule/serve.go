package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eveisesi/zrule/internal/search"

	"golang.org/x/oauth2"

	"github.com/eveisesi/zrule/internal/action"
	"github.com/eveisesi/zrule/internal/dispatcher"
	"github.com/eveisesi/zrule/internal/http"
	"github.com/eveisesi/zrule/internal/mdb"
	"github.com/eveisesi/zrule/internal/policy"
	"github.com/eveisesi/zrule/internal/token"
	"github.com/eveisesi/zrule/internal/user"

	"github.com/urfave/cli"
)

func serveCommand(c *cli.Context) {

	basics := basics("http")

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

	basics.logger.Info("tokenServ service initialized")

	repos := initializeRepositories(basics)

	basics.logger.Info("tokenServ service initialized")

	searchServ := search.NewService(basics.logger, fmt.Sprintf("%s:%d", basics.cfg.Redis.Host, basics.cfg.Redis.Port))
	if basics.cfg.Redis.InitializeAutocompleter {
		initializeAutocompleter(basics, repos, searchServ)
	}

	basics.logger.Info("searchServ service initialized")

	universeServ := newUniverseService(basics, repos)

	actionServ := action.NewService(actionRepo)
	userServ := user.NewService(basics.logger, basics.redis, tokenServ, universeServ, userRepo)
	policyServ := policy.NewService(universeServ, policyRepo)

	dispacther := dispatcher.NewService(
		basics.redis,
		basics.logger,
		basics.newrelic,
		basics.client,
		policyServ,
		actionServ,
	)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize and run dispatcher service")
	}

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
		universeServ,
		dispacther,
		searchServ,
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
