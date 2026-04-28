package kc_cli

import (
	"kyrokey/libs"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func KcDeleteDBCmd(confirm string) string {
	dir, err := libs.KeyChainDBFilePath(libs.KeyChainDB)
	if err != nil {
		zap.S().Error("Error getting db file path:", err.Error())
	}

	comment := "Secret Deleted"
	comment2 := "Write the `Confirm` flag"
	switch confirm {
	case "Confirm":
		err := libs.KeyChainDeleteDBFile(dir)
		if err != nil {
			zap.S().Error("Error deleting db file:", err.Error())
		}
		return comment
	case "C":
		err := libs.KeyChainDeleteDBFile(dir)
		if err != nil {
			zap.S().Error("Error deleting db file:", err.Error())
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
		KcDeleteDBCmd(confirm)

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliDeleteDBCmd)
	KCCliDeleteDBCmd.Flags().StringP("Confirm", "C", "", "Confirm flag for security")
	KCCliDeleteDBCmd.MarkFlagRequired("Confirm")

}
