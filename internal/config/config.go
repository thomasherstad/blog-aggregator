package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't get home directory. Error:", err)
		return "", err
	}

	filePath := homeDir + "/" + configFileName
	return filePath, nil
}

func Read() (Config, error) {

	var cfg Config
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Coultn't get config filepath. Error:", err)
		return Config{}, err
	}

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Couldn't read config file. Error:", err)
		return Config{}, err
	}

	err = json.Unmarshal(configData, &cfg)
	if err != nil {
		fmt.Println("Couldn't unmarshal config file", err)
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username

	err := write(c)
	if err != nil {
		fmt.Println("Error writing to file. Error:", err)
		return err
	}

	return nil
}
