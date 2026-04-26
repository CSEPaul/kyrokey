/*
Copyright © 2026 kyro
*/
package cmd

import (
	"kyrokey/cmd/keychain"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyrokey",
	Short: "Kyrokey sets keys and tracks them",
	Long:  ``,
}

func init() {

	rootCmd.AddCommand(keychain.KeyChainCmd)
	rootCmd.AddGroup(&cobra.Group{ID: "config", Title: "Config"})
}
