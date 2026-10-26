package status

import (
	"github.com/jterrazz/studio-workstation/src/internal/presentation/components"
)

func (m Model) loadingIndicator() string {
	if m.allLoaded {
		return components.Muted("·")
	}
	return m.spinner.View()
}
