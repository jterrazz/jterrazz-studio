package components

import (
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/theme"
)

// Muted renders text in muted style
func Muted(text string) string {
	return theme.Muted.Render(text)
}

// ColumnSeparator is the standard separator between columns
const ColumnSeparator = "  "

// RenderDescription renders a description in muted style with column separator
func RenderDescription(text string) string {
	if text == "" {
		return ""
	}
	return theme.Muted.Render(ColumnSeparator + text)
}

// PageIndent is the prefix for page-level headers and titles
const PageIndent = " "

// CommandHeaderHeight is the lines a CommandHeader occupies: blank, title, divider, blank.
const CommandHeaderHeight = 4

// CommandHeader renders the canonical header: bold title, muted right-aligned context, a ─ divider.
func CommandHeader(command, context string, width int) string {
	if width < 20 {
		width = 80
	}

	const indent = " "
	title := theme.SectionTitle.Render(command)
	left := indent + title
	leftW := VisibleLen(left)

	titleLine := left
	if context != "" {
		right := theme.Muted.Render(context)
		rightW := VisibleLen(right)
		gap := width - leftW - rightW - 1
		if gap < 1 {
			gap = 1
		}
		titleLine = left + strings.Repeat(" ", gap) + right
	}

	dividerW := width - 1
	if dividerW < 4 {
		dividerW = 4
	}
	divider := indent + theme.Muted.Render(strings.Repeat("─", dividerW))

	return "\n" + titleLine + "\n" + divider + "\n\n"
}
