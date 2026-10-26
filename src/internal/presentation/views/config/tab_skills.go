package configview

import (
	"fmt"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/studio-workstation/src/internal/config"
	"github.com/jterrazz/studio-workstation/src/internal/domain/skill"
)

type skillEntry struct {
	Repo string
	Name string
}

type skillSection struct {
	Title     string
	Items     []skillEntry
	Collapsed bool
}

type skillState int

const (
	skillStateNotInstalled skillState = iota
	// invariant: an installed skill whose currency is unknown reads installed: staleness
	// is never invented.
	skillStateInstalled
	skillStateOutdated
	skillStateChecking
)

func (s skillSection) installedCount(states map[string]skillState) (n, total int) {
	total = len(s.Items)
	for _, e := range s.Items {
		if states[e.Name] != skillStateNotInstalled {
			n++
		}
	}
	return
}

func skillsAvailable() bool {
	return skill.IsInstalled()
}

// invariant: no network I/O here; callers fire startSkillCurrencyCheck to resolve the
// checking entries.
func (m *Model) refreshSkillSections() {
	if !skillsAvailable() {
		m.skillSections = nil
		m.skillStates = nil
		return
	}

	installed := skill.ListInstalled()
	lock := skill.ReadLock()

	m.skillStates = make(map[string]skillState, len(installed))
	for _, name := range installed {
		if entry, ok := lock[name]; ok && entry.SourceType == "github" {
			m.skillStates[name] = skillStateChecking
		} else {
			m.skillStates[name] = skillStateInstalled
		}
	}

	m.skillSections = []skillSection{
		{Title: "Studio", Items: skillsToEntries(config.GetStudioSkills())},
		{Title: "Community", Items: skillsToEntries(config.GetCommunitySkills())},
	}

	var others []skillEntry
	for _, name := range installed {
		if !config.IsFavoriteSkill("", name) {
			others = append(others, skillEntry{Name: name})
		}
	}
	if len(others) > 0 {
		m.skillSections = append(m.skillSections, skillSection{Title: "Installed", Items: others})
	}

	m.clampSkillCursor()
}

func skillsToEntries(in []config.Skill) []skillEntry {
	out := make([]skillEntry, len(in))
	for i, s := range in {
		out[i] = skillEntry{Repo: s.Repo, Name: s.Skill}
	}
	return out
}

func (m Model) currentSkill() (skillEntry, bool) {
	if m.skillCursor.section < 0 || m.skillCursor.section >= len(m.skillSections) {
		return skillEntry{}, false
	}
	sec := m.skillSections[m.skillCursor.section]
	if m.skillCursor.item < 0 || m.skillCursor.item >= len(sec.Items) {
		return skillEntry{}, false
	}
	return sec.Items[m.skillCursor.item], true
}

func (m Model) skillInstalled(e skillEntry) bool {
	return m.skillStates[e.Name] != skillStateNotInstalled
}

func (m Model) skillOutdated(e skillEntry) bool {
	return m.skillStates[e.Name] == skillStateOutdated
}

func (m Model) skillsCheckingUpdates() bool {
	for _, st := range m.skillStates {
		if st == skillStateChecking {
			return true
		}
	}
	return false
}

func (m *Model) clampSkillCursor() {
	if len(m.skillSections) == 0 {
		m.skillCursor = cursorPos{section: -1, item: -1}
		return
	}
	if m.skillCursor.section < 0 {
		m.skillCursor.section = 0
	}
	if m.skillCursor.section >= len(m.skillSections) {
		m.skillCursor.section = len(m.skillSections) - 1
	}
	sec := m.skillSections[m.skillCursor.section]
	if len(sec.Items) == 0 {
		m.skillCursor = firstSkillCursor(m.skillSections)
		return
	}
	if m.skillCursor.item >= len(sec.Items) {
		m.skillCursor.item = len(sec.Items) - 1
	}
	if m.skillCursor.item < 0 {
		m.skillCursor.item = 0
	}
}

func firstSkillCursor(sections []skillSection) cursorPos {
	for s := range sections {
		if len(sections[s].Items) > 0 {
			return cursorPos{section: s, item: 0}
		}
	}
	return cursorPos{section: -1, item: -1}
}

func (m *Model) moveSkillCursorUp() {
	if m.skillCursor.section < 0 {
		return
	}
	if !m.skillSections[m.skillCursor.section].Collapsed && m.skillCursor.item > 0 {
		m.skillCursor.item--
		return
	}
	for s := m.skillCursor.section - 1; s >= 0; s-- {
		if m.skillSections[s].Collapsed || len(m.skillSections[s].Items) == 0 {
			continue
		}
		m.skillCursor.section = s
		m.skillCursor.item = len(m.skillSections[s].Items) - 1
		return
	}
}

func (m *Model) moveSkillCursorDown() {
	if m.skillCursor.section < 0 {
		return
	}
	sec := m.skillSections[m.skillCursor.section]
	if !sec.Collapsed && m.skillCursor.item+1 < len(sec.Items) {
		m.skillCursor.item++
		return
	}
	for s := m.skillCursor.section + 1; s < len(m.skillSections); s++ {
		if m.skillSections[s].Collapsed || len(m.skillSections[s].Items) == 0 {
			continue
		}
		m.skillCursor.section = s
		m.skillCursor.item = 0
		return
	}
}

