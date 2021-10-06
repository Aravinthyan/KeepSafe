package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Aravinthyan/KeepSafe/data"
	"github.com/Aravinthyan/KeepSafe/database"
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
	passwords := database.New()
	searchData := data.NewListingData()
	addData := data.NewListingData()
	removeData := data.NewListingData()
	window := keepSafe.NewWindow("Keep Safe")
	window.Resize(fyne.NewSize(800, 400))

	var (
		masterPassword []byte // holds master password
		passwordUI     *fyne.Container
		tabs           *container.AppTabs
	)

	// if the passwords file does not exist then a new password file will be created
	// and the UI should request the user to enter a new master password
	if _, err := os.Stat(database.PasswordFile); os.IsNotExist(err) {
		passwordEntryOne := widget.NewPasswordEntry()
		passwordEntryOne.SetPlaceHolder("Enter master password...")

		passwordEntryTwo := widget.NewPasswordEntry()
		passwordEntryTwo.SetPlaceHolder("Enter master password again...")

		errorMsg := canvas.NewText("", data.Red)

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

		incorrectPassword := canvas.NewText("Incorrect password", data.Red)
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

	tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.SearchIcon(), data.Search(searchData, passwords)),
		container.NewTabItemWithIcon("", theme.ContentAddIcon(), data.Add(addData, searchData, removeData, passwords)),
		container.NewTabItemWithIcon("", theme.ContentRemoveIcon(), data.Remove(removeData, searchData, addData, passwords)),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), data.Settings()),
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

	window.SetContent(content)
	window.ShowAndRun()
	passwords.WritePasswords(masterPassword)
}
