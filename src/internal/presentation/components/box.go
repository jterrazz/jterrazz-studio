package components

import (
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/theme"
)

// BoxStyle defines the border style for a box
type BoxStyle int

const (
	BoxRounded BoxStyle = iota
	BoxThick
)

// SectionHeader renders a section divider with horizontal lines
func SectionHeader(title string, width int) string {
	if width < 10 {
		width = 10
	}
	line := theme.SectionBorder.Render(strings.Repeat("━", width))
	label := " " + theme.SectionTitle.Render(strings.ToUpper(title))
	return line + "\n" + label + "\n" + line + "\n"
}

// SubsectionBox renders lines in a rounded box with the title set into its top border.
func SubsectionBox(title string, lines []string, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 20
	}

	borderStyle := theme.SectionBorder

	totalBorderChars := innerWidth + 2
	leftPart := 2
	rightPart := 2
	titleSpace := len(title)
	remainingDashes := totalBorderChars - leftPart - titleSpace - rightPart + 1
	if remainingDashes < 1 {
		remainingDashes = 1
	}

	top := borderStyle.Render(theme.BoxRoundedTopLeft+theme.BoxRoundedHorizontal+" ") +
		theme.SubSection.Render(title) +
		borderStyle.Render(" "+strings.Repeat(theme.BoxRoundedHorizontal, remainingDashes)+theme.BoxRoundedTopRight)

	bottom := borderStyle.Render(theme.BoxRoundedBottomLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedBottomRight)

	paddedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		paddedLines = append(paddedLines, padBoxLine(line, innerWidth))
	}

	return top + "\n" + strings.Join(paddedLines, "\n") + "\n" + bottom
}

// SubsectionBoxWithSeparator is SubsectionBox with a rule between topLines and bottomLines.
func SubsectionBoxWithSeparator(title string, topLines, bottomLines []string, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 20
	}

	borderStyle := theme.SectionBorder

	totalBorderChars := innerWidth + 2
	leftPart := 2
	rightPart := 2
	titleSpace := len(title)
	remainingDashes := totalBorderChars - leftPart - titleSpace - rightPart + 1
	if remainingDashes < 1 {
		remainingDashes = 1
	}

	top := borderStyle.Render(theme.BoxRoundedTopLeft+theme.BoxRoundedHorizontal+" ") +
		theme.SubSection.Render(title) +
		borderStyle.Render(" "+strings.Repeat(theme.BoxRoundedHorizontal, remainingDashes)+theme.BoxRoundedTopRight)

	separator := borderStyle.Render(theme.BoxRoundedTeeLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedTeeRight)

	bottom := borderStyle.Render(theme.BoxRoundedBottomLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedBottomRight)

	paddedTop := make([]string, 0, len(topLines))
	for _, line := range topLines {
		paddedTop = append(paddedTop, padBoxLine(line, innerWidth))
	}
	paddedBottom := make([]string, 0, len(bottomLines))
	for _, line := range bottomLines {
		paddedBottom = append(paddedBottom, padBoxLine(line, innerWidth))
	}

	result := top + "\n" + strings.Join(paddedTop, "\n") + "\n" + separator + "\n" + strings.Join(paddedBottom, "\n") + "\n" + bottom
	return result
}

// SubsectionBoxWithSections renders a box with multiple groups separated by ├──┤.
func SubsectionBoxWithSections(title string, groups [][]string, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 20
	}

	borderStyle := theme.SectionBorder

	totalBorderChars := innerWidth + 2
	leftPart := 2
	rightPart := 2
	titleSpace := len(title)
	remainingDashes := totalBorderChars - leftPart - titleSpace - rightPart + 1
	if remainingDashes < 1 {
		remainingDashes = 1
	}

	top := borderStyle.Render(theme.BoxRoundedTopLeft+theme.BoxRoundedHorizontal+" ") +
		theme.SubSection.Render(title) +
		borderStyle.Render(" "+strings.Repeat(theme.BoxRoundedHorizontal, remainingDashes)+theme.BoxRoundedTopRight)

	separator := borderStyle.Render(theme.BoxRoundedTeeLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedTeeRight)
	bottom := borderStyle.Render(theme.BoxRoundedBottomLeft + strings.Repeat(theme.BoxRoundedHorizontal, innerWidth+2) + theme.BoxRoundedBottomRight)

	var parts []string
	for i, group := range groups {
		if i > 0 {
			parts = append(parts, separator)
		}
		for _, line := range group {
			parts = append(parts, padBoxLine(line, innerWidth))
		}
	}

	return top + "\n" + strings.Join(parts, "\n") + "\n" + bottom
}

func padBoxLine(line string, innerWidth int) string {
	borderStyle := theme.SectionBorder
	visLen := VisibleLen(line)
	if visLen > innerWidth {
		line = truncateAnsi(line, innerWidth)
		visLen = innerWidth
	}
	padding := innerWidth - visLen
	return borderStyle.Render(theme.BoxRoundedVertical+" ") + line + strings.Repeat(" ", padding) + borderStyle.Render(" "+theme.BoxRoundedVertical)
}

func truncateAnsi(s string, maxWidth int) string {
	var result strings.Builder
	visible := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			result.WriteRune(r)
			continue
		}
		if inEscape {
			result.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		if visible >= maxWidth {
			break
		}
		result.WriteRune(r)
		visible++
	}
	result.WriteString("\x1b[0m")
	return result.String()
}

// VisibleLen returns the visible length of a string, stripping ANSI escape codes
func VisibleLen(s string) int {
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}
