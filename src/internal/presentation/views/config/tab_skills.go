package configview

import (
	"fmt"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/domain/skill"
)

// skillEntry is one row in the Skills tab — a (repo, name) pair plus the
// derived install state. Built by refreshSkillSections from the static
// favourites list and the live `skill list` output.
type skillEntry struct {
	Repo string
	Name string
}

// skillSection groups skillEntries under a collapsible header.
type skillSection struct {
	Title     string
	Items     []skillEntry
	Collapsed bool
}

// skillState is the currency state of one row in the Skills tab. The zero
// value (skillStateNotInstalled) is what a skill gets when it has no entry
// in Model.skillStates, so only installed skills need to be populated.
type skillState int

const (
	// skillStateNotInstalled: absent from ~/.agents/skills.
	skillStateNotInstalled skillState = iota
	// skillStateInstalled: present and (as far as is known) current. Also
	// used for installed skills whose currency can't be determined (no
	// github lock entry, or CheckUpToDate returned StatusUnknown) — we
	// never invent staleness, so "installed" is the safe default.
	skillStateInstalled
	// skillStateOutdated: installed, and CheckUpToDate confirmed the skill's
	// upstream folder moved on since it was installed — i.e. the `skills`
	// CLI has an update to apply.
	skillStateOutdated
	// skillStateChecking: installed, currency check in flight. Renders the
	// same as skillStateInstalled (plain ✓) until resolved.
	skillStateChecking
)

// installedCount counts how many of the section's skills are currently
// installed (any state other than skillStateNotInstalled counts).
func (s skillSection) installedCount(states map[string]skillState) (n, total int) {
	total = len(s.Items)
	for _, e := range s.Items {
		if states[e.Name] != skillStateNotInstalled {
			n++
		}
	}
	return
}

// skillsAvailable reports whether the underlying `skills` CLI is installed.
// When false the Skills tab renders a "not installed" placeholder instead of
// a list — saves a confusing empty pane.
func skillsAvailable() bool {
	return skill.IsInstalled()
}

// refreshSkillSections rebuilds skillSections + skillStates from the static
// favourites list and the filesystem/lock-file truth (skill.ListInstalled +
// skill.ReadLock — fast and offline). Called on startup and after each
// install/update/uninstall completes.
//
// Every installed skill starts in skillStateChecking if it has a github
// lock entry (eligible for a currency check) or skillStateInstalled
// otherwise. Callers are responsible for also firing
// startSkillCurrencyCheck to resolve the checking entries — this method
// never does network I/O itself.
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

	// "Installed" section: anything currently installed that isn't already
	// pinned as a Studio or Community favourite. Repo column blank — we
	// only know the name from `skills list`.
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

// currentSkill returns the entry under the skill cursor, or zero if invalid.
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

// skillInstalled reports whether e is installed, in any currency state
// (installed, outdated, or still-checking all count).
func (m Model) skillInstalled(e skillEntry) bool {
	return m.skillStates[e.Name] != skillStateNotInstalled
}

// skillOutdated reports whether e is installed and confirmed outdated.
func (m Model) skillOutdated(e skillEntry) bool {
	return m.skillStates[e.Name] == skillStateOutdated
}

// skillsCheckingUpdates reports whether a currency check is still in flight
// for any installed skill. Drives the muted footer hint.
func (m Model) skillsCheckingUpdates() bool {
	for _, st := range m.skillStates {
		if st == skillStateChecking {
			return true
		}
	}
	return false
}

// clampSkillCursor mirrors clampCursor for the skills tab.
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

// skillStartInstall handles the `i` key on the Skills tab. Behaviour depends
// on the entry's currency state:
//   - not installed: installs it (no-op if we don't know its repo — an
//     Installed-section orphan);
//   - outdated: updates it in place via `skills update`;
//   - installed / checking: already current as far as we know, no-op.
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

	default: // skillStateInstalled, skillStateChecking — already current
		return m, nil
	}
}

// skillStartUninstall queues an uninstall for the current entry.
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

// renderSkillsBody renders the Skills tab content. Mirrors the layout of
// the Configuration tab (collapsible sections, cursor row, install state
// icons) so the two tabs feel like the same TUI.
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

// skillCurrencyMsg carries the results of an async currency-check fan-out —
// one skill.UpdateStatus per skill that was checked. Delivered as a single
// message so Update() only needs one merge point, however many goroutines
// startSkillCurrencyCheck spawned.
type skillCurrencyMsg struct {
	results map[string]skill.UpdateStatus
}

// startSkillCurrencyCheck returns a tea.Cmd that fans out skill.CheckUpToDate
// over every skill currently in skillStateChecking, in parallel goroutines,
// and delivers one skillCurrencyMsg with the results. Returns nil if there's
// nothing to check (skills CLI unavailable, or no installed skill has a
// github lock entry) — never blocks the caller, since the actual network
// calls happen inside the returned tea.Cmd, off the UI goroutine.
//
// The fan-out is wider than the request count: CheckUpToDate shares one
// GitHub directory listing per source repo across all its skills, so N
// goroutines collapse to one request per distinct repo.
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

// applySkillCurrencyResults merges a skillCurrencyMsg into skillStates,
// resolving each checked skill to skillStateOutdated or skillStateInstalled
// (StatusUnknown settles to installed — we never invent staleness). Skills
// not present in results (e.g. removed since the check started) are left
// untouched.
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
