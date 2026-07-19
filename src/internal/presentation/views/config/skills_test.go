package configview

import (
	"testing"

	"github.com/jterrazz/jterrazz-studio/src/internal/domain/skill"
)

func modelWithSkills(sections []skillSection, states map[string]skillState) Model {
	m := Model{
		skillSections: sections,
		skillStates:   states,
	}
	m.skillCursor = firstSkillCursor(sections)
	return m
}

func TestSkillInstalledCountsAnyInstalledState(t *testing.T) {
	m := Model{skillStates: map[string]skillState{
		"a": skillStateInstalled,
		"b": skillStateOutdated,
		"c": skillStateChecking,
		"d": skillStateNotInstalled,
	}}
	entries := map[string]skillEntry{
		"a": {Name: "a"}, "b": {Name: "b"}, "c": {Name: "c"}, "d": {Name: "d"},
	}
	if !m.skillInstalled(entries["a"]) {
		t.Error("installed should count as installed")
	}
	if !m.skillInstalled(entries["b"]) {
		t.Error("outdated should count as installed")
	}
	if !m.skillInstalled(entries["c"]) {
		t.Error("checking should count as installed")
	}
	if m.skillInstalled(entries["d"]) {
		t.Error("not-installed should not count as installed")
	}
	// Absent from the map entirely — zero value must also read as not installed.
	if m.skillInstalled(skillEntry{Name: "missing"}) {
		t.Error("absent entry should default to not installed")
	}
}

func TestSkillOutdatedOnlyForOutdatedState(t *testing.T) {
	m := Model{skillStates: map[string]skillState{
		"a": skillStateOutdated,
		"b": skillStateInstalled,
	}}
	if !m.skillOutdated(skillEntry{Name: "a"}) {
		t.Error("expected outdated for a")
	}
	if m.skillOutdated(skillEntry{Name: "b"}) {
		t.Error("expected not outdated for b")
	}
}

func TestSkillsCheckingUpdatesReflectsPendingEntries(t *testing.T) {
	m := Model{skillStates: map[string]skillState{
		"a": skillStateInstalled,
		"b": skillStateChecking,
	}}
	if !m.skillsCheckingUpdates() {
		t.Error("expected true while a checking entry is pending")
	}

	m.skillStates["b"] = skillStateInstalled
	if m.skillsCheckingUpdates() {
		t.Error("expected false once all entries resolved")
	}
}

func TestInstalledCountAnyInstalledStateCounts(t *testing.T) {
	sec := skillSection{Items: []skillEntry{{Name: "a"}, {Name: "b"}, {Name: "c"}}}
	states := map[string]skillState{
		"a": skillStateInstalled,
		"b": skillStateOutdated,
		// "c" absent — not installed.
	}
	n, total := sec.installedCount(states)
	if n != 2 || total != 3 {
		t.Errorf("installedCount = (%d, %d), want (2, 3)", n, total)
	}
}

func TestApplySkillCurrencyResultsResolvesOutdatedAndUpToDate(t *testing.T) {
	m := Model{skillStates: map[string]skillState{
		"outdated-skill": skillStateChecking,
		"current-skill":  skillStateChecking,
		"unknown-skill":  skillStateChecking,
	}}
	m.applySkillCurrencyResults(map[string]skill.UpdateStatus{
		"outdated-skill": skill.StatusOutdated,
		"current-skill":  skill.StatusUpToDate,
		"unknown-skill":  skill.StatusUnknown,
	})

	if m.skillStates["outdated-skill"] != skillStateOutdated {
		t.Errorf("outdated-skill = %v, want skillStateOutdated", m.skillStates["outdated-skill"])
	}
	if m.skillStates["current-skill"] != skillStateInstalled {
		t.Errorf("current-skill = %v, want skillStateInstalled", m.skillStates["current-skill"])
	}
	// StatusUnknown must never be rendered as outdated — settle to installed.
	if m.skillStates["unknown-skill"] != skillStateInstalled {
		t.Errorf("unknown-skill = %v, want skillStateInstalled", m.skillStates["unknown-skill"])
	}
}

func TestApplySkillCurrencyResultsIgnoresSkillsNoLongerTracked(t *testing.T) {
	m := Model{skillStates: map[string]skillState{
		"still-here": skillStateChecking,
	}}
	m.applySkillCurrencyResults(map[string]skill.UpdateStatus{
		"still-here": skill.StatusUpToDate,
		"removed":    skill.StatusOutdated,
	})
	if _, ok := m.skillStates["removed"]; ok {
		t.Error("a skill removed since the check started should not be (re)added")
	}
	if m.skillStates["still-here"] != skillStateInstalled {
		t.Errorf("still-here = %v, want skillStateInstalled", m.skillStates["still-here"])
	}
}

func TestSkillStartInstallInstallsWhenNotInstalled(t *testing.T) {
	m := modelWithSkills(
		[]skillSection{{Items: []skillEntry{{Repo: "owner/repo", Name: "new-skill"}}}},
		map[string]skillState{},
	)
	_, cmd := m.skillStartInstall()
	if cmd == nil {
		t.Fatal("expected install cmd for not-installed skill with a repo")
	}
}

func TestSkillStartInstallNoOpWithoutRepo(t *testing.T) {
	m := modelWithSkills(
		[]skillSection{{Items: []skillEntry{{Name: "orphan"}}}},
		map[string]skillState{},
	)
	_, cmd := m.skillStartInstall()
	if cmd != nil {
		t.Error("expected nil cmd for not-installed orphan with no repo")
	}
}

func TestSkillStartInstallUpdatesWhenOutdated(t *testing.T) {
	m := modelWithSkills(
		[]skillSection{{Items: []skillEntry{{Repo: "owner/repo", Name: "stale-skill"}}}},
		map[string]skillState{"stale-skill": skillStateOutdated},
	)
	updated, cmd := m.skillStartInstall()
	if cmd == nil {
		t.Fatal("expected update cmd for outdated skill")
	}
	mm := updated.(Model)
	if mm.busyAction != "update stale-skill" {
		t.Errorf("busyAction = %q, want %q", mm.busyAction, "update stale-skill")
	}
}

func TestSkillStartInstallNoOpWhenUpToDate(t *testing.T) {
	m := modelWithSkills(
		[]skillSection{{Items: []skillEntry{{Repo: "owner/repo", Name: "current-skill"}}}},
		map[string]skillState{"current-skill": skillStateInstalled},
	)
	_, cmd := m.skillStartInstall()
	if cmd != nil {
		t.Error("expected nil cmd for an already up-to-date skill")
	}
}

func TestSkillStartInstallNoOpWhileChecking(t *testing.T) {
	m := modelWithSkills(
		[]skillSection{{Items: []skillEntry{{Repo: "owner/repo", Name: "checking-skill"}}}},
		map[string]skillState{"checking-skill": skillStateChecking},
	)
	_, cmd := m.skillStartInstall()
	if cmd != nil {
		t.Error("expected nil cmd while currency check is still pending")
	}
}
