package components

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/theme"
)

// SpinnerFPS is the animation speed for all spinners.
const SpinnerFPS = 80 * time.Millisecond

// NewSpinnerModel returns the styled spinner every TUI loads with.
func NewSpinnerModel() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: theme.BrailleSpinner,
		FPS:    SpinnerFPS,
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSpinner))
	return s
}
