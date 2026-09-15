package config

import (
	"errors"
	"fmt"

	output "github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
)

// PackageManager represents an upgradable package manager
type PackageManager struct {
	Name        string
	Flag        string // CLI flag name (e.g., "brew" for --brew)
	RequiresCmd string // Command that must exist
	UpgradeFn   func() // Function to run upgrades
}

// PackageManagers is the list of all package managers that can be upgraded
var PackageManagers = []PackageManager{
	{
		Name:        "homebrew",
		Flag:        "brew",
		RequiresCmd: "brew",
		UpgradeFn:   upgradeBrew,
	},
	{
		Name:        "npm",
		Flag:        "npm",
		RequiresCmd: "npm",
		UpgradeFn:   upgradeNpm,
	},
	{
		Name:        "bun",
		Flag:        "bun",
		RequiresCmd: "bun",
		UpgradeFn:   upgradeBun,
	},
	{
		Name:        "uv",
		Flag:        "uv",
		RequiresCmd: "uv",
		UpgradeFn:   upgradeUV,
	},
}

// GetPackageManagerByFlag returns a package manager by its flag name
func GetPackageManagerByFlag(flag string) *PackageManager {
	for i := range PackageManagers {
		if PackageManagers[i].Flag == flag {
			return &PackageManagers[i]
		}
	}
	return nil
}

// UpgradeAll upgrades all available package managers
func UpgradeAll() {
	for _, pm := range PackageManagers {
		if CommandExists(pm.RequiresCmd) {
			pm.UpgradeFn()
		}
	}
}

// UpgradePackageManager upgrades a specific package manager
func UpgradePackageManager(pm PackageManager) {
	if !CommandExists(pm.RequiresCmd) {
		fmt.Printf("%s %s not found, skipping\n", output.Cyan("Warning:"), pm.RequiresCmd)
		return
	}
	pm.UpgradeFn()
}

// UpgradePackageByName upgrades a specific package by name
func UpgradePackageByName(name string) error {
	// Find package in our tools list
	pkg := GetToolByName(name)
	if pkg != nil {
		switch pkg.Method {
		case InstallBrewFormula:
			if !CommandExists("brew") {
				return errors.New("brew not found")
			}
			fmt.Printf("  📥 Upgrading %s...\n", name)
			if err := ExecCommand("brew", "upgrade", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallBrewCask:
			if !CommandExists("brew") {
				return errors.New("brew not found")
			}
			fmt.Printf("  📥 Upgrading %s...\n", name)
			if err := ExecCommand("brew", "upgrade", "--cask", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallNpm:
			if !CommandExists("npm") {
				return errors.New("npm not found")
			}
			fmt.Printf("  📥 Upgrading %s...\n", name)
			if err := ExecCommand("npm", "update", "-g", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallBun:
			if !CommandExists("bun") {
				return errors.New("bun not found")
			}
			fmt.Printf("  📥 Upgrading %s...\n", name)
			if err := ExecCommand("bun", "update", "-g", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallUV:
			if !CommandExists("uv") {
				return errors.New("uv not found")
			}
			fmt.Printf("  📥 Upgrading %s...\n", name)
			if err := ExecCommand("uv", "tool", "upgrade", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallNvm, InstallXcode, InstallManual, InstallMAS:
			// No package manager of ours owns these — nvm-managed Node, Xcode, a
			// manual install and a Mac App Store app each upgrade elsewhere. The
			// brew fallback below is the only path j can still offer.
		}
	}

	// Try as a direct brew package name
	if CommandExists("brew") {
		fmt.Printf("  📥 Upgrading %s...\n", name)
		if err := ExecCommand("brew", "upgrade", name); err != nil {
			return fmt.Errorf("upgrading %s: %w", name, err)
		}
		fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
		return nil
	}

	return fmt.Errorf("unknown package: %s", name)
}

// =============================================================================
// Upgrade Functions
// =============================================================================

// makeUpgrader creates a standard upgrade function that prints status and runs commands.
func makeUpgrader(icon, label string, commands ...[]string) func() {
	return func() {
		fmt.Println(output.Cyan(icon + " Upgrading " + label + "..."))
		for _, cmd := range commands {
			if err := ExecCommand(cmd[0], cmd[1:]...); err != nil {
				fmt.Println(output.Cyan("  ⚠ " + cmd[0] + " failed: " + err.Error()))
			}
		}
		fmt.Println(output.Green("  ✅ " + label + " upgrade completed"))
	}
}

var (
	upgradeBrew = makeUpgrader("🍺", "Homebrew packages", []string{"brew", "update"}, []string{"brew", "upgrade"})
	upgradeNpm  = makeUpgrader("📦", "npm global packages", []string{"npm", "update", "-g"})
	upgradeBun  = makeUpgrader("📦", "bun global packages", []string{"bun", "update", "-g"})
	upgradeUV   = makeUpgrader("📦", "uv tools", []string{"uv", "tool", "upgrade", "--all"})
)
