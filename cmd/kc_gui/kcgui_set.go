package kc_gui

import (
	"kyrokey/cmd/kc_cli"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
		err := kc_cli.KcSetCmd(service, user, secret)
		if err != nil {
			status.SetText(err.Error())
			return
		}

		status.SetText("Saved successfully")
	})
	resetBtn := widget.NewButton("Reset", func() {
		dialog.ShowConfirm(
			"Clear Form",
			"Reset all fields?",
			func(ok bool) {
				if ok {
					userEntry.SetText("")
					secretEntry.SetText("")
					serviceEntry.SetText("")
				}
			},
			ui.window,
		)
	})

	backBtn := widget.NewButton("Back", func() {
		ui.ShowHome()
	})

	screen := container.NewVBox(
		widget.NewLabel("Set Secret"),
		form,
		saveBtn,
		resetBtn,
		backBtn,
		status,
	)

	ui.navigate(screen)
}
