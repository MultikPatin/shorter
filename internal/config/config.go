package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultStorageFilePath = "shorter"
)

type Config struct {
	PPofAddr         string
	Addr             string
	ShortLinkPrefix  string
	StorageFilePaths string
	PostgresDNS      *url.URL
}

type envConfig struct {
	StorageFilePaths string `env:"FILE_STORAGE_PATH"`
	Addr             string `env:"SERVER_ADDRESS"`
	ShortLinkPrefix  string `env:"BASE_URL"`
	PostgresDNS      string `env:"DATABASE_DSN"`
}
type cmdConfig struct {
	Addr             string
	StorageFilePaths string
	ShortLinkPrefix  string
	PostgresDNS      string
}

type ServHost struct {
	Host string
	Port int
}

func Parse(logger *zap.SugaredLogger) *Config {
	cfg := &Config{}

	envCfg, err := parseEnv()
	if err != nil {
		logger.Infow(
			"Parsed Env",
			"error", err.Error(),
		)
	}
	cmdCfg, err := parseCmd()
	if err != nil {
		logger.Infow(
			"Parsed CMD",
			"error", err.Error(),
		)
	}

	if envCfg.Addr == "" {
		cfg.Addr = cmdCfg.Addr
	} else {
		cfg.Addr = envCfg.Addr
	}
	if envCfg.ShortLinkPrefix == "" {
		cfg.ShortLinkPrefix = cmdCfg.ShortLinkPrefix
	} else {
		cfg.ShortLinkPrefix = envCfg.ShortLinkPrefix
	}
	if envCfg.StorageFilePaths == "" {
		cfg.StorageFilePaths = cmdCfg.StorageFilePaths
	} else {
		cfg.StorageFilePaths = envCfg.StorageFilePaths
	}
	if cfg.StorageFilePaths == "" {
		cfg.StorageFilePaths = defaultStorageFilePath
	}
	if envCfg.PostgresDNS != "" {
		cfg.PostgresDNS, _ = parseDSN(envCfg.PostgresDNS)
	} else if cmdCfg.PostgresDNS != "" {
		cfg.PostgresDNS, _ = parseDSN(cmdCfg.PostgresDNS)
	}
	cfg.PPofAddr = "localhost:6060"

	return cfg
}

func parseEnv() (*envConfig, error) {
	cfg := &envConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseCmd() (*cmdConfig, error) {
	cfg := &cmdConfig{}

	sh := new(ServHost)
	_ = flag.Value(sh)

	flag.StringVar(&cfg.PostgresDNS, "d", "", "Postgres DSN")
	flag.StringVar(&cfg.ShortLinkPrefix, "b", "", "short link server")
	flag.StringVar(&cfg.StorageFilePaths, "f", "", "Path to storage file")
	flag.Var(sh, "a", "Net address host:port")
	flag.Parse()

	cfg.Addr = sh.String()
	return cfg, nil
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

func parseDSN(dsn string) (*url.URL, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	return parsedURL, nil
}
