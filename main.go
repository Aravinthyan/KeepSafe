//go:generate fyne bundle -o icon.go Icon.png

package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Aravinthyan/KeepSafe/config"
	"github.com/Aravinthyan/KeepSafe/database"
	"github.com/Aravinthyan/KeepSafe/frontend"
)

// checkPassword will check if the user provided password is valid to decrypt the password database.
func checkPassword(passwords *database.PasswordDB, password []byte) bool {
	if err := passwords.ReadPasswords(password); err != nil {
		return false
	}
	return true
}

func main() {
	keepSafe := app.New()
	keepSafe.SetIcon(resourceIconPng)
	passwords := database.New()
	searchData := frontend.NewListingData()
	addData := frontend.NewListingData()
	removeData := frontend.NewListingData()
	window := keepSafe.NewWindow("Keep Safe")
	window.Resize(fyne.NewSize(800, 400))
	window.SetFixedSize(true)

	userConfig := config.NewConfig()
	userConfig.ReadConfig()

	var (
		masterPassword []byte // holds master password
		passwordUI     *fyne.Container
		tabs           *container.AppTabs
	)

	exeFilePath, _ := os.Executable()
	exeDirPath := filepath.Dir(exeFilePath)

	// if the passwords file does not exist then a new password file will be created
	// and the UI should request the user to enter a new master password
	if _, err := os.Stat(exeDirPath + database.PasswordFile); os.IsNotExist(err) {
		passwordEntryOne := widget.NewPasswordEntry()
		passwordEntryOne.SetPlaceHolder("Enter master password...")

		passwordEntryTwo := widget.NewPasswordEntry()
		passwordEntryTwo.SetPlaceHolder("Enter master password again...")

		errorMsg := canvas.NewText("", frontend.Red)

		enterButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
			masterPassword = []byte(passwordEntryOne.Text)
			passwords.CreateEmptyDB()
			passwordUI.Hide()
			tabs.Show()
		})
		enterButton.Disable()

		onChanged := func(password string) {
			if passwordEntryTwo.Text == "" {
				errorMsg.Text = ""
				enterButton.Disable()
			} else if passwordEntryOne.Text != passwordEntryTwo.Text {
				errorMsg.Text = "Passwords do not match"
				enterButton.Disable()
			} else {
				errorMsg.Text = ""
				enterButton.Enable()
			}
			errorMsg.Refresh()
		}

		passwordEntryOne.OnChanged = onChanged
		passwordEntryTwo.OnChanged = onChanged

		passwordUI = container.NewVBox(
			passwordEntryOne,
			passwordEntryTwo,
			enterButton,
			errorMsg,
		)
	} else { // password file exist, therefore ask to enter master password
		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Enter master password...")

		incorrectPassword := canvas.NewText("Incorrect password", frontend.Red)
		incorrectPassword.Hide()

		enterButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
			masterPassword = []byte(passwordEntry.Text)
			if checkPassword(passwords, masterPassword) {
				searchData.SearchResult = passwords.Services
				searchData.Services.Reload()
				addData.SearchResult = passwords.Services
				addData.Services.Reload()
				removeData.SearchResult = passwords.Services
				removeData.Services.Reload()
				passwordUI.Hide()
				tabs.Show()
			} else {
				incorrectPassword.Show()
			}
		})

		passwordUI = container.NewVBox(
			passwordEntry,
			enterButton,
			incorrectPassword,
		)
	}

	searchTab := frontend.Search(searchData, passwords)
	addTab := frontend.Add(addData, searchData, removeData, passwords)
	removeTab := frontend.Remove(removeData, searchData, addData, passwords)

	tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.SearchIcon(), searchTab),
		container.NewTabItemWithIcon("", theme.ContentAddIcon(), addTab),
		container.NewTabItemWithIcon("", theme.ContentRemoveIcon(), removeTab),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), frontend.Settings(searchTab, addTab, removeTab, userConfig)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.Hide()

	content := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		passwordUI,
		tabs,
	)

	frontend.LoadConfig(searchTab, addTab, removeTab, userConfig)
	window.SetContent(content)
	window.ShowAndRun()
	passwords.WritePasswords(masterPassword)
	userConfig.WriteConfig()
}
