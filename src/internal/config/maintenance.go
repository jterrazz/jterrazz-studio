package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type MaintenanceCheck struct {
	Name        string
	Description string
	CheckFn     func() CheckResult
	GoodWhen    bool
}

var MaintenanceChecks = []MaintenanceCheck{
	{
		Name:        "macOS",
		Description: "Pending macOS software updates",
		CheckFn:     checkMacOSUpdates,
		GoodWhen:    false,
	},
	{
		Name:        "brew",
		Description: "Outdated Homebrew formulae and casks",
		CheckFn:     checkBrewOutdated,
		GoodWhen:    false,
	},
}

// outside: softwareupdate touches /Library/Updates/index.plist once an update is
// downloaded, so reading it spares a `softwareupdate -l` network round trip.
func checkMacOSUpdates() CheckResult {
	info, err := os.Stat("/Library/Updates/index.plist")
	if err != nil {
		return CheckResult{Installed: false, Detail: "up to date"}
	}
	age := humanizeAge(time.Since(info.ModTime()))
	return CheckResult{Installed: true, Detail: "pending · " + age + " old"}
}

// outside: `brew outdated --quiet` prints one outdated name per line on stdout; its
// stderr is noise (permission warnings, cache messages).
func checkBrewOutdated() CheckResult {
	if _, err := exec.LookPath("brew"); err != nil {
		return CheckResult{}
	}
	cmd := exec.Command("brew", "outdated", "--quiet")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return CheckResult{Detail: "unknown"}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return CheckResult{Installed: false, Detail: "all current"}
	}
	return CheckResult{Installed: true, Detail: fmt.Sprintf("%d outdated", len(lines))}
}

func humanizeAge(d time.Duration) string {
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours())/24)
	}
}

// DaemonCheck is one LaunchAgent of ~/Library/LaunchAgents whose name matches daemonNamePrefixes.
type DaemonCheck struct {
	Label string // e.g. "ai.jterrazz.lock-after-login"
	Path  string // full plist path
}

// value: only the estate's own agents: macOS and third-party apps ship dozens more.
var daemonNamePrefixes = []string{"ai.jterrazz.", "com.jterrazz.", "ai.openclaw.", "ai.hermes."}

// DiscoverDaemons lists the matching LaunchAgents, sorted by label.
func DiscoverDaemons() []DaemonCheck {
	dir := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []DaemonCheck
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".plist") {
			continue
		}
		label := strings.TrimSuffix(name, ".plist")
		matched := false
		for _, p := range daemonNamePrefixes {
			if strings.HasPrefix(label, p) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		out = append(out, DaemonCheck{Label: label, Path: filepath.Join(dir, name)})
	}
	return out
}

// CheckDaemonState reads launchctl list: loaded is Installed, Detail is running PID, idle or exit N.
func CheckDaemonState(label string) CheckResult {
	out, err := exec.Command("launchctl", "list").Output()
	if err != nil {
		return CheckResult{Detail: "launchctl error"}
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[2] != label {
			continue
		}
		pid, exitCode := fields[0], fields[1]
		switch {
		case pid != "-":
			return CheckResult{Installed: true, Detail: "running · pid " + pid}
		case exitCode != "0" && exitCode != "-":
			return CheckResult{Installed: true, Detail: "idle · last exit " + exitCode}
		default:
			return CheckResult{Installed: true, Detail: "idle"}
		}
	}
	return CheckResult{Installed: false, Detail: "not loaded"}
}
