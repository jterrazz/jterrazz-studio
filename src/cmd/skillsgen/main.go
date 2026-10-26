// Command skillsgen projects the tool registry into skills/studio-toolbelt/SKILL.md,
// between its generated markers.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/jterrazz/studio-workstation/src/internal/config"
)

const (
	toolsStartMarker  = "<!-- GENERATED:tools by `make skills` from src/internal/config: DO NOT EDIT between markers -->"
	toolsEndMarker    = "<!-- /GENERATED:tools -->"
	preferStartMarker = "<!-- GENERATED:prefer by `make skills` from src/internal/config: DO NOT EDIT between markers -->"
	preferEndMarker   = "<!-- /GENERATED:prefer -->"

	defaultRelPath = "skills/studio-toolbelt/SKILL.md"
)

func main() {
	check := flag.Bool("check", false, "verify the projection is in sync without writing")
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		root, err := findRepoRoot()
		if err != nil {
			fmt.Fprintln(os.Stderr, "skillsgen:", err)
			os.Exit(1)
		}
		path = filepath.Join(root, defaultRelPath)
	}

	if err := run(path, *check); err != nil {
		fmt.Fprintln(os.Stderr, "skillsgen:", err)
		os.Exit(1)
	}
}

func run(path string, check bool) error {
	current, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	rendered, err := render(string(current))
	if err != nil {
		return err
	}

	if check {
		if rendered != string(current) {
			return fmt.Errorf("%s is out of sync with src/internal/config: run `make skills`", path)
		}
		return nil
	}

	//nolint:gosec // a tracked page of the repository, readable like everything else git checks out
	return os.WriteFile(path, []byte(rendered), 0o644)
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not locate repo root (no go.mod found above %s)", dir)
		}
		dir = parent
	}
}

func render(source string) (string, error) {
	withTools, err := spliceBlock(source, toolsStartMarker, toolsEndMarker, renderToolsBlock())
	if err != nil {
		return "", err
	}

	withPrefer, err := spliceBlock(withTools, preferStartMarker, preferEndMarker, renderPreferBlock())
	if err != nil {
		return "", err
	}

	return withPrefer, nil
}

func spliceBlock(source, startMarker, endMarker, block string) (string, error) {
	startIdx := strings.Index(source, startMarker)
	if startIdx == -1 {
		return "", fmt.Errorf("marker not found: %s", startMarker)
	}
	afterStart := startIdx + len(startMarker)

	endOffset := strings.Index(source[afterStart:], endMarker)
	if endOffset == -1 {
		return "", fmt.Errorf("marker not found: %s", endMarker)
	}
	endIdx := afterStart + endOffset

	var b strings.Builder
	b.WriteString(source[:afterStart])
	b.WriteString("\n\n")
	b.WriteString(block)
	b.WriteString("\n\n")
	b.WriteString(source[endIdx:])
	return b.String(), nil
}

func renderToolsBlock() string {
	var sections []string

	for _, category := range config.ToolCategories {
		tools := config.GetToolsByCategory(category)
		if len(tools) == 0 {
			continue
		}

		lines := make([]string, 0, len(tools)+2)
		lines = append(lines, "### "+string(category), "")
		for _, t := range tools {
			lines = append(lines, formatToolLine(t))
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}

	return strings.Join(sections, "\n\n")
}

func formatToolLine(t config.Tool) string {
	if t.Description == "" {
		return fmt.Sprintf("- **%s**", t.Name)
	}
	return fmt.Sprintf("- **%s**: %s", t.Name, t.Description)
}

// invariant: columns are padded to their widest cell, the shape the formatter writes,
// so the formatter and the sync test agree on every byte.
func renderPreferBlock() string {
	rows := [][2]string{{"Instead of", "Use"}}
	for _, t := range config.GetAllTools() {
		if t.Replaces == "" {
			continue
		}
		rows = append(rows, [2]string{t.Replaces, fmt.Sprintf("`%s`: %s", t.Name, t.Description)})
	}

	widths := [2]int{}
	for _, row := range rows {
		for column, cell := range row {
			widths[column] = max(widths[column], utf8.RuneCountInString(cell))
		}
	}

	lines := make([]string, 0, len(rows)+1)
	for i, row := range rows {
		lines = append(lines, fmt.Sprintf("| %s | %s |", pad(row[0], widths[0]), pad(row[1], widths[1])))
		if i == 0 {
			lines = append(lines, fmt.Sprintf("| %s | %s |", strings.Repeat("-", widths[0]), strings.Repeat("-", widths[1])))
		}
	}
	return strings.Join(lines, "\n")
}

func pad(cell string, width int) string {
	return cell + strings.Repeat(" ", width-utf8.RuneCountInString(cell))
}
