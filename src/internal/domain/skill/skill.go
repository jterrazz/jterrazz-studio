package skill

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/domain/tool"
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

// ListInstalled returns the list of globally installed skill names
func ListInstalled() []string {
	var installed []string

	cmd := exec.Command("skills", "list", "-g")
	output, err := cmd.Output()
	if err != nil {
		return installed
	}

	cleanOutput := tool.StripAnsi(string(output))
	lines := strings.Split(cleanOutput, "\n")
	for _, line := range lines {
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			continue
		}

		line = strings.TrimSpace(line)

		if line == "" ||
			strings.Contains(line, "No global skills") ||
			strings.Contains(line, "Global") ||
			strings.Contains(line, "Skills") ||
			strings.HasPrefix(line, "Try ") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 1 {
			skillName := parts[0]
			if !strings.HasPrefix(skillName, "/") &&
				!strings.HasPrefix(skillName, "~") &&
				!strings.Contains(skillName, ":") &&
				len(skillName) > 0 {
				installed = append(installed, skillName)
			}
		}
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
