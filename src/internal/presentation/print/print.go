// Package print renders the plain-text lines the commands emit outside a TUI:
// headers, role pills and the verdicts an install action reports.
package print

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/components"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/theme"
)

func Line(s string)                    { fmt.Println(s) }
func Linef(format string, args ...any) { fmt.Printf(format+"\n", args...) }
func Empty()                           { fmt.Println() }

func Error(msg string)   { fmt.Fprintf(os.Stderr, "%s %s\n", theme.Danger.Render("Error:"), msg) }
func Warning(msg string) { fmt.Printf("%s %s\n", theme.Warning.Render("Warning:"), msg) }
func Success(msg string) { fmt.Printf("%s %s\n", theme.Success.Render(theme.IconCheck), msg) }
func Info(msg string)    { fmt.Println(theme.Special.Render(msg)) }
func Dim(msg string)     { fmt.Println(theme.Muted.Render(msg)) }

func Action(msg string) { fmt.Println(theme.Special.Render(msg)) }
func Done(msg string)   { fmt.Println(theme.Success.Render(theme.IconCheck + " " + msg)) }

func Installing(name string) {
	fmt.Printf(components.PageIndent+"Installing %s...\n", name)
}

var (
	roleClientStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(theme.ColorClient)).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)
	roleServerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(theme.ColorServer)).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1)
)

// RenderRole renders a role as a colored pill, plain text for an unknown role.
func RenderRole(role string) string {
	switch role {
	case "client":
		return roleClientStyle.Render(role)
	case "server":
		return roleServerStyle.Render(role)
	}
	return role
}

// MutedText renders s in the muted style.
func MutedText(s string) string { return theme.Muted.Render(s) }

// Header prints the canonical command header to stdout; a TUI embeds RenderHeader instead.
func Header(command, context string) {
	fmt.Print(RenderHeader(command, context, termWidth()))
}

// RenderHeader returns the canonical header as a string sized to width.
func RenderHeader(command, context string, width int) string {
	return components.CommandHeader(command, context, width)
}

// Category prints a category header (dimmed)
func Category(name string) {
	fmt.Println(theme.Muted.Render(name))
}

func termWidth() int {
	if w, _ := strconv.Atoi(os.Getenv("COLUMNS")); w > 0 {
		return w
	}
	out, err := exec.Command("tput", "cols").Output()
	if err == nil {
		if w, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && w > 0 {
			return w
		}
	}
	return 80
}

func Row(ok bool, label, detail string) {
	icon := components.Badge(ok)
	if detail != "" {
		fmt.Printf(components.PageIndent+"%s %-14s %s\n", icon, label, theme.Muted.Render(detail))
	} else {
		fmt.Printf(components.PageIndent+"%s %s\n", icon, label)
	}
}

func Usage(lines ...string) {
	for _, line := range lines {
		fmt.Println(theme.Muted.Render(line))
	}
}

var (
	cyan  = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSpecial))
	green = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSuccess))
	dim   = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorMuted))
)

func Cyan(s string) string   { return cyan.Render(s) }
func Green(s string) string  { return green.Render(s) }
func Dimmed(s string) string { return dim.Render(s) }
