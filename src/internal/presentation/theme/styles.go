package theme

import "github.com/charmbracelet/lipgloss"

// =============================================================================
// Base Styles - Semantic styles for text
// =============================================================================

var (
	// Title is for main headings
	Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		Bold(true)

	// Section is for section headers with margin
	Section = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSecondary)).
		Bold(true).
		MarginTop(1)

	// SubSection is for subsection headers
	SubSection = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Italic(true)

	// Selected is for currently selected items
	Selected = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true)

	// Success is for positive states (installed, running, etc.)
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess))

	// Warning is for cautionary states
	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWarning))

	// Action is for actionable items
	Action = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWarning)).
		Bold(true)

	// Danger is for errors/negative states
	Danger = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorDanger))

	// Muted is for less important text
	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted))

	// Normal is for regular text
	Normal = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText))

	// Special is for highlighted values (versions, paths, etc.)
	Special = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSpecial))

	// Help is for help text at bottom
	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		MarginTop(1)
)

// =============================================================================
// Border Styles
// =============================================================================

var (
	// Border is for drawing borders
	Border = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBorder))

	// SectionBorder is for section header borders
	SectionBorder = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBorder))

	// SectionTitle is for titles inside section borders
	SectionTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true)

	// BoxStyle is for content boxes
	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorBorder)).
		Padding(0, 1)
)

// =============================================================================
// Breadcrumb Styles
// =============================================================================

var (
	// Breadcrumb is for inactive breadcrumb parts
	Breadcrumb = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// BreadcrumbActive is for the current breadcrumb
	BreadcrumbActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSecondary)).
				Bold(true)
)

// =============================================================================
// Progress & Status Styles
// =============================================================================

var (
	// ProgressFilled is for filled progress bar segments
	ProgressFilled = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSpinner))

	// ProgressEmpty is for empty progress bar segments
	ProgressEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBorder))

	// SpinnerStyle is for loading spinners
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSpinner))

	// BadgeOK is for success badges
	BadgeOK = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)).
		Bold(true)

	// BadgeError is for error badges
	BadgeError = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorDanger)).
			Bold(true)

	// BadgeLoading is for loading badges
	BadgeLoading = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning))
)

// =============================================================================
// Service Status Styles
// =============================================================================

var (
	// ServiceRunning is for running service indicator
	ServiceRunning = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess))

	// ServiceStopped is for stopped service indicator
	ServiceStopped = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning))
)

// =============================================================================
// Table Styles
// =============================================================================

var (
	// Cell is for table cell text
	Cell = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText))

	// CellMuted is for dimmed table cells
	CellMuted = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Method is for install method column (very dim)
	Method = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
)

// =============================================================================
// Style Resolver - Semantic style names to actual styles
// =============================================================================

// Resolve returns a lipgloss style for a semantic style name
func Resolve(styleName string) lipgloss.Style {
	switch styleName {
	case "success":
		return Success
	case "warning":
		return Warning
	case "danger":
		return Danger
	case "special":
		return Special
	case "normal":
		return Normal
	case "muted":
		return Muted
	default:
		return Muted
	}
}

// Render applies a semantic style to a value
func Render(value, styleName string) string {
	return Resolve(styleName).Render(value)
}
