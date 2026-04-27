package kc_cli

import (
	"kyrokey/libs"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rrCmd represents the rr command
var KCCliDeleteDBCmd = &cobra.Command{
	Use:     "deldb",
	Short:   "Delete the keychain db tracker only",
	Example: `deldb --Confirm Confirm`,

	Run: func(cmd *cobra.Command, args []string) {

		confirm, _ := cmd.Flags().GetString("Confirm")

		dir, err := libs.KeyChainDBFilePath(libs.KeyChainDB)
		if err != nil {
			zap.S().Error("Error getting db file path:", err.Error())
		}

		switch confirm {
		case "Confirm":
			// delete the db file
			libs.KeyChainDeleteDBFile(dir)
		case "C":
			// delete the db file
			libs.KeyChainDeleteDBFile(dir)
		default:
			println("You must use the Confirm flag to delete the db file for security.")
		}
	},
}

func init() {

	KCCliCmd.AddCommand(KCCliDeleteDBCmd)
	KCCliDeleteDBCmd.Flags().StringP("Confirm", "C", "", "Confirm flag for security")
	KCCliDeleteDBCmd.MarkFlagRequired("Confirm")

}
