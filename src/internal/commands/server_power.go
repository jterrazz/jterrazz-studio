package commands

import (
	"fmt"
	"strings"

	"github.com/jterrazz/studio-workstation/src/internal/config"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/print"
)

// outside: on Apple Silicon SleepOnPowerButton and hibernatemode are firmware-controlled
// and absent from `pmset -g custom`, so they are neither set nor checked.
var powerHardenSettings = [][2]string{
	{"autorestart", "1"},
	{"sleep", "0"},
	{"displaysleep", "5"},
	{"disksleep", "0"},
	{"powernap", "0"},
	{"womp", "1"},
}

func installPower() error {
	if err := requireDarwin(); err != nil {
		return err
	}
	if err := requireRoot(); err != nil {
		return err
	}

	print.Header("install power", machineContext())
	print.Category("Before")
	dumpPmset()
	print.Empty()

	print.Category("Applying")
	args := []string{"-a"}
	for _, kv := range powerHardenSettings {
		args = append(args, kv[0], kv[1])
	}
	if err := run("/usr/bin/pmset", args...); err != nil {
		return err
	}
	for _, kv := range powerHardenSettings {
		print.Success(kv[0] + "=" + kv[1])
	}

	print.Empty()
	print.Category("After")
	dumpPmset()
	return nil
}

func uninstallPower() error {
	if err := requireDarwin(); err != nil {
		return err
	}
	if err := requireRoot(); err != nil {
		return err
	}

	print.Header("uninstall power", machineContext())
	print.Category("Before")
	dumpPmset()
	print.Empty()

	print.Category("Applying")
	if err := run("/usr/bin/pmset", "-a", "-resetdefaults"); err != nil {
		return fmt.Errorf("pmset -resetdefaults: %w", err)
	}
	print.Success("pmset reset to macOS defaults")

	print.Empty()
	print.Category("After")
	dumpPmset()
	return nil
}

func checkPowerInstalled() config.CheckResult {
	out, err := runQuiet("/usr/bin/pmset", "-g", "custom")
	if err != nil {
		return config.CheckResult{}
	}
	settings := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 {
			settings[fields[0]] = fields[1]
		}
	}
	for _, kv := range powerHardenSettings {
		if settings[kv[0]] != kv[1] {
			return config.CheckResult{}
		}
	}
	return config.InstalledWithDetail("harden applied")
}

func dumpPmset() {
	out, err := runQuiet("/usr/bin/pmset", "-g", "custom")
	if err != nil {
		print.Warning("pmset -g custom failed: " + err.Error())
		return
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		print.Linef("  %s", line)
	}
}
