package main

import (
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

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter password...")

	incorrectPassword := canvas.NewText("Incorrect password", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	incorrectPassword.Hide()

	enterButton := widget.NewButton("Enter", func() {
		if checkPassword(passwords, []byte(passwordEntry.Text)) {
		} else {
			incorrectPassword.Show()
		}
	})

	content := container.NewGridWithRows(
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

	window.SetContent(content)
	window.ShowAndRun()
}
