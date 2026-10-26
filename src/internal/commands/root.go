package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "workstation",
	Short: "studio-workstation unified command system",
	Long:  "A unified CLI tool for development workflow automation.",
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func Execute() error {
	return rootCmd.Execute()
}
