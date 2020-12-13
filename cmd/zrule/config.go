package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type environment string

const production environment = "production"
const development environment = "development"

func (e environment) String() string {
	return string(e)
}

type config struct {
	Mongo struct {
		Host     string
		Port     int
		User     string
		Pass     string
		Name     string
		AuthMech string `default:"SCRAM-SHA-1"`
		Sleep    int    `default:"5"`
	}

	Redis struct {
		Host string
		Port uint
	}

	NewRelic struct {
		Enabled bool   `envconfig:"NEW_RELIC_ENABLED" default:"false"`
		AppName string `envconfig:"NEW_RELIC_APP_NAME"`
	}

	Env environment

	Developer struct {
		Name string
	}

	Log struct {
		Level string
	}

	Server struct {
		Port uint
	}

	Auth struct {
		ClientID         string
		ClientSecret     string
		RedirectURL      string
		AuthorizationURL string
		TokenURL         string
		JWKSURL          string
	}
}

var validEnvironment = []environment{production, development}

func (c config) validateEnvironment() bool {
	for _, env := range validEnvironment {
		if c.Env == env {
			return true
		}
	}

	return false
}

func loadConfig() (cfg config, err error) {
	_ = godotenv.Load(".env")

	err = envconfig.Process("", &cfg)
	if err != nil {
		return config{}, err
	}

	if !cfg.validateEnvironment() {
		return config{}, fmt.Errorf("invalid env %s declared", cfg.Env)
	}

	return

}
