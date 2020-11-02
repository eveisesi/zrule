package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const production = "production"
const prod = production

const development = "development"
const dev = development

type config struct {
	Mongo struct {
		Host     string
		Port     int
		User     string
		Pass     string
		Name     string
		AuthMech string `default:"SCRAM-SHA-1"`
	}

	Redis struct {
		Host string
		Port uint
	}

	NewRelic struct {
		AppName string `envconfig:"NEW_RELIC_APP_NAME"`
	}

	Env string

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

func loadConfig() (config config, err error) {
	_ = godotenv.Load("./cmd/zrule/.env")

	err = envconfig.Process("", &config)

	if config.Env != prod {
		config.Env = dev
	}

	return

}
