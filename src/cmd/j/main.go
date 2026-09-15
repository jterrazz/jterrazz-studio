// Command j is the workstation CLI: it runs the single Cobra root declared in
// internal/commands and exits non-zero when a verb fails.
package main

import (
	"os"

	"github.com/jterrazz/jterrazz-studio/src/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
