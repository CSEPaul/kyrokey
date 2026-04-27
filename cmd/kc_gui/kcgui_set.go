package kc_gui

import (
	"kyrokey/cmd/kc_cli"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SET SCREEN

func (ui *UI) ShowSet() {

	userEntry := widget.NewEntry()
	secretEntry := widget.NewPasswordEntry()
	serviceEntry := widget.NewEntry()

	status := widget.NewLabel("")

	form := widget.NewForm(
		widget.NewFormItem("User", userEntry),
		widget.NewFormItem("Secret", secretEntry),
		widget.NewFormItem("Service", serviceEntry),
	)

	saveBtn := widget.NewButton("Save", func() {

		user := userEntry.Text
		secret := secretEntry.Text
		service := serviceEntry.Text

		// call your Cobra/business logic here
		kc_cli.KcSetCmd(service, user, secret)

		status.SetText("Saved successfully")
	})

	backBtn := widget.NewButton("Back", func() {
		ui.ShowHome()
	})

	screen := container.NewVBox(
		widget.NewLabel("Set Secret"),
		form,
		saveBtn,
		backBtn,
		status,
	)

	ui.navigate(screen)
}
