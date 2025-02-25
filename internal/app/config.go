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
	defaultStorageFilePath = "shorter"
)

var sugar zap.SugaredLogger

type Config struct {
	Addr             string
	ShortLinkPrefix  string
	StorageFilePaths string
}

type envConfig struct {
	StorageFilePaths string `env:"FILE_STORAGE_PATH"`
	Addr             string `env:"SERVER_ADDRESS,required"`
	ShortLinkPrefix  string `env:"BASE_URL"`
}
type cmdConfig struct {
	ServHost    ServHost
	ShorLink    ShorLink
	FileStorage FileStorage
}
type ServHost struct {
	Host string
	Port int
}
type ShorLink struct {
	ShortLinkPrefix string
}
type FileStorage struct {
	FilePath string
}

func ParseConfig() (*Config, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	cfg := &Config{}
	if err := cfg.parseEnv(); err != nil {
		if err := cfg.parseFlags(); err != nil {
			return nil, err
		}
	}
	cfg.StorageFilePaths = filepath.Join(cfg.StorageFilePaths, defaultStorageFilePath)
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
	c.StorageFilePaths = cfg.StorageFilePaths
	sugar.Info(
		"Parsed Env",
		c,
	)
	return nil
}

func (c *Config) parseFlags() error {
	sv := new(ServHost)
	_ = flag.Value(sv)
	sh := new(ShorLink)
	_ = flag.Value(sh)
	fs := new(FileStorage)
	_ = flag.Value(fs)

	flag.Var(sv, "a", "Net address host:port")
	flag.Var(sh, "b", "short link server")
	flag.Var(fs, "f", "Path to storage file")
	flag.Parse()

	c.Addr = sv.String()
	c.ShortLinkPrefix = sh.String()
	c.StorageFilePaths = fs.String()
	sugar.Info(
		"Parsed Flags",
		c,
	)
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

func (a *FileStorage) String() string {
	return a.FilePath
}

func (a *FileStorage) Set(s string) error {
	hp := strings.Split(s, ":")
	a.FilePath = hp[0]
	return nil
}
