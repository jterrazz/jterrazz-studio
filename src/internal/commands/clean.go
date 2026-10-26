package commands

import (
	"github.com/jterrazz/studio-workstation/src/internal/config"
	"github.com/jterrazz/studio-workstation/src/internal/domain/tool"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/print"
	"github.com/spf13/cobra"
)

var cleanAll bool

var cleanCmd = &cobra.Command{
	Use:   "clean [item...]",
	Short: "Clean system caches, Docker, Multipass, and trash",
	Long: `Clean system caches and resources.

Examples:
  workstation clean --all              Clean everything
  workstation clean brew               Clean Homebrew cache
  workstation clean docker             Clean Docker resources
  workstation clean multipass          Clean Multipass instances
  workstation clean trash              Empty trash
  workstation clean brew docker        Clean specific items
  workstation clean                    List available clean items`,
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		all := make([]string, 0, len(config.Cleanables))
		for _, c := range config.Cleanables {
			all = append(all, c.Name)
		}
		return tool.FilterStrings(all, args), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		if cleanAll {
			print.Action("Cleaning everything...")
			for _, c := range config.Cleanables {
				runCleanable(c)
			}
			print.Done("System cleanup completed")
			return
		}

		if len(args) == 0 {
			listCleanItems()
			return
		}

		print.Action("Cleaning selected items...")
		for _, name := range args {
			c := config.GetCleanableByName(name)
			if c == nil {
				print.Error("Unknown clean item: " + name)
				continue
			}
			runCleanable(*c)
		}
		print.Done("Cleanup completed")
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean everything")
	rootCmd.AddCommand(cleanCmd)
}

func listCleanItems() {
	print.Header("workstation clean", machineContext())

	for _, c := range config.Cleanables {
		available := c.RequiresCmd == "" || config.CommandExists(c.RequiresCmd)
		print.Row(available, c.Name, c.Description)
	}

	print.Empty()
	print.Usage(
		"Usage: workstation clean <item> [item...]",
		"       workstation clean --all",
	)
}

func runCleanable(c config.Cleanable) {
	if c.RequiresCmd != "" && !config.CommandExists(c.RequiresCmd) {
		print.Warning(c.RequiresCmd + " not found, skipping")
		return
	}

	print.Action(c.Description + "...")
	if c.CleanFn != nil {
		if err := c.CleanFn(); err != nil {
			print.Error("Failed to clean " + c.Name + ": " + err.Error())
			return
		}
	}
	print.Row(true, c.Name, "completed")
}
