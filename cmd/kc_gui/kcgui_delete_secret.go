package kc_gui

import (
	"kyrokey/cmd/kc_cli"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ui *UI) ShowDeleteSecret() {

	userEntry := widget.NewEntry()
	serviceEntry := widget.NewEntry()

	serviceEntry.SetPlaceHolder("Service name")
	userEntry.SetPlaceHolder("Username")

	result := widget.NewLabel("")
	form := widget.NewForm(
		widget.NewFormItem("User", userEntry),
		widget.NewFormItem("Service", serviceEntry),
	)

	delBtn := widget.NewButton("Delete Secret", func() {
		user := userEntry.Text
		service := serviceEntry.Text

		secret := kc_cli.KcDeleteSecretCmd(service, user)

		// replace with your real lookup logic
		result.SetText(
			"Found secret for: " + secret,
		)
	})

	screen := container.NewVBox(
		widget.NewLabel("Delete Secret"),
		widget.NewLabel("Please Enter the User and Service related to the Secret"),
		form,
		delBtn,
		result,
		widget.NewButton("Back", func() {
			ui.ShowHome()
		}),
	)

	ui.navigate(screen)
}
