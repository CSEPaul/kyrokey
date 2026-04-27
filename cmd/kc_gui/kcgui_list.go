package kc_gui

import (
	"kyrokey/cmd/kc_cli"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LIST SCREEN

func (ui *UI) ShowList() {

	entries := kc_cli.KcListCmd()

	table := widget.NewTable(

		// rows + columns
		func() (int, int) {
			return len(entries) + 1, 2
		},

		// create each cell
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},

		// update cell content
		func(id widget.TableCellID, cell fyne.CanvasObject) {

			label := cell.(*widget.Label)

			// Header row
			if id.Row == 0 {
				if id.Col == 0 {
					label.SetText("Service")
				} else {
					label.SetText("User")
				}
				return
			}

			entry := entries[id.Row-1]

			if id.Col == 0 {
				label.SetText(entry[0]) // service
			} else {
				label.SetText(entry[1]) // user
			}
		},
	)
	// Important: size columns
	table.SetColumnWidth(0, 350) // Service column
	table.SetColumnWidth(1, 350) // User column

	// Optional: taller rows
	table.SetRowHeight(0, 35) // header
	for i := 1; i <= len(entries); i++ {
		table.SetRowHeight(i, 30)
	}

	screen := container.NewBorder(
		widget.NewLabel("Stored Services"),
		widget.NewButton("Back", func() {
			ui.ShowHome()
		}),
		nil,
		nil,
		container.NewScroll(table),
	)

	ui.navigate(screen)
}
