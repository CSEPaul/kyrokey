package kc_cli

import (
	"github.com/spf13/cobra"
)

// rrCmd represents the rr command
var KCCliCmd = &cobra.Command{
	Use:     "k",
	Short:   "Write a secret to keychain",
	Example: `k`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	KCCliCmd.GroupID = "cli"

}
