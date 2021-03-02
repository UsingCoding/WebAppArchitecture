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
	LogPath          string `envconfig:"log_dir" default:"./var/log/dev.log"`
	ServeHTTPAddress string `envconfig:"serve_http_address" default:":8000"`
}
