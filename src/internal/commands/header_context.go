package commands

import (
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
)

func machineContext() string {
	alias, m, ok := config.SelfMachine()
	if !ok {
		return print.MutedText("(unregistered)")
	}
	return alias + " " + print.RenderRole(string(m.Role))
}

func targetContext(alias string) string {
	m, ok := config.GetMachine(alias)
	if !ok {
		return alias
	}
	return alias + " " + print.RenderRole(string(m.Role))
}
