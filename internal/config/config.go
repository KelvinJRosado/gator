package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Get the full path to the config file
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := filepath.Join(homeDir, configFileName)

	return configFilePath, nil
}

func write(c Config) error {
	// Convert from struct to raw JSON
	data, err := json.Marshal(&c)
	if err != nil {
		return err
	}

	// Get the config file path
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Write to the file
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Read in the config file as a usable format
func Read() (*Config, error) {
	// Get the config file path
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	// Read the file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	// Convert from raw json to struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}

func (c *Config) SetUser(name string) error {

	// Check base case
	if name == c.CurrentUserName {
		return nil
	}

	// Update user in struct
	c.CurrentUserName = name

	// Update config file
	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
