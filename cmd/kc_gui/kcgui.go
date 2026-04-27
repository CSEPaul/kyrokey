package kc_gui

import (
	"fmt"
	"kyrokey/cmd/kc_cli"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
)

type UI struct {
	window  fyne.Window
	content *fyne.Container
}

func (ui *UI) navigate(screen fyne.CanvasObject) {
	ui.content.Objects = []fyne.CanvasObject{screen}
	ui.content.Refresh()
}

// HOME SCREEN

func (ui *UI) ShowHome() {

	screen := container.NewVBox(
		widget.NewLabel("Secret Manager"),

		widget.NewButton("Set", func() {
			ui.ShowSet()
		}),

		widget.NewButton("Get", func() {
			ui.ShowGet()
		}),

		widget.NewButton("List", func() {
			ui.ShowList()
		}),

		widget.NewButton("DeleteDB", func() {
			fmt.Println("delete db")
		}),

		widget.NewButton("DeleteSecret", func() {
			fmt.Println("delete secret")
		}),
	)

	ui.navigate(screen)
}

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

// GET SCREEN

func (ui *UI) ShowGet() {

	service := widget.NewEntry()
	service.SetPlaceHolder("Service name")

	result := widget.NewLabel("")

	getBtn := widget.NewButton("Lookup", func() {

		// replace with your real lookup logic
		result.SetText(
			"Found secret for: " + service.Text,
		)
	})

	screen := container.NewVBox(
		widget.NewLabel("Get Secret"),
		service,
		getBtn,
		result,
		widget.NewButton("Back", func() {
			ui.ShowHome()
		}),
	)

	ui.navigate(screen)
}

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

var KCGuiCmd = &cobra.Command{
	Use:     "g",
	Short:   "Write a secret to keychain using a gui",
	Example: `g`,

	Run: func(cmd *cobra.Command, args []string) {
		a := app.New()
		w := a.NewWindow("Secrets Manager")
		w.Resize(fyne.NewSize(700, 500))

		ui := &UI{
			window:  w,
			content: container.NewStack(),
		}

		ui.ShowHome()

		w.SetContent(ui.content)
		w.ShowAndRun()

	},
}

func init() {

	KCGuiCmd.GroupID = "gui"

}
