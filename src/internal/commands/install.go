package commands

import (
	"fmt"
	"os"
	"sync"

	xterm "github.com/charmbracelet/x/term"
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/domain/tool"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
	installview "github.com/jterrazz/jterrazz-studio/src/internal/presentation/views/install"
	"github.com/spf13/cobra"
)

// installListOnly forces the plain-text catalog listing even when stdout is
// a terminal — the --list escape hatch out of the TUI.
var installListOnly bool

var installCmd = &cobra.Command{
	Use:   "install [tool...]",
	Short: "Install development tools",
	Long: `Install development tools.

With no arguments and an interactive terminal, opens the install TUI: tools
are grouped into four pages — Dev, Infra, AI, Apps — each split into family
sections in registry order. Use arrows/jk to move, tab to fold a section,
space for details, i to install, u to uninstall (with a confirm prompt),
1-4 to jump pages, q/esc to quit.

Falls back to the classic plain-text catalog listing when stdout isn't a
terminal (scripts, CI, piped output) or when --list is passed.

Examples:
  j install                 Open the interactive install TUI
  j install --list          Print the plain-text tool catalog instead
  j install homebrew        Install Homebrew
  j install nvm             Install NVM
  j install go python node  Install specific tools`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var all []string
		for _, t := range config.Tools {
			all = append(all, t.Name)
		}
		return tool.FilterStrings(all, args), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if !installListOnly && xterm.IsTerminal(os.Stdout.Fd()) {
				installview.RunOrExit()
				return
			}
			listAvailableTools()
			return
		}

		print.Action("📦", "Installing selected tools...")
		for _, name := range args {
			installToolByName(name)
		}
		print.Done("Done")
	},
}

func init() {
	installCmd.Flags().BoolVar(&installListOnly, "list", false, "print the plain-text tool catalog instead of opening the TUI")
	rootCmd.AddCommand(installCmd)
}

func listAvailableTools() {
	print.Header("j install", machineContext())

	// Check all tools in parallel
	results := make(map[string]config.CheckResult, len(config.Tools))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := range config.Tools {
		wg.Add(1)
		go func(t *config.Tool) {
			defer wg.Done()
			result := t.Check()
			mu.Lock()
			results[t.Name] = result
			mu.Unlock()
		}(&config.Tools[i])
	}
	wg.Wait()

	knownCategories := make(map[config.ToolCategory]bool, len(config.ToolCategories))
	for _, category := range config.ToolCategories {
		knownCategories[category] = true
		tools := config.GetToolsByCategory(category)
		if len(tools) == 0 {
			continue
		}

		print.Category(string(category))
		for _, t := range tools {
			print.Row(results[t.Name].Installed, t.Name, t.Method.String())
		}
	}

	// Fallback: show any tools using categories not listed in ToolCategories.
	currentCategory := config.ToolCategory("")
	for _, t := range config.Tools {
		if knownCategories[t.Category] {
			continue
		}
		if t.Category != currentCategory {
			currentCategory = t.Category
			print.Category(string(currentCategory))
		}
		print.Row(results[t.Name].Installed, t.Name, t.Method.String())
	}

	print.Empty()
	print.Usage("Usage: j install <tool> [tool...]")
}

func installToolByName(name string) {
	// Handle "brew" as alias for "homebrew"
	if name == "brew" {
		name = "homebrew"
	}

	t := config.GetToolByName(name)
	if t == nil {
		print.Error("Unknown tool: " + name)
		return
	}

	result := t.Check()
	if result.Installed {
		print.Row(true, t.Name, "already installed")
		return
	}

	// Check dependencies
	for _, depName := range t.Dependencies {
		depTool := config.GetToolByName(depName)
		if depTool == nil {
			continue
		}
		depResult := depTool.Check()
		if !depResult.Installed {
			print.Error(depName + " required for " + t.Name + ". Run: j install " + depName)
			return
		}
	}

	print.Installing(t.Name)
	if err := t.Install(); err != nil {
		print.Error(fmt.Sprintf("Failed to install %s: %v", t.Name, err))
	} else {
		print.Row(true, t.Name, "installed")
		// Run post-install scripts
		for _, scriptName := range t.Scripts {
			runScript(scriptName)
		}
	}
}
