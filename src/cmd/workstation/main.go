// Command workstation is the workstation CLI.
package main

import (
	"os"

	"github.com/jterrazz/studio-workstation/src/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
