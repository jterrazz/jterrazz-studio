// Package installview implements the `j install` full-screen TUI. It mirrors
// the architecture of the j config TUI (src/internal/presentation/views/config)
// closely — same checkCache + parallel refresh, same tea.Exec-wrapped action
// pattern, same huh modal machinery — adapted to browse config.Tools instead
// of config.Scripts, and grouped into four pages (Dev/Infra/AI/Apps) instead
// of a single flat list of sections.
package installview

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

// Model is the bubbletea model backing `j install`.
type Model struct {
	tabs components.Tabs

	sections []Section
	cursor   cursorPos
	expanded map[string]bool

	// checkCache stores the result of each Tool's Check() so we don't
	// fork+exec a sub-process (brew, npm, uv, …) on every render. Refreshed
	// at startup and after each install/uninstall completes — never during
	// navigation. Spans every page's tools, not just the active one — there's
	// only one item type here (unlike config's Configuration+Skills tabs).
	checkCache map[string]config.CheckResult

	selfAlias string
	selfRole  config.Role

	width, height int

	busy       bool
	busyAction string

	lastResult string
	lastErr    error

	// Modal state — non-nil while the uninstall-confirmation huh form is up.
	form           *huh.Form
	formTitle      string
	formHelp       string
	formOnComplete func() tea.Cmd
}

// cursorPos identifies the highlighted item: section index (into
// Model.sections) + item index within that section's Tools slice.
type cursorPos struct {
	section int
	item    int
}

// NewModel constructs a fresh Model with all sections expanded, the cursor
// on the first item of the first page, and a fresh check cache.
func NewModel() Model {
	alias, m, ok := config.SelfMachine()
	role := config.Role("")
	if ok {
		role = m.Role
	}
	if !ok {
		alias = "(unregistered)"
	}

	sections := buildSections()
	model := Model{
		tabs:       components.Tabs{Labels: PageLabels, Active: 0},
		sections:   sections,
		cursor:     firstItemCursorForPage(sections, 0),
		expanded:   map[string]bool{},
		checkCache: map[string]config.CheckResult{},
		selfAlias:  alias,
		selfRole:   role,
	}
	model.refreshCheckCache()
	return model
}

// refreshCheckCache invokes every Tool's Check() once (in parallel) and
// stores the results. Called at startup and after each install/uninstall
// lifecycle.
func (m *Model) refreshCheckCache() {
	if m.checkCache == nil {
		m.checkCache = map[string]config.CheckResult{}
	}

	type result struct {
		name string
		cr   config.CheckResult
	}
	var pending []*config.Tool
	for _, sec := range m.sections {
		pending = append(pending, sec.Tools...)
	}
	if len(pending) == 0 {
		return
	}
	results := make([]result, len(pending))
	var wg sync.WaitGroup
	wg.Add(len(pending))
	for i, t := range pending {
		go func(i int, t *config.Tool) {
			defer wg.Done()
			results[i] = result{name: t.Name, cr: t.Check()}
		}(i, t)
	}
	wg.Wait()
	for _, r := range results {
		m.checkCache[r.name] = r.cr
	}
}

// cachedCheck returns the cached CheckResult for t, falling back to a live
// call if the tool isn't in the cache (shouldn't normally happen).
func (m Model) cachedCheck(t *config.Tool) config.CheckResult {
	if t == nil {
		return config.CheckResult{}
	}
	if r, ok := m.checkCache[t.Name]; ok {
		return r
	}
	return t.Check()
}

// installed reports whether the tool is currently installed, reading from
// the cache.
func (m Model) installed(t *config.Tool) bool {
	return m.cachedCheck(t).Installed
}

// dependencyInstalled reports whether the named tool is installed, reading
// from the cache when possible. An unknown dependency name is treated as
// satisfied (mirrors the CLI's installToolByName, which silently skips
// dependencies it can't resolve).
func (m Model) dependencyInstalled(name string) bool {
	dep := config.GetToolByName(name)
	if dep == nil {
		return true
	}
	if r, ok := m.checkCache[dep.Name]; ok {
		return r.Installed
	}
	return dep.Check().Installed
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

// currentTool returns the tool under the cursor, or nil if the cursor isn't
// on a valid item.
func (m Model) currentTool() *config.Tool {
	if m.cursor.section < 0 || m.cursor.section >= len(m.sections) {
		return nil
	}
	sec := m.sections[m.cursor.section]
	if m.cursor.item < 0 || m.cursor.item >= len(sec.Tools) {
		return nil
	}
	return sec.Tools[m.cursor.item]
}

// sectionInstalledCount returns how many of the section's tools are
// currently installed, reading from the check cache.
func (m Model) sectionInstalledCount(sec Section) (installed, total int) {
	total = len(sec.Tools)
	for _, t := range sec.Tools {
		if m.installed(t) {
			installed++
		}
	}
	return
}

// rebuildSections re-runs buildSections (after an install/uninstall changed
// state), refreshes the check cache, and clamps the cursor onto a valid item
// on the current page.
func (m *Model) rebuildSections() {
	m.sections = buildSections()
	m.refreshCheckCache()
	m.clampCursor()
}

// clampCursor ensures the cursor points at a real tool on the current page.
// If it falls off (items moved/vanished, or the page changed), snaps to the
// first item of the active page.
func (m *Model) clampCursor() {
	page := m.tabs.Active
	if m.cursor.section >= 0 && m.cursor.section < len(m.sections) &&
		m.sections[m.cursor.section].Page == page {
		sec := m.sections[m.cursor.section]
		if len(sec.Tools) == 0 {
			m.cursor = firstItemCursorForPage(m.sections, page)
			return
		}
		if m.cursor.item >= len(sec.Tools) {
			m.cursor.item = len(sec.Tools) - 1
		}
		if m.cursor.item < 0 {
			m.cursor.item = 0
		}
		return
	}
	m.cursor = firstItemCursorForPage(m.sections, page)
}

// fnExecCommand adapts a Go func() error into a tea.ExecCommand so install
// and uninstall actions can release the terminal (sudo, prompts, etc. all
// need raw TTY access).
type fnExecCommand struct {
	fn func() error
}

func (f *fnExecCommand) Run() error        { return f.fn() }
func (*fnExecCommand) SetStdin(io.Reader)  {}
func (*fnExecCommand) SetStdout(io.Writer) {}
func (*fnExecCommand) SetStderr(io.Writer) {}

// actionDoneMsg is dispatched when a tea.Exec'd install/uninstall completes.
type actionDoneMsg struct {
	toolName string
	verb     string // "install" or "uninstall"
	err      error
}

// runAction returns a tea.Cmd that releases the terminal, runs fn, then
// emits an actionDoneMsg. Used for both install and uninstall paths.
func runAction(toolName, verb string, fn func() error) tea.Cmd {
	return tea.Exec(&fnExecCommand{fn: fn}, func(err error) tea.Msg {
		return actionDoneMsg{toolName: toolName, verb: verb, err: err}
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
