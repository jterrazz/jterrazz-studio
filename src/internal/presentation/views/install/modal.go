package installview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// modalActive reports whether the uninstall-confirmation modal is currently
// being shown.
func (m Model) modalActive() bool {
	return m.form != nil
}

// buildConfirmModal opens a Yes/No huh confirm before running a destructive
// action (currently: uninstall). Mirrors configview's buildModal: the
// onComplete callback built here sets busy state and kicks off runAction
// exactly like the direct (non-modal) install path does — but only if the
// user picked "Yes". Declining ("No", or aborting with esc) just closes the
// modal without running anything.
func (m *Model) buildConfirmModal(title, verb, name string, run func() error) {
	confirmed := new(bool)
	m.formTitle = title
	m.formHelp = ""
	m.formOnComplete = func() tea.Cmd {
		if !*confirmed {
			return nil
		}
		m.busy = true
		m.busyAction = verb + " " + name
		return runAction(name, verb, run)
	}
	m.form = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(title).
			Affirmative("Yes").
			Negative("No").
			Value(confirmed),
	)).WithTheme(huh.ThemeBase()).WithShowHelp(false)
}

// closeModal clears modal state. Called on completion (after queueing the
// uninstall action, or after a decline) or on abort (esc).
func (m *Model) closeModal() {
	m.form = nil
	m.formTitle = ""
	m.formHelp = ""
	m.formOnComplete = nil
}
