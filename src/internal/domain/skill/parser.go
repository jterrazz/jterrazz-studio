package skill

import (
	"regexp"
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/domain/tool"
)

// boxCharsRegex matches box-drawing characters used in skills CLI output
var boxCharsRegex = regexp.MustCompile(`[│├└┌◇]`)

// ParseSkillsListOutput parses the output of `skills add <repo> --list`
func ParseSkillsListOutput(output string) []string {
	var skills []string
	cleanOutput := tool.StripAnsi(output)
	lines := strings.Split(cleanOutput, "\n")

	inSkillsSection := false
	for _, line := range lines {
		if strings.Contains(line, "Available Skills") {
			inSkillsSection = true
			continue
		}

		if !inSkillsSection {
			continue
		}

		if strings.Contains(line, "Use --skill") {
			break
		}

		cleaned := boxCharsRegex.ReplaceAllString(line, "")

		leadingSpaces := len(cleaned) - len(strings.TrimLeft(cleaned, " "))
		trimmed := strings.TrimSpace(cleaned)

		if trimmed == "" {
			continue
		}

		if leadingSpaces <= 5 && !strings.Contains(trimmed, " ") && len(trimmed) > 0 {
			if IsValidName(trimmed) {
				skills = append(skills, trimmed)
			}
		}
	}

	return skills
}

// IsValidName checks if a string is a valid skill name (lowercase, numbers, hyphens, underscores)
func IsValidName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}
