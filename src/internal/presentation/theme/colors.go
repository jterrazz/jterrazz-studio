package theme

import "github.com/charmbracelet/lipgloss"

// value: the client blue is also the primary accent (titles, selection, spinner,
// progress), so the CLI reads as one product.
const (
	ColorClient = "#5fafd7" // Cyan blue — client role pill, primary accent
	ColorServer = "#5fd75f" // Green — server role pill

	ColorPrimary   = ColorClient // primary accent (titles, selection, spinner)
	ColorSecondary = "99"        // Purple for breadcrumbs / secondary headers
	ColorSuccess   = "42"        // Green for success/installed
	ColorWarning   = "214"       // Orange for actions
	ColorDanger    = "196"       // Red for errors/not configured
	ColorMuted     = "241"       // Gray for dimmed text
	ColorText      = "252"       // Light gray for normal text
	ColorSpecial   = "86"        // Cyan for special highlights
	ColorBorder    = "238"       // Dark gray for borders
	ColorSpinner   = ColorClient // Loader spinner = primary accent
)

// Color returns a lipgloss.Color for the given color code
func Color(code string) lipgloss.Color {
	return lipgloss.Color(code)
}
