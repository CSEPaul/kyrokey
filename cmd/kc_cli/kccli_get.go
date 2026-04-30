package kc_cli

import (
	"os"

	"kyrokey/libs"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func KcGetCmd(service string, user string) string {
	var secret string

	secret, err := libs.KeyChainGet(service, user)
	if err != nil {
		zap.S().Error("failed to get secret: %v", err)
	}
	return secret
}

var KCCliGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get a secret from the keychain",
	Example: `get -s myapp -u username`,

	Run: func(cmd *cobra.Command, args []string) {

		service, _ := cmd.Flags().GetString("service")
		user, _ := cmd.Flags().GetString("user")

		secret := KcGetCmd(service, user)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Secret"})
		t.AppendRow(table.Row{secret})
		t.Render()

	},
}

func init() {

	KCCliCmd.AddCommand(KCCliGetCmd)
	KCCliGetCmd.Flags().StringP("service", "s", "inquiry", "The service name for the keyring")
	KCCliGetCmd.Flags().StringP("user", "u", "inquiry", "The user name for the keyring")
	KCCliGetCmd.MarkFlagRequired("service")
	KCCliGetCmd.MarkFlagRequired("user")

}
