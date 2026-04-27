package kc_gui

import (
	"kyrokey/cmd/kc_cli"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// GET SCREEN

func (ui *UI) ShowGet() {

	userEntry := widget.NewEntry()
	serviceEntry := widget.NewEntry()

	serviceEntry.SetPlaceHolder("Service name")
	userEntry.SetPlaceHolder("Username")

	result := widget.NewLabel("")
	form := widget.NewForm(
		widget.NewFormItem("User", userEntry),
		widget.NewFormItem("Service", serviceEntry),
	)

	getBtn := widget.NewButton("Lookup", func() {
		user := userEntry.Text
		service := serviceEntry.Text

		secret := kc_cli.KcGetCmd(service, user)

		// replace with your real lookup logic
		result.SetText(
			"Found secret for: " + secret,
		)
	})

	screen := container.NewVBox(
		widget.NewLabel("Get Secret"),
		form,
		getBtn,
		result,
		widget.NewButton("Back", func() {
			ui.ShowHome()
		}),
	)

	ui.navigate(screen)
}
