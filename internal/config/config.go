package config

import (
	"go.uber.org/zap"
	"net/url"
)

const (
	defaultStorageFilePath = "shorter"        // Default path for storage file if no custom path is provided.
	defaultPProfAddr       = "localhost:6060" // Address for pprof profiling endpoint.
	defaultConfFileName    = "conf.json"      // Name of the configuration file in json format
)

// Config stores all the necessary configurations from both environment variables and command line inputs.
type Config struct {
	PostgresDSN      *url.URL // Database connection details (Data Source Name).
	PProfAddr        string   // Address for pprof profiling endpoint.
	Addr             string   // Server listening address.
	ShortLinkPrefix  string   // Base URL for short links.
	StorageFilePaths string   // Path where storage files are located.
	ExecutableDir    string   // Project directory
	HTTPSEnable      bool     // Indicates whether HTTPS is enabled for the server.
}

// Parse merges environment variables and command-line options into a single configuration object.
func Parse(exeDir string, logger *zap.SugaredLogger) *Config {
	envCfg, err := parseEnv()
	if err != nil {
		logger.Infow("Error while parsing environment variables", "error", err.Error())
	}

	cmdCfg, err := parseCmd()
	if err != nil {
		logger.Infow("Error while parsing command-line arguments", "error", err.Error())
	}

	confDir := "."
	confFileName := defaultConfFileName
	if envCfg.ConfFile != "" {
		confFileName = envCfg.ConfFile
	} else if cmdCfg.ConfFile != "" {
		confFileName = cmdCfg.ConfFile
	}

	jsonCfg, err := parseJSON(confDir, confFileName)
	if err != nil {
		logger.Infow("Error while parsing JSON configuration file", "error", err.Error())
	}

	cfg, err := mergeConfigs(exeDir, envCfg, cmdCfg, jsonCfg)
	if err != nil {
		logger.Infow("Error while merging configurations", "error", err.Error())
	}

	return cfg
}
