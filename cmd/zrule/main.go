package main

import (
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/urfave/cli"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	err    error
)

type zrule struct {
	cfg      config
	newrelic *newrelic.Application
	logger   *logrus.Logger
	db       *mongo.Database
	redis    *redis.Client
	client   *nethttp.Client
}

// basics initializes the following
// loadConfig - parses environment variables and applies them to a struct
// loadLogger - takes in a configuration and intializes a logrus logger
// loadDB - takes in a configuration and establishes a connection with our datastore, in this application that is mongoDB
// loadRedis - takes in a configuration and establises a connection with our cache, in this application that is Redis
// loadNewrelic - takes in a configuration and configures a NR App to report metrics to NewRelic for monitoring
// loadClient - create a client from the net/http library that is used on all outgoing http requests
func basics(command string) *zrule {

	app := zrule{}

	app.cfg, err = loadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	app.logger, err = loadLogger(app.cfg, command)
	if err != nil {
		log.Fatalf("failed to load logger: %s", err)
	}

	app.cfg.NewRelic.AppName = fmt.Sprintf("%s-%s", app.cfg.NewRelic.AppName, command)

	app.newrelic, err = configNRApplication(app.cfg, app.logger)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure NR App")
	}

	app.db, err = makeMongoDB(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to make mongo db connection")
	}

	app.redis = makeRedis(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure redis client")
	}

	app.client = &nethttp.Client{
		Timeout:   time.Second * 5,
		Transport: newrelic.NewRoundTripper(nil),
	}
	return &app

}

func main() {
	app := cli.NewApp()
	app.Name = "ZRule"
	app.UsageText = "zrule"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "serve",
			Usage:  "Initializes the http server that handle http requests to this application",
			Action: serveCommand,
		},
		cli.Command{
			Name:    "processor",
			Aliases: []string{"p"},

			Usage:  "Initializes the processor that analyzes killmails to see if there is a policy match",
			Action: processorCommand,
		},
		cli.Command{
			Name:    "listener",
			Aliases: []string{"l"},
			Usage:   "Initializes the websocket client that listens to the ZKillboard Websocket",
			Action:  listenerCommand,
		},
		cli.Command{
			Name:    "dispatcher",
			Aliases: []string{"d"},
			// Usage:   "Initializes the websocket client that listens to the ZKillboard Websocket",
			Action: dispatcherCommand,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
