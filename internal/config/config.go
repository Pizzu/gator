package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, nil
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configData := Config{}
	err = decoder.Decode(&configData)
	if err != nil {
		return Config{}, err
	}

	return configData, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	// Read the file
	file, err := os.Create(configFilePath)

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

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}
