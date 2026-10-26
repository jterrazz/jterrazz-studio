package config

import (
	"errors"
	"fmt"

	output "github.com/jterrazz/studio-workstation/src/internal/presentation/print"
)

// PackageManager represents an upgradable package manager
type PackageManager struct {
	Name        string
	Flag        string
	RequiresCmd string
	UpgradeFn   func()
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
	pkg := GetToolByName(name)
	if pkg != nil {
		switch pkg.Method {
		case InstallBrewFormula:
			if !CommandExists("brew") {
				return errors.New("brew not found")
			}
			fmt.Printf("  Upgrading %s...\n", name)
			if err := ExecCommand("brew", "upgrade", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallBrewCask:
			if !CommandExists("brew") {
				return errors.New("brew not found")
			}
			fmt.Printf("  Upgrading %s...\n", name)
			if err := ExecCommand("brew", "upgrade", "--cask", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallNpm:
			if !CommandExists("npm") {
				return errors.New("npm not found")
			}
			fmt.Printf("  Upgrading %s...\n", name)
			if err := ExecCommand("npm", "update", "-g", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallBun:
			if !CommandExists("bun") {
				return errors.New("bun not found")
			}
			fmt.Printf("  Upgrading %s...\n", name)
			if err := ExecCommand("bun", "update", "-g", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallUV:
			if !CommandExists("uv") {
				return errors.New("uv not found")
			}
			fmt.Printf("  Upgrading %s...\n", name)
			if err := ExecCommand("uv", "tool", "upgrade", pkg.Formula); err != nil {
				return fmt.Errorf("upgrading %s: %w", name, err)
			}
			fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
			return nil
		case InstallNvm, InstallXcode, InstallManual, InstallMAS:
			// invariant: no package manager of workstation's owns these (nvm Node, Xcode, manual, App Store);
			// the brew fallback below is the only path workstation still offers.
		}
	}

	if CommandExists("brew") {
		fmt.Printf("  Upgrading %s...\n", name)
		if err := ExecCommand("brew", "upgrade", name); err != nil {
			return fmt.Errorf("upgrading %s: %w", name, err)
		}
		fmt.Printf("  %s %s upgraded\n", output.Green("✓"), name)
		return nil
	}

	return fmt.Errorf("unknown package: %s", name)
}

func makeUpgrader(label string, commands ...[]string) func() {
	return func() {
		fmt.Println(output.Cyan("Upgrading " + label + "..."))
		for _, cmd := range commands {
			if err := ExecCommand(cmd[0], cmd[1:]...); err != nil {
				fmt.Println(output.Cyan("  ⚠ " + cmd[0] + " failed: " + err.Error()))
			}
		}
		fmt.Println(output.Green("  ✓ " + label + " upgrade completed"))
	}
}

var (
	upgradeBrew = makeUpgrader("Homebrew packages", []string{"brew", "update"}, []string{"brew", "upgrade"})
	upgradeNpm  = makeUpgrader("npm global packages", []string{"npm", "update", "-g"})
	upgradeBun  = makeUpgrader("bun global packages", []string{"bun", "update", "-g"})
	upgradeUV   = makeUpgrader("uv tools", []string{"uv", "tool", "upgrade", "--all"})
)
