package app

import (
	"github.com/caarlos0/env/v6"
	"log"
)

var EnvConfig envConfig

type envConfig struct {
	ServHost string `env:"SERVER_ADDRESS"`
	ShorLink string `env:"BASE_URL"`
}

func (c *envConfig) Parse() error {
	err := env.Parse(&EnvConfig)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
