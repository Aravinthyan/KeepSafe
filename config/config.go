package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config contains the information about the user configuartion for the app.
type Config struct {
	Theme string
}

const (
	DarkTheme  = "dark"
	LightTheme = "light"
)

var configFile = "/appConfig"

// NewConfig creates a new Config.
func NewConfig() *Config {
	return new(Config)
}

// ReadConfig will read the configFile and populate the user configuration into the struct and if
// the config files does not exist or if there is an error when unmarshal then set default values.
func (cfg *Config) ReadConfig() {
	exeFilePath, _ := os.Executable()
	exeDirPath := filepath.Dir(exeFilePath)
	data, err := ioutil.ReadFile(exeDirPath + configFile)
	if err == nil && json.Unmarshal(data, &cfg) == nil {
		return
	}

	// if file does not exist or if unmarshal was unsuccesful then set default config
	cfg.Theme = LightTheme
}

// WriteConfig will write the user config to the configFile.
func (cfg *Config) WriteConfig() {
	jsonData, _ := json.Marshal(cfg)
	exeFilePath, _ := os.Executable()
	exeDirPath := filepath.Dir(exeFilePath)
	ioutil.WriteFile(exeDirPath+configFile, jsonData, 0600)
}
