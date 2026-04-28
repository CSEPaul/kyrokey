package kc_gui

import (
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
			ui.ShowDeleteDb()
		}),

		widget.NewButton("DeleteSecret", func() {
			ui.ShowDeleteSecret()
		}),
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
