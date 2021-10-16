package frontend

import (
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aravinthyan/KeepSafe/config"
	"github.com/Aravinthyan/KeepSafe/database"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sethvargo/go-password/password"
)

// red colour is used for error messages so declared once so that it can be used for all cases
var red = color.NRGBA{R: 255, G: 0, B: 0, A: 255}

// ListingData has information about the data that is currently being shown by a list widget.
type ListingData struct {
	SearchResult []string
	Services     binding.ExternalStringList
}

// NewListingData creates a new ListingData and binds Services to SearchResult.
func NewListingData() *ListingData {
	data := new(ListingData)
	data.Services = binding.BindStringList(&data.SearchResult)
	return data
}

// search creates the UI that will allow a user to search and find the desired password for a service.
func search(data *ListingData, passwords *database.PasswordDB) fyne.CanvasObject {
	serviceEntry := widget.NewEntry()
	serviceEntry.SetPlaceHolder("Enter service...")
	serviceEntry.OnChanged = func(service string) {
		data.SearchResult = fuzzy.Find(service, passwords.Services)
		data.Services.Reload()
	}

	servicesList := widget.NewListWithData(
		data.Services,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	serviceLiteralText := widget.NewLabel("Service")
	passwordLiteralText := widget.NewLabel("Password")
	serviceText := widget.NewLabel("")
	passwordText := widget.NewLabel("")

	serviceLiteralText.TextStyle.Monospace = true
	passwordLiteralText.TextStyle.Monospace = true
	serviceText.TextStyle.Bold = true
	passwordText.TextStyle.Bold = true

	servicesList.OnSelected = func(id int) {
		serviceText.SetText(data.SearchResult[id])
		passwordText.SetText(passwords.Password(data.SearchResult[id]))
	}

	infoCombined := container.NewCenter(
		container.NewVBox(
			serviceLiteralText,
			serviceText,
			passwordLiteralText,
			passwordText,
		),
	)

	passwordVisibility := widget.NewSlider(0, 1)
	passwordVisibility.Step = 1
	passwordVisibility.OnChanged = func(currentValue float64) {
		if infoCombined.Hidden {
			infoCombined.Show()
		} else {
			infoCombined.Hide()
		}
	}

	left := container.NewBorder(
		serviceEntry,
		nil,
		nil,
		nil,
		servicesList,
	)

	right := container.NewBorder(
		container.NewHBox(layout.NewSpacer(), passwordVisibility),
		nil,
		nil,
		nil,
		infoCombined,
	)

	return container.NewGridWithColumns(2, left, right)
}

// add creates the UI that will allow a user to add a new service and the corresponding password.
func add(data, searchData, removeData *ListingData, passwords *database.PasswordDB) fyne.CanvasObject {
	servicesList := widget.NewListWithData(
		data.Services,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	right := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		servicesList,
	)

	serviceEntry := widget.NewEntry()
	serviceEntry.SetPlaceHolder("Enter service...")

	serviceEntry.OnChanged = func(service string) {
		data.SearchResult = fuzzy.Find(service, passwords.Services)
		data.Services.Reload()
	}

	displayMsg := canvas.NewText("", red)

	passwordEntryOne := widget.NewPasswordEntry()
	passwordEntryOne.SetPlaceHolder("Enter password...")

	passwordEntryTwo := widget.NewPasswordEntry()
	passwordEntryTwo.SetPlaceHolder("Enter password again...")

	createButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		if err := passwords.Insert(serviceEntry.Text, passwordEntryOne.Text); err != nil {
			displayMsg.Text = "Service already exists"
			displayMsg.Refresh()
			return
		}
		// need to copy the slice to SearchResult, otherwise will use old one which is an error
		searchData.SearchResult = passwords.Services
		// reload to show updated list
		searchData.Services.Reload()
		data.SearchResult = passwords.Services
		data.Services.Reload()
		removeData.SearchResult = passwords.Services
		removeData.Services.Reload()
		serviceEntry.SetText("")
		passwordEntryOne.SetText("")
		passwordEntryTwo.SetText("")
	})
	createButton.Disable()

	onChanged := func(password string) {
		if serviceEntry.Text == "" || passwordEntryTwo.Text == "" {
			displayMsg.Text = ""
			createButton.Disable()
		} else if passwordEntryOne.Text != passwordEntryTwo.Text {
			displayMsg.Text = "Passwords do not match"
			createButton.Disable()
		} else {
			displayMsg.Text = ""
			createButton.Enable()
		}
		displayMsg.Refresh()
	}

	passwordEntryOne.OnChanged = onChanged
	passwordEntryTwo.OnChanged = onChanged

	generateButton := widget.NewButton("Generate Password", func() {
		generatedPassword, err := password.Generate(20, 3, 3, false, false)
		if err != nil {
			displayMsg.Text = "Failed to generate a password"
			displayMsg.Refresh()
		}
		passwordEntryOne.SetText(generatedPassword)
		passwordEntryTwo.SetText(generatedPassword)
	})

	left := container.NewVBox(
		serviceEntry,
		passwordEntryOne,
		passwordEntryTwo,
		createButton,
		generateButton,
		displayMsg,
	)

	return container.NewGridWithColumns(2, left, right)
}

