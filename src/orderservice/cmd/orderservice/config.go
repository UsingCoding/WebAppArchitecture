package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	ServeHTTPAddress string `envconfig:"serve_http_address" default:":8000"`
	DatabaseDriver   string `envconfig:"db_driver" default:"mysql"`
	DSN              string `envconfig:"dsn" default:"root:1234@/orderservice"`
}
