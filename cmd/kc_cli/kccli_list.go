package kc_cli

import (
	"kyrokey/libs"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func KcListCmd() [][2]string {
	db, err := libs.KeychainOpenDB(libs.KeyChainDB)
	if err != nil {
		zap.S().Error(err)
	}
	defer db.Close()

	entries, err := libs.KeychainDBListServicesUsers(db)
	if err != nil {
		zap.S().Error("failed to list kc_cli entries: %v", err)
	}

	return entries
}

var KCCliListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List services and users from local keychain db to this app only",
	Example: `list`,

	Run: func(cmd *cobra.Command, args []string) {
		entries := KcListCmd()

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
