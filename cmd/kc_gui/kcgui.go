package kc_gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
)

func doSet() error {
	return nil
}

func doGet() (string, error) {
	return "example-secret", nil
}

func doList() ([]string, error) {
	return []string{
		"secret1",
		"secret2",
		"secret3",
	}, nil
}

func doDeleteDB() error {
	return nil
}

func doDeleteSecret() error {
	return nil
}

// rrCmd represents the rr command
var KCGuiCmd = &cobra.Command{
	Use:     "g",
	Short:   "Write a secret to keychain using a gui",
	Example: `g <secretpass>`,

	Run: func(cmd *cobra.Command, args []string) {
		a := app.New()
		w := a.NewWindow("Secret Manager")

		output := widget.NewMultiLineEntry()
		output.SetPlaceHolder("Output...")
		output.Disable()

		appendOutput := func(msg string) {
			output.SetText(output.Text + msg + "\n")
		}

		// --- Button handlers ---

		setBtn := widget.NewButton("Set", func() {
			// Replace with your real set logic
			appendOutput("Running SET...")
			err := doSet()
			if err != nil {
				appendOutput("Error: " + err.Error())
				return
			}
			appendOutput("Secret stored successfully.")
		})

		getBtn := widget.NewButton("Get", func() {
			appendOutput("Running GET...")
			val, err := doGet()
			if err != nil {
				appendOutput("Error: " + err.Error())
				return
			}
			appendOutput("Secret: " + val)
		})

		listBtn := widget.NewButton("List", func() {
			appendOutput("Running LIST...")
			items, err := doList()
			if err != nil {
				appendOutput("Error: " + err.Error())
				return
			}

			for _, item := range items {
				appendOutput(item)
			}
		})

		deleteDBBtn := widget.NewButton("DeleteDB", func() {
			appendOutput("Running DELETE DB...")
			err := doDeleteDB()
			if err != nil {
				appendOutput("Error: " + err.Error())
				return
			}
			appendOutput("Database deleted.")
		})

		deleteSecretBtn := widget.NewButton("DeleteSecret", func() {
			appendOutput("Running DELETE SECRET...")
			err := doDeleteSecret()
			if err != nil {
				appendOutput("Error: " + err.Error())
				return
			}
			appendOutput("Secret deleted.")
		})

		buttons := container.NewVBox(
			setBtn,
			getBtn,
			listBtn,
			deleteDBBtn,
			deleteSecretBtn,
		)

		content := container.NewBorder(
			nil,
			nil,
			buttons,
			nil,
			output,
		)

		w.SetContent(content)
		w.Resize(fyne.NewSize(700, 400))
		w.ShowAndRun()
	},
}

func init() {

	KCGuiCmd.GroupID = "gui"

}
