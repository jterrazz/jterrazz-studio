package commands

import (
	statusview "github.com/jterrazz/jterrazz-studio/src/internal/presentation/views/status"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show comprehensive system status",
	Run: func(cmd *cobra.Command, args []string) {
		statusview.RunOrExit()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
