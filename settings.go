package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aravinthyan/KeepSafe/config"
	"github.com/Aravinthyan/KeepSafe/database"
)

func setAppTheme(selectedTheme string) {
	if selectedTheme == config.DarkTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	} else if selectedTheme == config.LightTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}
}

func loadConfig() {
	setAppTheme(config.AppConfig.Theme)
}

func settings(passwords *database.PasswordDB, window fyne.Window) {
	themeText := canvas.NewText("Theme", theme.ForegroundColor())
	themeSelect := widget.NewSelect([]string{config.DarkTheme, config.LightTheme}, func(selectedTheme string) {
		setAppTheme(selectedTheme)
		themeText.Color = theme.ForegroundColor()
		config.AppConfig.Theme = selectedTheme
	})
	themeSelect.PlaceHolder = config.AppConfig.Theme

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		accessPasswords(passwords, window)
	})

	content := container.NewBorder(
		container.NewHBox(backButton, layout.NewSpacer()),
		nil,
		nil,
		nil,
		container.NewVBox(themeText, themeSelect),
	)
	window.SetContent(content)
}
