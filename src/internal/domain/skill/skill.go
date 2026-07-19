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

// Update updates a single globally installed skill to its latest upstream
// version via the `skills` CLI.
func Update(name string) error {
	cmd := exec.Command("skills", "update", "-g", "-y", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return nil
}

// skillsDir returns ~/.agents/skills, the directory `skills add -g` installs
// into. Returns "" if the home directory can't be resolved.
func skillsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".agents", "skills")
}

// ListInstalled returns the names of globally installed skills.
//
// This reads the filesystem directly — entries of ~/.agents/skills that are
// directories (or symlinks to directories) containing a SKILL.md — rather
// than parsing `skills list -g` output. The CLI's list-output format has
// changed shape across versions (indentation, headers, columns), which made
// output-parsing fragile; the filesystem is ground truth for what's actually
// installed. A missing skills directory returns an empty list.
func ListInstalled() []string {
	return listInstalledIn(skillsDir())
}

// listInstalledIn is the testable variant of ListInstalled — takes the
// skills directory explicitly instead of resolving $HOME.
func listInstalledIn(dir string) []string {
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

		// os.Stat (unlike DirEntry.IsDir) follows symlinks, so symlinked
		// skill folders count too.
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
