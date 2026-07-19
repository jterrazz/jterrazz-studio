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
