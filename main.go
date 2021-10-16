//go:generate fyne bundle -o icon.go Icon.png

package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/Aravinthyan/KeepSafe/config"
	"github.com/Aravinthyan/KeepSafe/database"
	"github.com/Aravinthyan/KeepSafe/frontend"
)

func main() {
	keepSafe := app.New()
	keepSafe.SetIcon(resourceIconPng)
	passwords := database.New()
	userConfig := config.NewConfig()
	userConfig.ReadConfig()

	masterPassword := frontend.LoadUI(keepSafe, passwords, userConfig)

	keepSafe.Run()

	// write data back to files
	passwords.WritePasswords(*masterPassword)
	userConfig.WriteConfig()
}
