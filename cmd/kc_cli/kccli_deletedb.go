package kc_cli

import (
	"fmt"
	"kyrokey/libs"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func KcDeleteDBCmd(confirm string) string {
	dir, err := libs.KeyChainDBFilePath(libs.KeyChainDB)
	if err != nil {
		zap.S().Error("Error getting db file path:", err.Error())
	}

	comment := "Secret DB Deleted"
	comment2 := "Write the `Confirm` flag"
	switch confirm {
	case "Confirm", "C":
		err := libs.KeyChainDeleteDBFile(dir)
		if err != nil {
			zap.S().Error("Error deleting db file:", err.Error())
		}
		// check if db is present - keychain_entries.db
		exists := libs.FileExists(libs.KeyChainDB)
		if exists {
			statement := "DB Not Deleted"
			return statement
		}
		return comment

	default:
		println("You must use the Confirm flag to delete the db file for security.")
		return comment2
	}

}

var KCCliDeleteDBCmd = &cobra.Command{
	Use:     "deldb",
	Short:   "Delete the keychain db tracker only",
	Example: `deldb --Confirm Confirm`,

	Run: func(cmd *cobra.Command, args []string) {

		confirm, _ := cmd.Flags().GetString("Confirm")
		result := KcDeleteDBCmd(confirm)
		fmt.Println(result)

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliDeleteDBCmd)
	KCCliDeleteDBCmd.Flags().StringP("Confirm", "C", "", "Confirm flag for security")
	KCCliDeleteDBCmd.MarkFlagRequired("Confirm")

}
