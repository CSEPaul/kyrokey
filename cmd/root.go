/*
Copyright © 2026 kyro
*/
package cmd

import (
	"kyrokey/cmd/kc_cli"
	"kyrokey/cmd/kc_gui"

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
	Use:   "kyro",
	Short: "Kyro sets keys and tracks them",
	Long:  ``,
}

func init() {

	rootCmd.AddCommand(kc_cli.KCCliCmd)
	rootCmd.AddCommand(kc_gui.KCGuiCmd)
	rootCmd.AddGroup(&cobra.Group{ID: "cli", Title: "Cli"})
	rootCmd.AddGroup(&cobra.Group{ID: "gui", Title: "Gui"})
}
