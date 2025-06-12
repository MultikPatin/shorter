package config

import (
	"github.com/caarlos0/env/v6"
)

// envConfig holds configuration settings retrieved from environment variables.
type envConfig struct {
	StorageFilePaths string `env:"FILE_STORAGE_PATH"` // File storage paths specified via an environment variable.
	Addr             string `env:"SERVER_ADDRESS"`    // Server address defined by an environment variable.
	ShortLinkPrefix  string `env:"BASE_URL"`          // Short link base URL configured via an environment variable.
	PostgresDSN      string `env:"DATABASE_DSN"`      // PostgreSQL Data Source Name received from an environment variable.
	HTTPSEnable      string `env:"ENABLE_HTTPS"`      // Indicates whether HTTPS is enabled for the server.
	ConfFile         string `env:"CONFIG"`            // Name of the configuration file.
}

// parseEnv extracts configuration from environment variables.
func parseEnv() (*envConfig, error) {
	cfg := &envConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
