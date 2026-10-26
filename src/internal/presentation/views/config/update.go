package configview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/jterrazz/studio-workstation/src/internal/config"
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
	Install:   key.NewBinding(key.WithKeys("i")),
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
			m.lastResult = fmt.Sprintf("Failed to %s %s: %s", msg.verb, msg.scriptName, msg.err)
		} else {
			m.lastErr = nil
			m.lastResult = fmt.Sprintf("%sed %s", capitalize(msg.verb), msg.scriptName)
		}
		// invariant: actionDoneMsg does not carry its tab, so both caches refresh, and the
		// currency check reruns so no stale verdict outlives an install, update or removal.
		m.rebuildSections()
		m.refreshSkillSections()
		return m, m.startSkillCurrencyCheck()

	case skillCurrencyMsg:
		m.applySkillCurrencyResults(msg.results)
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
	// invariant: handled before the per-tab keys, so a stray digit never lands on a list item.
	if s := msg.String(); len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
		idx := int(s[0] - '1')
		if idx < len(TabLabels) {
			m.tabs.SetActive(idx)
			m.lastResult = ""
			return m, nil
		}
	}

	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.TabPrev):
		m.tabs.Prev()
		m.lastResult = ""
		return m, nil

	case key.Matches(msg, keys.TabNext):
		m.tabs.Next()
		m.lastResult = ""
		return m, nil

	case key.Matches(msg, keys.Up):
		m.cursorUpForActiveTab()
		return m, nil

	case key.Matches(msg, keys.Down):
		m.cursorDownForActiveTab()
		return m, nil

	case key.Matches(msg, keys.Toggle):
		m.toggleSectionForActiveTab()
		return m, nil

	case key.Matches(msg, keys.Details):
		s := m.currentScript()
		if s != nil {
			m.expanded[s.Name] = !m.expanded[s.Name]
		}
		return m, nil

	case key.Matches(msg, keys.Install):
		return m.installForActiveTab()

	case key.Matches(msg, keys.Uninstall):
		return m.uninstallForActiveTab()
	}
	return m, nil
}

func (m *Model) cursorUpForActiveTab() {
	switch m.tabs.Active {
	case tabConfiguration:
		m.moveCursorUp()
	case tabSkills:
		m.moveSkillCursorUp()
	}
}

func (m *Model) cursorDownForActiveTab() {
	switch m.tabs.Active {
	case tabConfiguration:
		m.moveCursorDown()
	case tabSkills:
		m.moveSkillCursorDown()
	}
}

func (m *Model) toggleSectionForActiveTab() {
	switch m.tabs.Active {
	case tabConfiguration:
		m.toggleCurrentSection()
	case tabSkills:
		m.toggleCurrentSkillSection()
	}
}

func (m Model) installForActiveTab() (tea.Model, tea.Cmd) {
	switch m.tabs.Active {
	case tabConfiguration:
		return m.startInstall()
	case tabSkills:
		return m.skillStartInstall()
	case tabRemote:
		return m.remoteStartConfigure()
	}
	return m, nil
}

func (m Model) uninstallForActiveTab() (tea.Model, tea.Cmd) {
	switch m.tabs.Active {
	case tabConfiguration:
		return m.startUninstall()
	case tabSkills:
		return m.skillStartUninstall()
	}
	return m, nil
}

func (m Model) startInstall() (tea.Model, tea.Cmd) {
	s := m.currentScript()
	if s == nil || s.InstallFn == nil {
		return m, nil
	}
	if m.installed(s) {
		return m, nil
	}
	if len(s.Inputs) > 0 {
		m.buildModal(s)
		return m, m.form.Init()
	}
	m.busy = true
	m.busyAction = "install " + s.Name
	m.lastResult = ""
	m.lastErr = nil
	install := s.InstallFn
	return m, runAction(s.Name, "install", func() error { return install(config.InputValues{}) })
}

func (m Model) startUninstall() (tea.Model, tea.Cmd) {
	s := m.currentScript()
	if s == nil || s.UninstallFn == nil {
		return m, nil
	}
	if !m.installed(s) {
		return m, nil
	}
	m.busy = true
	m.busyAction = "uninstall " + s.Name
	m.lastResult = ""
	m.lastErr = nil
	return m, runAction(s.Name, "uninstall", s.UninstallFn)
}

// invariant: collapsing keeps the cursor on the section's first item, so expanding
// restores the view.
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
	if !m.sections[m.cursor.section].Collapsed && m.cursor.item > 0 {
		m.cursor.item--
		return
	}
	for s := m.cursor.section - 1; s >= 0; s-- {
		if m.sections[s].Collapsed || len(m.sections[s].Scripts) == 0 {
			continue
		}
		m.cursor.section = s
		m.cursor.item = len(m.sections[s].Scripts) - 1
		return
	}
}

func (m *Model) moveCursorDown() {
	if m.cursor.section < 0 {
		return
	}
	sec := m.sections[m.cursor.section]
	if !sec.Collapsed && m.cursor.item+1 < len(sec.Scripts) {
		m.cursor.item++
		return
	}
	for s := m.cursor.section + 1; s < len(m.sections); s++ {
		if m.sections[s].Collapsed || len(m.sections[s].Scripts) == 0 {
			continue
		}
		m.cursor.section = s
		m.cursor.item = 0
		return
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]-32) + s[1:]
}
