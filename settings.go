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

func loadConfig() string {
	setAppTheme(config.AppConfig.Theme)
	currentView := config.AppConfig.View
	return currentView
}

func settings(passwords *database.PasswordDB, window fyne.Window) {
	homeText := canvas.NewText("Default home view", theme.ForegroundColor())
	homeSelect := widget.NewSelect([]string{config.SerachView, config.ListView}, func(selectedDefaultView string) {
		config.AppConfig.View = selectedDefaultView
	})
	homeSelect.PlaceHolder = config.AppConfig.View

	themeText := canvas.NewText("Theme", theme.ForegroundColor())
	themeSelect := widget.NewSelect([]string{config.DarkTheme, config.LightTheme}, func(selectedTheme string) {
		setAppTheme(selectedTheme)
		themeText.Color = theme.ForegroundColor()
		homeText.Color = theme.ForegroundColor()
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
		container.NewVBox(
			themeText,
			themeSelect,
			homeText,
			homeSelect,
		),
	)
	window.SetContent(content)
}
