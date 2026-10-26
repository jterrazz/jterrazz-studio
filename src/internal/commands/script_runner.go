package commands

import (
	"github.com/jterrazz/studio-workstation/src/internal/config"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/print"
)

func runScript(name string) {
	verb, err := config.RunScript(name)
	if err == nil {
		return
	}
	switch verb {
	case "uninstall":
		print.Error("Failed to uninstall " + name + ": " + err.Error())
	case "install":
		print.Error("Failed to install " + name + ": " + err.Error())
	default:
		if config.GetScriptByName(name) == nil {
			print.Error("Unknown script: " + name)
		} else {
			print.Error("No runner for script: " + name)
		}
	}
}
