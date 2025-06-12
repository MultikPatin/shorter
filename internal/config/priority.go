package config

import (
	"fmt"
	"net/url"
	"strings"
)

// mergeConfigs merges configuration values with priority: environment > command line > JSON file.
// It returns a final Config object and any error that occurred during the merge.
// The priority order is:
// 1. Environment variables (envCfg)
// 2. Command line flags (cmdCfg)
// 3. JSON configuration file (jsonCfg)
func mergeConfigs(exeDir string, envCfg *envConfig, cmdCfg *cmdConfig, jsonCfg *JSONConfig) (*Config, error) {
	var finalConfig Config

	if envCfg.PostgresDSN != "" {
		finalConfig.PostgresDSN, _ = parseDSN(envCfg.PostgresDSN)
	} else if cmdCfg.PostgresDSN != "" {
		finalConfig.PostgresDSN, _ = parseDSN(cmdCfg.PostgresDSN)
	} else if jsonCfg.PostgresDSN != "" {
		finalConfig.PostgresDSN, _ = parseDSN(jsonCfg.PostgresDSN)
	}

	if envCfg.Addr != "" {
		finalConfig.Addr = envCfg.Addr
	} else if cmdCfg.Addr != "" {
		finalConfig.Addr = cmdCfg.Addr
	} else if jsonCfg.Addr != "" {
		finalConfig.Addr = jsonCfg.Addr
	}

	if envCfg.ShortLinkPrefix != "" {
		finalConfig.ShortLinkPrefix = envCfg.ShortLinkPrefix
	} else if cmdCfg.ShortLinkPrefix != "" {
		finalConfig.ShortLinkPrefix = cmdCfg.ShortLinkPrefix
	} else if jsonCfg.ShortLinkPrefix != "" {
		finalConfig.ShortLinkPrefix = jsonCfg.ShortLinkPrefix
	}

	if envCfg.StorageFilePaths != "" {
		finalConfig.StorageFilePaths = envCfg.StorageFilePaths
	} else if cmdCfg.StorageFilePaths != "" {
		finalConfig.StorageFilePaths = cmdCfg.StorageFilePaths
	} else if jsonCfg.StorageFilePaths != "" {
		finalConfig.StorageFilePaths = jsonCfg.StorageFilePaths
	}

	if envCfg.HTTPSEnable != "" {
		finalConfig.HTTPSEnable = resolveBool(envCfg.HTTPSEnable)
	} else if cmdCfg.HTTPSEnable != "" {
		finalConfig.HTTPSEnable = resolveBool(cmdCfg.HTTPSEnable)
	} else if jsonCfg.HTTPSEnable != false {
		finalConfig.HTTPSEnable = jsonCfg.HTTPSEnable
	}

	finalConfig.PProfAddr = defaultPProfAddr
	finalConfig.ExecutableDir = exeDir

	if finalConfig.StorageFilePaths == "" {
		finalConfig.StorageFilePaths = defaultStorageFilePath
	}

	return &finalConfig, nil
}

// parseDSN converts a raw Data Source Name (DSN) string into a structured URL object.
func parseDSN(dsn string) (*url.URL, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	return u, nil
}

// resolveBool converts a string representation of a boolean into a bool value.
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
