package frontend

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aravinthyan/KeepSafe/database"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sethvargo/go-password/password"
)

// Red colour is used for error messages so declared once so that it can be used for all cases
var Red = color.NRGBA{R: 255, G: 0, B: 0, A: 255}

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

// Search creates the UI that will allow a user to search and find the desired password for a service.
func Search(data *ListingData, passwords *database.PasswordDB) fyne.CanvasObject {
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

// Add creates the UI that will allow a user to add a new service and the corresponding password.
func Add(data, searchData, removeData *ListingData, passwords *database.PasswordDB) fyne.CanvasObject {
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

	displayMsg := canvas.NewText("", Red)

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

// Remove creates the UI that will allow a user to remove an existing service and password.
func Remove(data, searchData, addData *ListingData, passwords *database.PasswordDB) fyne.CanvasObject {
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
