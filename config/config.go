package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Theme string
}

const (
	DarkTheme  = "dark"
	LightTheme = "light"
)

var AppConfig = newConfig()

var configFile = "./.config"

func newConfig() *Config {
	cfg := new(Config)
	cfg.Theme = ""
	return cfg
}

func (cfg *Config) ReadConfig() {
	if data, err := ioutil.ReadFile(configFile); err == nil {
		if json.Unmarshal(data, &cfg) == nil {
			return
		}
	}

	// if file does not exist or if unmarshal was unsuccesful then set default config
	cfg.Theme = LightTheme
}

func (cfg *Config) WriteConfig() {
	if jsonData, err := json.Marshal(cfg); err == nil {
		if err = ioutil.WriteFile(configFile, jsonData, 0666); err == nil {
			return
		}
	}
}