// remove creates the UI that will allow a user to remove an existing service and password.
func remove(data, searchData, addData *ListingData, passwords *database.PasswordDB) fyne.CanvasObject {
	serviceEntry := widget.NewEntry()
	serviceEntry.SetPlaceHolder("Enter service...")
	serviceEntry.OnChanged = func(service string) {
		data.SearchResult = fuzzy.Find(service, passwords.Services)
		data.Services.Reload()
	}

	servicesList := widget.NewListWithData(
		data.Services,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	left := container.NewBorder(
		serviceEntry,
		nil,
		nil,
		nil,
		servicesList,
	)

	prompt := widget.NewLabel("")
	prompt.Hide()
	var serviceToRemove string
	var yesButton *widget.Button

	servicesList.OnSelected = func(id int) {
		prompt.SetText("Are you sure you want to delete " + data.SearchResult[id] + "?")
		serviceToRemove = data.SearchResult[id]
		prompt.Refresh()
		prompt.Show()
		yesButton.Show()
	}

	yesButton = widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		passwords.Remove(serviceToRemove)
		data.SearchResult = passwords.Services
		data.Services.Reload()
		searchData.SearchResult = passwords.Services
		searchData.Services.Reload()
		addData.SearchResult = passwords.Services
		addData.Services.Reload()
		serviceEntry.SetText("")
		prompt.Hide()
		yesButton.Hide()
	})
	yesButton.Hide()

	right := container.NewVBox(
		prompt,
		yesButton,
	)

	return container.NewGridWithColumns(2, left, right)
}

// checkPassword will check if the user provided password is valid to decrypt the password database.
func checkPassword(passwords *database.PasswordDB, password []byte) bool {
	if err := passwords.ReadPasswords(password); err != nil {
		return false
	}
	return true
}

// LoadUI will load the UI only without running the application.
func LoadUI(keepSafe fyne.App, passwords *database.PasswordDB, userConfig *config.Config) *[]byte {
	searchData := NewListingData()
	addData := NewListingData()
	removeData := NewListingData()
	window := keepSafe.NewWindow("Keep Safe")
	window.Resize(fyne.NewSize(800, 400))
	window.SetFixedSize(true)

	var (
		masterPassword []byte
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

		errorMsg := canvas.NewText("", red)

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

		incorrectPassword := canvas.NewText("Incorrect password", red)
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

	searchTab := search(searchData, passwords)
	addTab := add(addData, searchData, removeData, passwords)
	removeTab := remove(removeData, searchData, addData, passwords)

	tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.SearchIcon(), searchTab),
		container.NewTabItemWithIcon("", theme.ContentAddIcon(), addTab),
		container.NewTabItemWithIcon("", theme.ContentRemoveIcon(), removeTab),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), settings(searchTab, addTab, removeTab, userConfig)),
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

	loadConfig(searchTab, addTab, removeTab, userConfig)
	window.SetContent(content)
	window.Show()

	return &masterPassword
}
