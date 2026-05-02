package kc_gui

import (
	in "kyrokey/internal"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ui *UI) ShowDeleteDb() {

	confirmEntry := widget.NewEntry()
	confirmEntry.SetPlaceHolder("Confirm")

	result := widget.NewLabel("")
	form := widget.NewForm(
		widget.NewFormItem("Confirm", confirmEntry),
	)

	delBtn := widget.NewButton("Enter", func() {

		confirm := confirmEntry.Text
		comment := in.KcDeleteDB(confirm)

		// replace with your real lookup logic
		result.SetText(
			"Results: " + comment,
		)
	})

	screen := container.NewVBox(
		widget.NewLabel("Delete DB"),
		form,
		delBtn,
		result,
		widget.NewButton("Back", func() {
			ui.ShowHome()
		}),
	)

	ui.navigate(screen)
}
