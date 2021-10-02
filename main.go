package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"image/color"

	"github.com/Aravinthyan/KeepSafe/database"
)

// checkPassword will check if the user provided password is valid to decrypt the password database
func checkPassword(passwords *database.PasswordDB, password []byte) bool {
	if err := passwords.ReadPasswords(password); err != nil {
		return false
	}
	return true
}

func main() {
	keepSafe := app.New()
	passwords := database.New()

	window := keepSafe.NewWindow("Keep Safe")
	window.Resize(fyne.NewSize(500, 700))

	var masterPassword []byte
	var content *fyne.Container

	if _, err := os.Stat(database.PasswordFile); os.IsNotExist(err) {
		passwordEntryOne := widget.NewPasswordEntry()
		passwordEntryOne.SetPlaceHolder("Enter master password...")

		passwordEntryTwo := widget.NewPasswordEntry()
		passwordEntryTwo.SetPlaceHolder("Enter master password again...")

		passwordsDoNotMatch := canvas.NewText("Passwords do not match", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		passwordsDoNotMatch.Hide()

		enterButton := widget.NewButton("Enter", func() {
			if passwordEntryOne.Text == passwordEntryTwo.Text {
				masterPassword = []byte(passwordEntryOne.Text)
				passwords.CreateEmptyDB()
			} else {
				passwordsDoNotMatch.Show()
			}
		})

		content = container.NewGridWithRows(
			12,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			passwordEntryOne,
			passwordEntryTwo,
			enterButton,
			passwordsDoNotMatch,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
		)

	} else {

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Enter master password...")

		incorrectPassword := canvas.NewText("Incorrect password", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		incorrectPassword.Hide()

		enterButton := widget.NewButton("Enter", func() {
			masterPassword = []byte(passwordEntry.Text)
			if checkPassword(passwords, masterPassword) {
			} else {
				incorrectPassword.Show()
			}
		})

		content = container.NewGridWithRows(
			11,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			passwordEntry,
			enterButton,
			incorrectPassword,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
		)
	}

	window.SetContent(content)
	window.ShowAndRun()
	passwords.WritePasswords(masterPassword)
}
