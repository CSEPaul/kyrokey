package kc_cli

import (
	in "kyrokey/internal"

	"github.com/spf13/cobra"
)

var KCCliDeleteCmd = &cobra.Command{
	Use:     "del",
	Short:   "Delete a secret for a service and user",
	Example: `del -s myapp -u username`,

	Run: func(cmd *cobra.Command, args []string) {

		service, _ := cmd.Flags().GetString("service")
		user, _ := cmd.Flags().GetString("user")

		in.KcDeleteSecret(service, user)

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliDeleteCmd)
	KCCliDeleteCmd.Flags().StringP("service", "s", "", "The service name for the keyring")
	KCCliDeleteCmd.Flags().StringP("user", "u", "", "The user name for the keyring")

}
