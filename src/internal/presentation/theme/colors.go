package theme

import "github.com/charmbracelet/lipgloss"

// value: the client blue is also the primary accent (titles, selection, spinner,
// progress), so the CLI reads as one product.
const (
	ColorClient = "#5fafd7"
	ColorServer = "#5fd75f"

	ColorPrimary   = ColorClient
	ColorSecondary = "99"
	ColorSuccess   = "42"
	ColorWarning   = "214"
	ColorDanger    = "196"
	ColorMuted     = "241"
	ColorText      = "252"
	ColorSpecial   = "86"
	ColorBorder    = "238"
	ColorSpinner   = ColorClient
)

// Color returns a lipgloss.Color for the given color code
func Color(code string) lipgloss.Color {
	return lipgloss.Color(code)
}
