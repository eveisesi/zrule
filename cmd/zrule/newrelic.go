package main

import (
	"fmt"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func configNRApplication(cfg config, logger *logrus.Logger) (app *newrelic.Application, err error) {

	appName := cfg.NewRelic.AppName

	if cfg.Env != production {
		appName = fmt.Sprintf("%s-%s", cfg.Env, appName)
		if cfg.Developer.Name != "" {
			appName = fmt.Sprintf("%s-%s", cfg.Developer.Name, appName)
		}
	}

	opts := []newrelic.ConfigOption{}
	opts = append(opts, newrelic.ConfigFromEnvironment())
	opts = append(opts, newrelic.ConfigAppName(appName))
	opts = append(opts, newrelic.ConfigInfoLogger(logger.Writer()))

	app, err = newrelic.NewApplication(opts...)
	if err != nil {
		return nil, err
	}

	err = app.WaitForConnection(time.Second * 20)

	return

}
