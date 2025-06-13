package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// JSONConfig represents the structure of the JSON configuration file.
type JSONConfig struct {
	StorageFilePaths string `json:"file_storage_path,omitempty"`
	Addr             string `json:"server_address,omitempty"`
	ShortLinkPrefix  string `json:"base_url,omitempty"`
	PostgresDSN      string `json:"database_dsn,omitempty"`
	HTTPSEnable      bool   `json:"enable_https,omitempty"`
}

// parseJSON reads and parses the JSON configuration file from the given directory.
func parseJSON(dirPath string, confFile string) (*JSONConfig, error) {
	confPath := filepath.Join(dirPath, confFile)
	jsonData, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	cfg := &JSONConfig{}
	err = json.Unmarshal(jsonData, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
