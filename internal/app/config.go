package app

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultStorageFilePath = "shorter.json"
)

var sugar zap.SugaredLogger

type Config struct {
	Addr             string
	ShortLinkPrefix  string
	StorageFilePaths string
}

type envConfig struct {
	StorageFilePaths string `env:"FILE_STORAGE_PATH"`
	Addr             string `env:"SERVER_ADDRESS"`
	ShortLinkPrefix  string `env:"BASE_URL"`
}
type cmdConfig struct {
	Addr             string
	StorageFilePaths string
	ShortLinkPrefix  string
}

type ServHost struct {
	Host string
	Port int
}

func ParseConfig() (*Config, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	cfg := &Config{}

	envCfg, err := parseEnv()
	if err != nil {
		sugar.Infow(
			"Parsed Env",
			"error", err.Error(),
		)
	}
	cmdCfg, err := parseCmd()
	if err != nil {
		sugar.Infow(
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
		cfg.StorageFilePaths = filepath.Join(cmdCfg.StorageFilePaths, defaultStorageFilePath)
	} else {
		cfg.StorageFilePaths = filepath.Join(envCfg.StorageFilePaths, defaultStorageFilePath)
	}
	return cfg, nil
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