func (m *Model) toggleCurrentSkillSection() {
	if m.skillCursor.section < 0 || m.skillCursor.section >= len(m.skillSections) {
		return
	}
	m.skillSections[m.skillCursor.section].Collapsed = !m.skillSections[m.skillCursor.section].Collapsed
}

func (m Model) skillStartInstall() (tea.Model, tea.Cmd) {
	e, ok := m.currentSkill()
	if !ok {
		return m, nil
	}

	switch m.skillStates[e.Name] {
	case skillStateOutdated:
		m.busy = true
		m.busyAction = "update " + e.Name
		m.lastResult = ""
		m.lastErr = nil
		name := e.Name
		return m, runAction(name, "update", func() error { return skill.Update(name) })

	case skillStateNotInstalled:
		if e.Repo == "" {
			return m, nil
		}
		m.busy = true
		m.busyAction = "install " + e.Name
		m.lastResult = ""
		m.lastErr = nil
		repo, name := e.Repo, e.Name
		return m, runAction(name, "install", func() error { return skill.Install(repo, name) })

	default:
		return m, nil
	}
}

func (m Model) skillStartUninstall() (tea.Model, tea.Cmd) {
	e, ok := m.currentSkill()
	if !ok {
		return m, nil
	}
	if !m.skillInstalled(e) {
		return m, nil
	}
	m.busy = true
	m.busyAction = "uninstall " + e.Name
	m.lastResult = ""
	m.lastErr = nil
	name := e.Name
	return m, runAction(name, "uninstall", func() error { return skill.Remove(name) })
}

func (m Model) renderSkillsBody() string {
	if !skillsAvailable() {
		return contextStyle.Render(" The `skills` CLI isn't installed. Run: npm install -g skills")
	}
	if len(m.skillSections) == 0 {
		return contextStyle.Render(" No skills available.")
	}
	var b strings.Builder
	for sIdx, section := range m.skillSections {
		if sIdx > 0 {
			b.WriteString("\n")
		}
		b.WriteString(m.renderSkillSectionHeader(section))
		b.WriteString("\n")
		if section.Collapsed {
			continue
		}
		for iIdx, entry := range section.Items {
			b.WriteString(m.renderSkillItem(entry, sIdx, iIdx))
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m Model) renderSkillSectionHeader(s skillSection) string {
	caret := "▾"
	if s.Collapsed {
		caret = "▸"
	}
	installed, total := s.installedCount(m.skillStates)
	count := fmt.Sprintf("%d/%d", installed, total)
	return fmt.Sprintf(" %s %s   %s",
		dividerStyle.Render(caret),
		sectionHeaderStyle.Render(s.Title),
		sectionCountStyle.Render(count),
	)
}

func (m Model) renderSkillItem(e skillEntry, sectionIdx, itemIdx int) string {
	isCursor := m.tabs.Active == tabSkills && m.skillCursor.section == sectionIdx && m.skillCursor.item == itemIdx

	state := m.skillStates[e.Name]

	icon := iconMissing
	iconStyle := stateMissingStyle
	suffix := ""
	switch {
	case m.busy && isCursor:
		icon = iconBusy
	case state == skillStateOutdated:
		icon = iconUpdateAvailable
		iconStyle = stateOutdatedStyle
		suffix = "  " + stateOutdatedStyle.Render("update available")
	case state == skillStateInstalled || state == skillStateChecking:
		icon = iconInstalled
		iconStyle = stateInstalledStyle
	}

	cursorMark := "  "
	if isCursor {
		cursorMark = cursorStyle.Render("▶ ")
	}

	nameStyle := itemNameStyle
	if state == skillStateNotInstalled {
		nameStyle = itemNameMutedStyle
	}

	row := fmt.Sprintf(" %s%s %s%s",
		cursorMark,
		iconStyle.Render(icon),
		nameStyle.Render(e.Name),
		suffix,
	)
	if isCursor {
		row = cursorRowStyle.Render(padToWidth(row, m.contentWidth()))
	}
	return row
}

type skillCurrencyMsg struct {
	results map[string]skill.UpdateStatus
}

// value: CheckUpToDate shares one GitHub listing per source repo, so the goroutines
// collapse to one request per distinct repo.
func (m Model) startSkillCurrencyCheck() tea.Cmd {
	if len(m.skillStates) == 0 {
		return nil
	}

	lock := skill.ReadLock()
	type job struct {
		name  string
		entry skill.LockEntry
	}
	var jobs []job
	for name, state := range m.skillStates {
		if state != skillStateChecking {
			continue
		}
		if entry, ok := lock[name]; ok {
			jobs = append(jobs, job{name: name, entry: entry})
		}
	}
	if len(jobs) == 0 {
		return nil
	}

	return func() tea.Msg {
		results := make(map[string]skill.UpdateStatus, len(jobs))
		var mu sync.Mutex
		var wg sync.WaitGroup
		wg.Add(len(jobs))
		for _, j := range jobs {
			go func(j job) {
				defer wg.Done()
				status := skill.CheckUpToDate(j.name, j.entry)
				mu.Lock()
				results[j.name] = status
				mu.Unlock()
			}(j)
		}
		wg.Wait()
		return skillCurrencyMsg{results: results}
	}
}

func (m *Model) applySkillCurrencyResults(results map[string]skill.UpdateStatus) {
	for name, status := range results {
		if _, ok := m.skillStates[name]; !ok {
			continue
		}
		if status == skill.StatusOutdated {
			m.skillStates[name] = skillStateOutdated
		} else {
			m.skillStates[name] = skillStateInstalled
		}
	}
}
