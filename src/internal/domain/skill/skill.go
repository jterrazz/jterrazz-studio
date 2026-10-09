package skill

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Install installs a skill from a repo globally
func Install(repo, skill string) error {
	cmd := exec.Command("skills", "add", repo, "-g", "-y", "--skill", skill)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return nil
}

// InstallAll installs all skills from a repo globally
func InstallAll(repo string) error {
	cmd := exec.Command("skills", "add", repo, "-g", "-y", "--all")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return nil
}

// Remove removes a skill globally
func Remove(skill string) error {
	cmd := exec.Command("skills", "remove", "-g", "-y", skill)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove %s: %w", skill, err)
	}
	return nil
}

// RemoveAll removes all skills globally
func RemoveAll() error {
	cmd := exec.Command("skills", "remove", "-g", "-y", "--all")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove all skills: %w", err)
	}
	return nil
}

func Update(name string) error {
	cmd := exec.Command("skills", "update", "-g", "-y", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return nil
}

func skillsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".agents", "skills")
}

// ListInstalled lists the skill folders of ~/.agents/skills that hold a SKILL.md.
func ListInstalled() []string {
	return listInstalledIn(skillsDir())
}

func listInstalledIn(dir string) []string {
	// outside: `skills list -g` changed its output shape across versions, so the
	// filesystem is read instead of the CLI's listing.
	var installed []string
	if dir == "" {
		return installed
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return installed
	}

	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		// outside: os.Stat, unlike DirEntry.IsDir, follows symlinks, so a symlinked skill folder counts.
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil || !info.IsDir() {
			continue
		}

		if _, err := os.Stat(filepath.Join(dir, name, "SKILL.md")); err != nil {
			continue
		}

		installed = append(installed, name)
	}

	return installed
}

// ListFromRepo fetches available skills from a repo
func ListFromRepo(repo string) ([]string, error) {
	cmd := exec.Command("skills", "add", repo, "--list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return ParseSkillsListOutput(string(output)), nil
}

// IsInstalled checks if the skills CLI is available
func IsInstalled() bool {
	_, err := exec.LookPath("skills")
	return err == nil
}
