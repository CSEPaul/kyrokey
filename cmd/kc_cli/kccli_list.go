package kc_cli

import (
	in "kyrokey/internal"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var KCCliListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List services and users from local keychain db to this app only",
	Example: `list`,

	Run: func(cmd *cobra.Command, args []string) {
		entries := in.KcList()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Service", "User"})

		for _, entry := range entries {
			t.AppendRow(table.Row{entry[0], entry[1]})
		}

		t.Render()

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliListCmd)

}
