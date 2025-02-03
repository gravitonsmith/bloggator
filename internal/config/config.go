package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUser = username
	return write(*cfg)
}

func write(cfg Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, configFileName)
	return path, nil
}
