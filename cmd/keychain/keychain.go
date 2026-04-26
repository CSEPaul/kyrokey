package keychain

import (
	"github.com/spf13/cobra"
)

// rrCmd represents the rr command
var KeyChainCmd = &cobra.Command{
	Use:     "k",
	Short:   "Write a secret to keychain",
	Example: `k <secretpass>`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	KeyChainCmd.GroupID = "config"

}
