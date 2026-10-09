package installview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m Model) modalActive() bool {
	return m.form != nil
}

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

func (m *Model) closeModal() {
	m.form = nil
	m.formTitle = ""
	m.formHelp = ""
	m.formOnComplete = nil
}
