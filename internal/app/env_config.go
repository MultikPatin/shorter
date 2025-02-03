package app

import (
	"github.com/caarlos0/env/v6"
	"log"
	"net/url"
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
	_, err = url.Parse(c.ShorLink)
	if err != nil {
		c.ShorLink = ""
	}
	return err
}
