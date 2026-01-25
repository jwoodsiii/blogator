package config

import (
	"os"
	"fmt"
	"encoding/json"
)

const configFileName = "/.blogatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	// read blogatorconfig.json from home directory and marshal into Config struct
	fp, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(fp)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
	    fmt.Println("error decoding config")
	    return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser() {
	// write config struct to json after setting current_user_name
	cfg.Current_user_name = "jwoodsiii"
}

func getConfigFilePath() (string, error) {
	// return filepath to blogator config filez
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + configFileName, nil
}

func write(cfg Config) error {
	// write config struct to json file at config file
	fp, err := getConfigFilePath()
	if err != nil {
		return err
	}
	cfg_json, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(fp, cfg_json, 0666); err != nil {
		return err
	}
	return nil
}
