package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"image/color"

	"github.com/Aravinthyan/KeepSafe/database"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// checkPassword will check if the user provided password is valid to decrypt the password database
func checkPassword(passwords *database.PasswordDB, password []byte) bool {
	if err := passwords.ReadPasswords(password); err != nil {
		return false
	}
	return true
}

func showServicesAsList(services []string, passwords *database.PasswordDB, window fyne.Window) {
	data := binding.BindStringList(&services)

	servicesList := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	servicesList.OnSelected = func(id int) {
		password := passwords.Password(services[id])

		serviceLiteralText := canvas.NewText("Service", color.White)
		passwordLiteralText := canvas.NewText("Password", color.White)

		serviceText := canvas.NewText(services[id], color.White)
		passwordText := canvas.NewText(password, color.White)

		serviceText.TextStyle.Bold = true
		passwordText.TextStyle.Bold = true

		serviceLiteralText.TextSize = 35
		passwordLiteralText.TextSize = 35

		serviceText.TextSize = 35
		passwordText.TextSize = 35

		backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			showServicesAsList(services, passwords, window)
		})

		content := container.NewBorder(
			container.New(layout.NewHBoxLayout(), backButton, layout.NewSpacer()),
			nil,
			nil,
			nil,
			container.NewCenter(
				container.NewVBox(
					serviceLiteralText,
					serviceText,
					passwordLiteralText,
					passwordText,
				),
			),
		)

		window.SetContent(content)
	}

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		accessPasswords(passwords, window)
	})

	content := container.NewBorder(
		container.New(layout.NewHBoxLayout(), backButton, layout.NewSpacer()),
		nil,
		nil,
		nil,
		servicesList,
	)

	window.SetContent(content)
}

func add(passwords *database.PasswordDB, window fyne.Window) {
	serviceEntry := widget.NewEntry()
	serviceEntry.SetPlaceHolder("Enter service...")

	passwordEntryOne := widget.NewPasswordEntry()
	passwordEntryOne.SetPlaceHolder("Enter password...")

	passwordEntryTwo := widget.NewPasswordEntry()
	passwordEntryTwo.SetPlaceHolder("Enter password again...")

	errorMsg := canvas.NewText("", color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	createButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		if serviceEntry.Text == "" || passwordEntryOne.Text == "" || passwordEntryTwo.Text == "" {
			errorMsg.Text = "Please fill out all fields"
			errorMsg.Refresh()
			return
		}

		if passwordEntryOne.Text != passwordEntryTwo.Text {
			errorMsg.Text = "Passwords do not match"
			errorMsg.Refresh()
			return
		}

		if err := passwords.Insert(serviceEntry.Text, passwordEntryOne.Text); err != nil {
			errorMsg.Text = "Service already exists"
			errorMsg.Refresh()
		} else {
			accessPasswords(passwords, window)
		}
	})

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		accessPasswords(passwords, window)
	})

	content := container.NewGridWithRows(
		16,
		container.New(layout.NewHBoxLayout(), backButton, layout.NewSpacer()),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		serviceEntry,
		passwordEntryOne,
		passwordEntryTwo,
		createButton,
		errorMsg,
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
	)

	window.SetContent(content)
}

func delete(passwords *database.PasswordDB, window fyne.Window) {
	serviceEntry := widget.NewEntry()
	serviceEntry.SetPlaceHolder("Enter service...")

	errorMsg := canvas.NewText("", color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	deleteButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		if serviceEntry.Text == "" {
			errorMsg.Text = "Please enter service"
			errorMsg.Refresh()
			return
		}

		if err := passwords.Remove(serviceEntry.Text); err != nil {
			errorMsg.Text = "Service does not exists"
			errorMsg.Refresh()
		} else {
			accessPasswords(passwords, window)
		}
	})

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		accessPasswords(passwords, window)
	})

	content := container.NewGridWithRows(
		16,
		container.New(layout.NewHBoxLayout(), backButton, layout.NewSpacer()),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		serviceEntry,
		deleteButton,
		errorMsg,
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
	)

	window.SetContent(content)
}

func accessPasswords(passwords *database.PasswordDB, window fyne.Window) {
	serviceEntry := widget.NewEntry()
	serviceEntry.SetPlaceHolder("Enter service...")

	errorMsg := canvas.NewText("", color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	searchButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		if serviceEntry.Text == "" {
			errorMsg.Text = "Please enter a service"
			errorMsg.Refresh()
			return
		}
		searchResult := fuzzy.Find(serviceEntry.Text, passwords.Services)
		if searchResult == nil {
			errorMsg.Text = "Service does not exist"
			errorMsg.Refresh()
			return
		}
		showServicesAsList(searchResult, passwords, window)
	})

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		add(passwords, window)
	})

	deleteButton := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		delete(passwords, window)
	})

	listButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		showServicesAsList(passwords.Services, passwords, window)
	})

	content := container.NewBorder(
		nil,
		container.NewGridWithColumns(
			2,
			addButton,
			deleteButton,
		),
		nil,
		nil,
		container.NewGridWithRows(
			15,
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), listButton),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			serviceEntry,
			searchButton,
			errorMsg,
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
		),
	)

	window.SetContent(content)
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
				accessPasswords(passwords, window)
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
				accessPasswords(passwords, window)
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
