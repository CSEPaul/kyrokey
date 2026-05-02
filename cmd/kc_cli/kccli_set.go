package kc_cli

import (
	in "kyrokey/internal"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var KCCliSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Write a secret to keychain",
	Example: `set -s myapp -u username -S secretpass`,

	Run: func(cmd *cobra.Command, args []string) {

		service, _ := cmd.Flags().GetString("service")
		user, _ := cmd.Flags().GetString("user")
		secret, _ := cmd.Flags().GetString("secret")

		// Write the secret to the kc_cli
		err := in.KcSet(service, user, secret)
		if err != nil {
			zap.S().Error(err)
		}

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
