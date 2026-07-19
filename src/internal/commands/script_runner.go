package commands

import (
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
)

// runScript runs a script by name from a non-TUI context — currently used by
// `j install` for post-install scripts attached to a Tool. The j config TUI
// dispatches actions itself (with proper Inputs collection); this helper is
// the headless fallback and ignores any Inputs the script declares.
//
// The dispatch logic (install vs. uninstall, toggle convention) lives in
// config.RunScript, shared with the install TUI's post-install wiring — this
// wrapper only adds the CLI's print-based reporting on top.
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
