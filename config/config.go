package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

// setAppTheme will set the app theme to the selectedTheme (which is currently dark or light)
// and it will refresh all the other tabs (search, add, remove).
func setAppTheme(searchTab, addTab, removeTab fyne.CanvasObject, selectedTheme string) {
	if selectedTheme == DarkTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	} else if selectedTheme == LightTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}

	// tabs need to be refreshed so that the theme can be applied to all the CanvasObjects
	searchTab.Refresh()
	addTab.Refresh()
	removeTab.Refresh()
}

// LoadConfig will load the user config that was read by ReadConfig().
func (cfg *Config) LoadConfig(searchTab, addTab, removeTab fyne.CanvasObject) {
	setAppTheme(searchTab, addTab, removeTab, cfg.Theme)
}

// Settings will implement the settings UI so that the user can choose the configs/settings that they
// want for the app.
func (cfg *Config) Settings(searchTab, addTab, removeTab fyne.CanvasObject) fyne.CanvasObject {
	themeText := widget.NewLabel("Theme")
	themeText.TextStyle.Monospace = true

	themeSelect := widget.NewSelect([]string{DarkTheme, LightTheme}, func(selectedTheme string) {
		setAppTheme(searchTab, addTab, removeTab, selectedTheme)
		cfg.Theme = selectedTheme
	})
	themeSelect.PlaceHolder = cfg.Theme

	content := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewVBox(
			themeText,
			themeSelect,
		),
	)

	return content
}
