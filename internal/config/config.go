package config // Package config handles parsing environment variables and command-line arguments into a unified configuration structure.

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
	defaultStorageFilePath = "shorter" // Default path for storage file if no custom path is provided.
)

// Config stores all the necessary configurations from both environment variables and command line inputs.
type Config struct {
	PostgresDSN      *url.URL // Database connection details (Data Source Name).
	PProfAddr        string   // Address for pprof profiling endpoint.
	Addr             string   // Server listening address.
	ShortLinkPrefix  string   // Base URL for short links.
	StorageFilePaths string   // Path where storage files are located.
	HttpsEnable      bool     // Indicates whether HTTPS is enabled for the server.
}

// envConfig holds configuration settings retrieved from environment variables.
type envConfig struct {
	StorageFilePaths string `env:"FILE_STORAGE_PATH"` // File storage paths specified via an environment variable.
	Addr             string `env:"SERVER_ADDRESS"`    // Server address defined by an environment variable.
	ShortLinkPrefix  string `env:"BASE_URL"`          // Short link base URL configured via an environment variable.
	PostgresDSN      string `env:"DATABASE_DSN"`      // PostgreSQL Data Source Name received from an environment variable.
	HttpsEnable      string `env:"ENABLE_HTTPS"`      // Indicates whether HTTPS is enabled for the server.
}

// cmdConfig holds configuration settings obtained from command-line flags.
type cmdConfig struct {
	Addr             string // Command-line argument for server address.
	StorageFilePaths string // Command-line option specifying file storage paths.
	ShortLinkPrefix  string // Base URL for short links passed via command-line.
	PostgresDSN      string // Postgres DSN given on the command line.
	HttpsEnable      string // Indicates whether HTTPS is enabled for the server.
}

// servHost encapsulates information about the network service's host and port.
type servHost struct {
	Host string // Hostname part of the net address.
	Port int    // Port number component of the net address.
}

// Parse merges environment variables and command-line options into a single configuration object.
func Parse(logger *zap.SugaredLogger) *Config {
	cfg := &Config{}

	envCfg, err := parseEnv()
	if err != nil {
		logger.Infow("Error while parsing environment variables", "error", err.Error())
	}

	cmdCfg, err := parseCmd()
	if err != nil {
		logger.Infow("Error while parsing command-line arguments", "error", err.Error())
	}

	if envCfg.Addr == "" {
		cfg.Addr = cmdCfg.Addr
	} else {
		cfg.Addr = envCfg.Addr
	}
	if envCfg.HttpsEnable == "" {
		cfg.HttpsEnable = resolveBool(cmdCfg.HttpsEnable)
	} else {
		cfg.HttpsEnable = resolveBool(envCfg.HttpsEnable)
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
	if envCfg.PostgresDSN != "" {
		cfg.PostgresDSN, _ = parseDSN(envCfg.PostgresDSN)
	} else if cmdCfg.PostgresDSN != "" {
		cfg.PostgresDSN, _ = parseDSN(cmdCfg.PostgresDSN)
	}
	cfg.PProfAddr = "localhost:6060"

	return cfg
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

// parseCmd processes command-line flags to extract configuration values.
func parseCmd() (*cmdConfig, error) {
	cfg := &cmdConfig{}

	hostPort := new(servHost)
	_ = flag.Value(hostPort)

	flag.StringVar(&cfg.PostgresDSN, "d", "", "Postgres DSN")
	flag.StringVar(&cfg.ShortLinkPrefix, "b", "", "Short link server")
	flag.StringVar(&cfg.StorageFilePaths, "f", "", "Path to storage file")
	flag.StringVar(&cfg.HttpsEnable, "s", "0", "HTTPS is enabled")
	flag.Var(hostPort, "a", "Network address host:port")
	flag.Parse()

	cfg.Addr = hostPort.String()
	return cfg, nil
}

// String returns the formatted representation of the ServHost instance.
func (a *servHost) String() string {
	a.normalize()
	return a.Host + ":" + strconv.Itoa(a.Port)
}

// Set parses a string input representing a host-port pair and updates the ServHost fields accordingly.
func (a *servHost) Set(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return errors.New("address must be in format host:port")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	a.Host = parts[0]
	a.Port = port
	a.normalize()
	return nil
}

// normalize ensures valid defaults for empty or zero-value fields.
func (a *servHost) normalize() {
	if a.Port == 0 {
		a.Port = 8080
	}
	if a.Host == "" {
		a.Host = "localhost"
	}
}

// parseDSN converts a raw Data Source Name (DSN) string into a structured URL object.
func parseDSN(dsn string) (*url.URL, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	return u, nil
}

func resolveBool(arg string) bool {
	switch strings.ToLower(arg) {
	case "true", "yes", "1":
		return true
	case "false", "no", "0":
		return false
	default:
		return false
	}
}
