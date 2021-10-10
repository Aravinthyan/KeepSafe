package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aravinthyan/KeepSafe/config"
)

// setAppTheme will set the app theme to the selectedTheme (which is currently dark or light)
// and it will refresh all the other tabs (search, add, remove).
func setAppTheme(searchTab, addTab, removeTab fyne.CanvasObject, selectedTheme string) {
	if selectedTheme == config.DarkTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	} else if selectedTheme == config.LightTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}

	// tabs need to be refreshed so that the theme can be applied to all the CanvasObjects
	searchTab.Refresh()
	addTab.Refresh()
	removeTab.Refresh()
}

// LoadConfig will load the user config that was read by ReadConfig().
func LoadConfig(searchTab, addTab, removeTab fyne.CanvasObject, cfg *config.Config) {
	setAppTheme(searchTab, addTab, removeTab, cfg.Theme)
}

// Settings will implement the settings UI so that the user can choose the configs/settings that they
// want for the app.
func Settings(searchTab, addTab, removeTab fyne.CanvasObject, cfg *config.Config) fyne.CanvasObject {
	themeText := widget.NewLabel("Theme")
	themeText.TextStyle.Monospace = true

	themeSelect := widget.NewSelect([]string{config.DarkTheme, config.LightTheme}, func(selectedTheme string) {
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
