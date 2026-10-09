package components

import (
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/theme"
)

// Tabs is the one-line tab strip every tabbed TUI shares, the active label in the primary accent.
type Tabs struct {
	Labels []string
	Active int
}

// Render returns the strip as one line; it spans only as far as its labels.
func (t Tabs) Render(width int) string {
	_ = width
	var b strings.Builder
	b.WriteString(" ")
	for i, label := range t.Labels {
		if i > 0 {
			b.WriteString("  ")
		}
		if i == t.Active {
			b.WriteString(theme.Selected.Render("● " + label))
		} else {
			b.WriteString(theme.Muted.Render("○ " + label))
		}
	}
	return b.String()
}

// Next moves the active tab right. No-op at the end.
func (t *Tabs) Next() {
	if t.Active < len(t.Labels)-1 {
		t.Active++
	}
}

// Prev moves the active tab left. No-op at the start.
func (t *Tabs) Prev() {
	if t.Active > 0 {
		t.Active--
	}
}

// SetActive jumps directly to an index, clamped to bounds.
func (t *Tabs) SetActive(i int) {
	if i < 0 {
		i = 0
	}
	if i >= len(t.Labels) {
		i = len(t.Labels) - 1
	}
	t.Active = i
}
