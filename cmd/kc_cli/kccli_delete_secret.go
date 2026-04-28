package kc_cli

import (
	"fmt"

	"kyrokey/libs"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func KcDeleteSecretCmd(service string, user string) string {
	err := libs.KeyChainDelete(service, user)
	if err != nil {
		zap.S().Error("failed to delete secret: %v", err)
	}

	fmt.Println("Secret deleted from keyring related to:--- service:", service, "and user:", user)

	//delete the service + user entry from sqlite db
	db, err := libs.KeychainOpenDB(libs.KeyChainDB)
	if err != nil {
		zap.S().Error(err)
	}
	defer db.Close()
	err = libs.KeychainDBDeleteSecret(db, service, user)
	if err != nil {
		zap.S().Error(err)
	}
	fmt.Println("All users related to service:", service, "deleted from kc_cli tracking db.")
	return "Secret deleted from kc_cli tracking db."
}

var KCCliDeleteCmd = &cobra.Command{
	Use:     "del",
	Short:   "Delete a secret for a service and user",
	Example: `kc_cli del -s myapp -u username`,

	Run: func(cmd *cobra.Command, args []string) {

		service, _ := cmd.Flags().GetString("service")
		user, _ := cmd.Flags().GetString("user")

		KcDeleteSecretCmd(service, user)

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliDeleteCmd)
	KCCliDeleteCmd.Flags().StringP("service", "s", "", "The service name for the keyring")
	KCCliDeleteCmd.Flags().StringP("user", "u", "", "The user name for the keyring")

}
