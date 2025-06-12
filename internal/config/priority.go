package config

import (
	"fmt"
	"net/url"
	"strings"
)

func mergeConfigs(exeDir string, envCfg *envConfig, cmdCfg *cmdConfig, jsonCfg *JSONConfig) (*Config, error) {
	var finalConfig Config

	dsn := ""
	if envCfg.PostgresDSN != "" {
		dsn = envCfg.PostgresDSN
	} else if cmdCfg.PostgresDSN != "" {
		dsn = cmdCfg.PostgresDSN
	} else if jsonCfg.PostgresDSN != "" {
		dsn = jsonCfg.PostgresDSN
	}
	parsedDsn, _ := parseDSN(dsn)
	finalConfig.PostgresDSN = parsedDsn

	if envCfg.Addr != "" {
		finalConfig.Addr = envCfg.Addr
	} else if cmdCfg.Addr != "" {
		finalConfig.Addr = cmdCfg.Addr
	} else {
		finalConfig.Addr = jsonCfg.Addr
	}

	if envCfg.ShortLinkPrefix != "" {
		finalConfig.ShortLinkPrefix = envCfg.ShortLinkPrefix
	} else if cmdCfg.ShortLinkPrefix != "" {
		finalConfig.ShortLinkPrefix = cmdCfg.ShortLinkPrefix
	} else {
		finalConfig.ShortLinkPrefix = jsonCfg.ShortLinkPrefix
	}

	if envCfg.StorageFilePaths != "" {
		finalConfig.StorageFilePaths = envCfg.StorageFilePaths
	} else if cmdCfg.StorageFilePaths != "" {
		finalConfig.StorageFilePaths = cmdCfg.StorageFilePaths
	} else {
		finalConfig.StorageFilePaths = jsonCfg.StorageFilePaths
	}

	httpEnabledStr := ""
	if envCfg.HTTPSEnable != "" {
		httpEnabledStr = envCfg.HTTPSEnable
	} else if cmdCfg.HTTPSEnable != "" {
		httpEnabledStr = cmdCfg.HTTPSEnable
	} else {
		httpEnabledStr = fmt.Sprintf("%t", jsonCfg.HTTPSEnable)
	}
	finalConfig.HTTPSEnable = resolveBool(httpEnabledStr)

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
