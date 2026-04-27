package kc_cli

import (
	"fmt"
	"kyrokey/libs"

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

// rrCmd represents the rr command
var KCCliListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List services and users from local keychain db to this app only",
	Example: `list`,

	Run: func(cmd *cobra.Command, args []string) {
		entries := KcListCmd()

		fmt.Println("Stored Keychain Entries:")
		for _, entry := range entries {
			fmt.Printf("- service: %s, user: %s\n", entry[0], entry[1])
		}

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliListCmd)

}
