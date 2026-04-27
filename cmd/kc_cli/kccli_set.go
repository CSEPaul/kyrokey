package kc_cli

import (
	"fmt"

	"kyrokey/libs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func KcSetCmd(service string, user string, secret string) {
	comment, err := libs.KeyChainSet(service, user, secret)
	if err != nil {
		zap.S().Error("failed to set secret: %v", err)
	}

	fmt.Println(comment)

	//write the service name and user to a sqlite db to track kc_cli entries

	db, err := libs.KeychainOpenDB(libs.KeyChainDB)
	if err != nil {
		zap.S().Error(err)
	}
	defer db.Close()

	if err := libs.KeychainDBEnsureSchema(db); err != nil {
		zap.S().Error(err)
	}

	_ = libs.KeychainDBSaveSecret(db, service, user)
}

var KCCliSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Write a secret to keychain",
	Example: `set -s myapp -u api_key_name -S secretpass`,

	Run: func(cmd *cobra.Command, args []string) {

		service, _ := cmd.Flags().GetString("service")
		user, _ := cmd.Flags().GetString("user")
		secret, _ := cmd.Flags().GetString("secret")

		// Write the secret to the kc_cli
		KcSetCmd(service, user, secret)

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliSetCmd)

	KCCliSetCmd.Flags().StringP("service", "s", "", "The service name for the keyring")
	KCCliSetCmd.Flags().StringP("user", "u", "", "The user name for the keyring")
	KCCliSetCmd.Flags().StringP("secret", "S", "", "The secret value to store in the keyring")
	KCCliSetCmd.MarkFlagRequired("service")
	KCCliSetCmd.MarkFlagRequired("user")
	KCCliSetCmd.MarkFlagRequired("secret")

}
