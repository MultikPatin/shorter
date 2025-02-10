package app

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"strconv"
	"strings"
)

const (
	urlPrefix   = "http://"
	delimiter   = "/"
	contentType = "text/plain; charset=utf-8"
)

type Config struct {
	Addr            string
	ShortLinkPrefix string
}

type envConfig struct {
	Addr            string `env:"SERVER_ADDRESS,required"`
	ShortLinkPrefix string `env:"BASE_URL,required"`
}
type cmdConfig struct {
	ServHost ServHost
	ShorLink ShorLink
}
type ServHost struct {
	Host string
	Port int
}
type ShorLink struct {
	ShortLinkPrefix string
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.parseEnv(); err != nil {
		if err := cfg.parseFlags(); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func (c *Config) parseEnv() error {
	cfg := &envConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	c.Addr = cfg.Addr
	c.ShortLinkPrefix = cfg.ShortLinkPrefix
	return nil
}

func (c *Config) parseFlags() error {
	sv := new(ServHost)
	_ = flag.Value(sv)
	sh := new(ShorLink)
	_ = flag.Value(sh)

	flag.Var(sv, "a", "Net address host:port")
	flag.Var(sh, "b", "short link server")
	flag.Parse()

	c.Addr = sv.String()
	c.ShortLinkPrefix = sh.String()
	return nil
}

func (a *ServHost) String() string {
	a.normalize()
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *ServHost) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	a.normalize()
	return nil
}

func (a *ServHost) normalize() {
	if a.Port == 0 {
		a.Port = 8080
	}
	if a.Host == "" {
		a.Host = "localhost"
	}
}

func (a *ShorLink) String() string {
	return a.ShortLinkPrefix
}

func (a *ShorLink) Set(s string) error {
	hp := strings.Split(s, ":")
	a.ShortLinkPrefix = hp[0]
	return nil
}
