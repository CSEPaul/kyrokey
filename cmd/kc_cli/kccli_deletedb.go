package kc_cli

import (
	"fmt"
	in "kyrokey/internal"

	"github.com/spf13/cobra"
)

var KCCliDeleteDBCmd = &cobra.Command{
	Use:     "deldb",
	Short:   "Delete the keychain db tracker only",
	Example: `deldb --Confirm Confirm`,

	Run: func(cmd *cobra.Command, args []string) {

		confirm, _ := cmd.Flags().GetString("Confirm")
		result := in.KcDeleteDB(confirm)
		fmt.Println(result)

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliDeleteDBCmd)
	KCCliDeleteDBCmd.Flags().StringP("Confirm", "C", "", "Confirm flag for security")
	KCCliDeleteDBCmd.MarkFlagRequired("Confirm")

}
