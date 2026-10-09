// Package tool answers whether a tool is installed and at which version, by
// probing commands and reading package-manager output. It imports no other
// package of this binary.
package tool

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// VersionFromCmd creates a version func that runs a command and parses output
func VersionFromCmd(cmd string, args []string, parser func(string) string) func() string {
	return func() string {
		out, err := exec.Command(cmd, args...).CombinedOutput()
		if err != nil {
			return ""
		}
		return parser(string(out))
	}
}

// VersionFromBrewFormula creates a version func that gets version from brew info
func VersionFromBrewFormula(formula string) func() string {
	return func() string {
		out, err := exec.Command("brew", "list", "--versions", formula).Output()
		if err != nil {
			return ""
		}
		// outside: brew prints "formula 1.2.3" or "formula 1.2.3 1.2.2".
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			return parts[1]
		}
		return ""
	}
}

// VersionFromBrewCask creates a version func that gets version from brew cask info
func VersionFromBrewCask(cask string) func() string {
	return func() string {
		out, err := exec.Command("brew", "list", "--cask", "--versions", cask).Output()
		if err != nil {
			return ""
		}
		// outside: brew prints "cask 1.2.3" or "cask 0.34.1,HASH".
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			version := parts[1]
			if idx := strings.Index(version, ","); idx > 0 {
				version = version[:idx]
			}
			return version
		}
		return ""
	}
}

// VersionFromAppPlist creates a version func that reads version from app's Info.plist
func VersionFromAppPlist(appName string) func() string {
	return func() string {
		plistPath := fmt.Sprintf("/Applications/%s.app/Contents/Info.plist", appName)
		out, err := exec.Command("defaults", "read", plistPath, "CFBundleShortVersionString").Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
}

// VersionFromBunGlobal creates a version func that gets version from bun global packages
func VersionFromBunGlobal(pkg string) func() string {
	return func() string {
		out, err := exec.Command("bun", "pm", "ls", "-g").Output()
		if err != nil {
			return ""
		}
		// outside: the global list prints "├── package@version" or "└── package@version".
		for _, line := range strings.Split(string(out), "\n") {
			idx := strings.Index(line, "── ")
			if idx < 0 {
				continue
			}
			entry := strings.TrimSpace(line[idx+len("── "):])
			// outside: an entry reads "qmd@1.2.3" or "qmd@github:tobi/qmd#hash".
			atIdx := strings.LastIndex(entry, "@")
			if atIdx <= 0 {
				continue
			}
			name := entry[:atIdx]
			version := entry[atIdx+1:]
			if name == pkg {
				version = strings.TrimPrefix(version, "github:")
				return version
			}
		}
		return ""
	}
}

// CommandExists checks if a command exists in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// GetCommandOutput runs a command and returns its trimmed output, or empty string on error
func GetCommandOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GetCommandOutputWithTimeout runs a command with a timeout and returns its trimmed output
func GetCommandOutputWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command timed out after %v", timeout)
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
