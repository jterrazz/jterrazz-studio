package installview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
)

type keymap struct {
	Up        key.Binding
	Down      key.Binding
	TabPrev   key.Binding
	TabNext   key.Binding
	Toggle    key.Binding
	Details   key.Binding
	Install   key.Binding
	Uninstall key.Binding
	Quit      key.Binding
	Cancel    key.Binding
}

var keys = keymap{
	Up:        key.NewBinding(key.WithKeys("up", "k")),
	Down:      key.NewBinding(key.WithKeys("down", "j")),
	TabPrev:   key.NewBinding(key.WithKeys("left", "shift+tab")),
	TabNext:   key.NewBinding(key.WithKeys("right")),
	Toggle:    key.NewBinding(key.WithKeys("tab")),
	Details:   key.NewBinding(key.WithKeys(" ")),
	Install:   key.NewBinding(key.WithKeys("i", "enter")),
	Uninstall: key.NewBinding(key.WithKeys("u")),
	Quit:      key.NewBinding(key.WithKeys("q", "esc")),
	Cancel:    key.NewBinding(key.WithKeys("ctrl+c")),
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case actionDoneMsg:
		m.busy = false
		m.busyAction = ""
		if msg.err != nil {
			m.lastErr = msg.err
			m.lastResult = fmt.Sprintf("Failed to %s %s: %s", msg.verb, msg.toolName, msg.err)
		} else {
			m.lastErr = nil
			m.lastResult = fmt.Sprintf("%sed %s", capitalize(msg.verb), msg.toolName)
		}
		m.rebuildSections()
		return m, nil
	}

	if m.modalActive() {
		return m.updateModal(msg)
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, keys.Cancel) {
			return m, tea.Quit
		}
		if m.busy {
			return m, nil
		}
		return m.handleKey(keyMsg)
	}
	return m, nil
}

func (m Model) updateModal(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	switch m.form.State {
	case huh.StateCompleted:
		onComplete := m.formOnComplete
		m.closeModal()
		m.lastResult = ""
		m.lastErr = nil
		if onComplete == nil {
			return m, nil
		}
		return m, onComplete()

	case huh.StateAborted:
		m.closeModal()
		return m, nil

	case huh.StateNormal:
	}
	return m, cmd
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if s := msg.String(); len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
		idx := int(s[0] - '1')
		if idx < len(PageLabels) {
			m.tabs.SetActive(idx)
			m.clampCursor()
			m.lastResult = ""
			return m, nil
		}
	}

	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.TabPrev):
		m.tabs.Prev()
		m.clampCursor()
		m.lastResult = ""
		return m, nil

	case key.Matches(msg, keys.TabNext):
		m.tabs.Next()
		m.clampCursor()
		m.lastResult = ""
		return m, nil

	case key.Matches(msg, keys.Up):
		m.moveCursorUp()
		return m, nil

	case key.Matches(msg, keys.Down):
		m.moveCursorDown()
		return m, nil

	case key.Matches(msg, keys.Toggle):
		m.toggleCurrentSection()
		return m, nil

	case key.Matches(msg, keys.Details):
		t := m.currentTool()
		if t != nil {
			m.expanded[t.Name] = !m.expanded[t.Name]
		}
		return m, nil

	case key.Matches(msg, keys.Install):
		return m.startInstall()

	case key.Matches(msg, keys.Uninstall):
		return m.startUninstall()
	}
	return m, nil
}

func (m *Model) toggleCurrentSection() {
	if m.cursor.section < 0 || m.cursor.section >= len(m.sections) {
		return
	}
	m.sections[m.cursor.section].Collapsed = !m.sections[m.cursor.section].Collapsed
}

func (m *Model) moveCursorUp() {
	if m.cursor.section < 0 {
		return
	}
	page := m.tabs.Active
	if !m.sections[m.cursor.section].Collapsed && m.cursor.item > 0 {
		m.cursor.item--
		return
	}
	for s := m.cursor.section - 1; s >= 0; s-- {
		if m.sections[s].Page != page {
			continue
		}
		if m.sections[s].Collapsed || len(m.sections[s].Tools) == 0 {
			continue
		}
		m.cursor.section = s
		m.cursor.item = len(m.sections[s].Tools) - 1
		return
	}
}

func (m *Model) moveCursorDown() {
	if m.cursor.section < 0 {
		return
	}
	page := m.tabs.Active
	sec := m.sections[m.cursor.section]
	if !sec.Collapsed && m.cursor.item+1 < len(sec.Tools) {
		m.cursor.item++
		return
	}
	for s := m.cursor.section + 1; s < len(m.sections); s++ {
		if m.sections[s].Page != page {
			continue
		}
		if m.sections[s].Collapsed || len(m.sections[s].Tools) == 0 {
			continue
		}
		m.cursor.section = s
		m.cursor.item = 0
		return
	}
}

// invariant: post-install scripts run inside the same action, so a script failure
// surfaces as an install failure.
func (m Model) startInstall() (tea.Model, tea.Cmd) {
	t := m.currentTool()
	if t == nil || !t.Installable() {
		return m, nil
	}
	if m.installed(t) {
		return m, nil
	}
	for _, dep := range t.Dependencies {
		if !m.dependencyInstalled(dep) {
			err := fmt.Errorf("install %s first", dep)
			m.lastErr = err
			m.lastResult = err.Error()
			return m, nil
		}
	}
	m.busy = true
	m.busyAction = "install " + t.Name
	m.lastResult = ""
	m.lastErr = nil
	tool := *t
	return m, runAction(tool.Name, "install", func() error {
		if err := tool.Install(); err != nil {
			return err
		}
		return config.RunScripts(tool.Scripts)
	})
}

func (m Model) startUninstall() (tea.Model, tea.Cmd) {
	t := m.currentTool()
	if t == nil || !m.installed(t) || !t.Uninstallable() {
		return m, nil
	}
	tool := *t
	title := fmt.Sprintf("Uninstall %s? This removes it via %s.", tool.Name, tool.Method.String())
	m.buildConfirmModal(title, "uninstall", tool.Name, tool.Uninstall)
	return m, m.form.Init()
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]-32) + s[1:]
}
