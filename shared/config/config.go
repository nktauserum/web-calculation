package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Port int
}

func GetConfig() (*Config, error) {
	path := "shared/config/config.json"

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	info, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = json.Unmarshal(info, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
