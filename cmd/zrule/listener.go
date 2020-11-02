package main

import (
	"github.com/eveisesi/zrule/internal/websocket"
	"github.com/urfave/cli"
)

func listenerCommand(c *cli.Context) {

	basics := basics("websocket")

	err = websocket.NewService(
		basics.redis,
		basics.logger,
		basics.newrelic,
	).Run()
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to start websocket listener")
	}

}
