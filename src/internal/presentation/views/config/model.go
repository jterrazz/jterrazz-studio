package configview

import (
	"fmt"
	"io"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/components"
)

// TabLabels lists the j config tabs in display order; an index here is Model.tabs.Active.
var TabLabels = []string{"Configuration", "Skills", "Remote"}

const (
	tabConfiguration = 0
	tabSkills        = 1
	tabRemote        = 2
)

// Model is the Bubble Tea model behind j config: the Configuration, Skills and Remote tabs.
type Model struct {
	tabs components.Tabs

	sections []Section
	cursor   cursorPos
	expanded map[string]bool

	// value: a CheckFn forks a sub-process (defaults, pmset, systemsetup), so its results
	// are cached and refreshed only at startup and after an install or uninstall.
	checkCache map[string]config.CheckResult

	skillSections []skillSection
	skillCursor   cursorPos

	// invariant: a skill with no entry is skillStateNotInstalled; checking entries are
	// resolved by the async currency check fired from Init and after every change.
	skillStates map[string]skillState

	selfAlias string
	selfRole  config.Role

	width, height int

	busy       bool
	busyAction string

	lastResult string
	lastErr    error

	form           *huh.Form
	formScript     *config.Script
	formBindings   []*string
	formTitle      string         // header label, e.g. "install autologin"
	formHelp       string         // optional sub-text shown above the form
	formOnComplete func() tea.Cmd // fires when the user submits
}

type cursorPos struct {
	section int
	item    int
}

func NewModel() Model {
	alias, m, ok := config.SelfMachine()
	role := config.Role("")
	if ok {
		role = m.Role
	}
	if !ok {
		alias = "(unregistered)"
	}

	sections := buildSections(role)
	model := Model{
		tabs:       components.Tabs{Labels: TabLabels, Active: tabConfiguration},
		sections:   sections,
		cursor:     firstItemCursor(sections),
		expanded:   map[string]bool{},
		checkCache: map[string]config.CheckResult{},
		selfAlias:  alias,
		selfRole:   role,
	}
	model.refreshCheckCache()
	model.refreshSkillSections()
	model.skillCursor = firstSkillCursor(model.skillSections)
	return model
}

func (m *Model) refreshCheckCache() {
	if m.checkCache == nil {
		m.checkCache = map[string]config.CheckResult{}
	}

	type result struct {
		name string
		cr   config.CheckResult
	}
	var pending []*config.Script
	for _, sec := range m.sections {
		for _, s := range sec.Scripts {
			if s.CheckFn != nil {
				pending = append(pending, s)
			}
		}
	}
	if len(pending) == 0 {
		return
	}
	results := make([]result, len(pending))
	var wg sync.WaitGroup
	wg.Add(len(pending))
	for i, s := range pending {
		go func(i int, s *config.Script) {
			defer wg.Done()
			results[i] = result{name: s.Name, cr: s.CheckFn()}
		}(i, s)
	}
	wg.Wait()
	for _, r := range results {
		m.checkCache[r.name] = r.cr
	}
}

func (m Model) cachedCheck(s *config.Script) config.CheckResult {
	if s == nil {
		return config.CheckResult{}
	}
	if r, ok := m.checkCache[s.Name]; ok {
		return r
	}
	if s.CheckFn != nil {
		return s.CheckFn()
	}
	return config.CheckResult{}
}

func (m Model) installed(s *config.Script) bool {
	return m.cachedCheck(s).Installed
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.WindowSize(), m.startSkillCurrencyCheck())
}

func firstItemCursor(sections []Section) cursorPos {
	for s := range sections {
		if len(sections[s].Scripts) > 0 {
			return cursorPos{section: s, item: 0}
		}
	}
	return cursorPos{section: -1, item: -1}
}

func (m Model) currentScript() *config.Script {
	if m.cursor.section < 0 || m.cursor.section >= len(m.sections) {
		return nil
	}
	sec := m.sections[m.cursor.section]
	if m.cursor.item < 0 || m.cursor.item >= len(sec.Scripts) {
		return nil
	}
	return sec.Scripts[m.cursor.item]
}

func (m Model) sectionInstalledCount(sec Section) (installed, total int) {
	for _, sc := range sec.Scripts {
		if sc.CheckFn == nil {
			continue
		}
		total++
		if m.installed(sc) {
			installed++
		}
	}
	return
}

func (m *Model) rebuildSections() {
	m.sections = buildSections(m.selfRole)
	m.refreshCheckCache()
	m.clampCursor()
}

func (m *Model) clampCursor() {
	if len(m.sections) == 0 {
		m.cursor = cursorPos{section: -1, item: -1}
		return
	}
	if m.cursor.section < 0 {
		m.cursor.section = 0
	}
	if m.cursor.section >= len(m.sections) {
		m.cursor.section = len(m.sections) - 1
	}
	sec := m.sections[m.cursor.section]
	if len(sec.Scripts) == 0 {
		m.cursor = firstItemCursor(m.sections)
		return
	}
	if m.cursor.item >= len(sec.Scripts) {
		m.cursor.item = len(sec.Scripts) - 1
	}
	if m.cursor.item < 0 {
		m.cursor.item = 0
	}
}

// outside: sudo, prompts and key generation need the raw TTY, so an action runs as a tea.ExecCommand.
type fnExecCommand struct {
	fn func() error
}

func (f *fnExecCommand) Run() error        { return f.fn() }
func (*fnExecCommand) SetStdin(io.Reader)  {}
func (*fnExecCommand) SetStdout(io.Writer) {}
func (*fnExecCommand) SetStderr(io.Writer) {}

type actionDoneMsg struct {
	scriptName string
	verb       string
	err        error
}

func runAction(scriptName, verb string, fn func() error) tea.Cmd {
	return tea.Exec(&fnExecCommand{fn: fn}, func(err error) tea.Msg {
		return actionDoneMsg{scriptName: scriptName, verb: verb, err: err}
	})
}

// Run starts the TUI. Use RunOrExit for the standard CLI entry point.
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// RunOrExit runs the TUI and exits with status 1 on error.
func RunOrExit() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
