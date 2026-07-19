package installview

import (
	"testing"

	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/components"
)

// makeSection builds a Section with `names` tools installable via
// InstallBrewFormula (and thus also Uninstallable), all reporting
// installed=false unless overridden by the caller.
func makeSection(page int, cat config.ToolCategory, names ...string) Section {
	tools := make([]*config.Tool, len(names))
	for i, n := range names {
		tools[i] = &config.Tool{
			Name:     n,
			Category: cat,
			Method:   config.InstallBrewFormula,
			Formula:  n,
			CheckFn:  checkInstalled(false),
		}
	}
	return Section{Page: page, Category: cat, Tools: tools}
}

func modelWithSections(sections []Section, page int) Model {
	m := Model{
		tabs:       components.Tabs{Labels: PageLabels, Active: page},
		sections:   sections,
		expanded:   map[string]bool{},
		checkCache: map[string]config.CheckResult{},
	}
	m.cursor = firstItemCursorForPage(sections, page)
	return m
}

func TestCursorMoveDownWithinSection(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a", "b", "c"),
	}, 0)
	m.moveCursorDown()
	if m.cursor.item != 1 {
		t.Errorf("cursor.item = %d, want 1", m.cursor.item)
	}
	m.moveCursorDown()
	if m.cursor.item != 2 {
		t.Errorf("cursor.item = %d, want 2", m.cursor.item)
	}
	m.moveCursorDown()
	if m.cursor.item != 2 {
		t.Errorf("cursor.item = %d after bottom, want 2 (no-op at end)", m.cursor.item)
	}
}

func TestCursorMoveDownAcrossSectionsSamePage(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a"),
		makeSection(0, config.CategoryCLITools, "b", "c"),
	}, 0)
	m.moveCursorDown()
	if m.cursor.section != 1 || m.cursor.item != 0 {
		t.Errorf("cursor = (%d, %d), want (1, 0)", m.cursor.section, m.cursor.item)
	}
}

func TestCursorMoveDownDoesNotCrossPageBoundary(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a"),
		makeSection(1, config.CategoryDeploy, "b"),
	}, 0)
	m.moveCursorDown()
	if m.cursor.section != 0 || m.cursor.item != 0 {
		t.Errorf("cursor = (%d, %d), want (0, 0) — must not cross into another page's section",
			m.cursor.section, m.cursor.item)
	}
}

func TestCursorMoveDownSkipsCollapsedSection(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a"),
		makeSection(0, config.CategoryCLITools, "b", "c"),
		makeSection(0, config.CategoryRuntimes, "d"),
	}, 0)
	m.sections[1].Collapsed = true
	m.moveCursorDown()
	if m.cursor.section != 2 || m.cursor.item != 0 {
		t.Errorf("cursor = (%d, %d), want (2, 0) — should skip collapsed section",
			m.cursor.section, m.cursor.item)
	}
}

func TestCursorMoveUpAcrossSectionsSamePage(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a", "b"),
		makeSection(0, config.CategoryCLITools, "c"),
	}, 0)
	m.cursor = cursorPos{section: 1, item: 0}
	m.moveCursorUp()
	if m.cursor.section != 0 || m.cursor.item != 1 {
		t.Errorf("cursor = (%d, %d), want (0, 1)", m.cursor.section, m.cursor.item)
	}
}

func TestCursorMoveUpSkipsCollapsedSection(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a"),
		makeSection(0, config.CategoryCLITools, "b"),
		makeSection(0, config.CategoryRuntimes, "c"),
	}, 0)
	m.sections[1].Collapsed = true
	m.cursor = cursorPos{section: 2, item: 0}
	m.moveCursorUp()
	if m.cursor.section != 0 || m.cursor.item != 0 {
		t.Errorf("cursor = (%d, %d), want (0, 0) — should skip collapsed section",
			m.cursor.section, m.cursor.item)
	}
}

func TestToggleCurrentSection(t *testing.T) {
	m := modelWithSections([]Section{
		makeSection(0, config.CategoryGit, "a"),
	}, 0)
	if m.sections[0].Collapsed {
		t.Fatal("section starts collapsed, want expanded")
	}
	m.toggleCurrentSection()
	if !m.sections[0].Collapsed {
		t.Error("toggle didn't collapse")
	}
	m.toggleCurrentSection()
	if m.sections[0].Collapsed {
		t.Error("second toggle didn't expand")
	}
}

func TestStartInstallNoOpWhenAlreadyInstalled(t *testing.T) {
	tools := []*config.Tool{{
		Name: "x", Category: config.CategoryGit, Method: config.InstallBrewFormula,
		CheckFn: checkInstalled(true),
	}}
	m := modelWithSections([]Section{{Page: 0, Category: config.CategoryGit, Tools: tools}}, 0)
	_, cmd := m.startInstall()
	if cmd != nil {
		t.Error("expected nil cmd when already installed")
	}
}

func TestStartInstallNoOpWhenNotInstallable(t *testing.T) {
	tools := []*config.Tool{{
		Name: "xcode", Category: config.CategoryEditorsIDEs, Method: config.InstallMAS,
		CheckFn: checkInstalled(false),
	}}
	m := modelWithSections([]Section{{Page: 0, Category: config.CategoryEditorsIDEs, Tools: tools}}, 0)
	_, cmd := m.startInstall()
	if cmd != nil {
		t.Error("expected nil cmd for a read-only (non-Installable) tool")
	}
}

func TestStartInstallDependencyGate(t *testing.T) {
	// "homebrew" resolves against the real catalog (config.GetToolByName),
	// so dependencyInstalled finds a non-nil dep and then consults
	// checkCache first — seeding checkCache here keeps the test
	// deterministic regardless of whether homebrew is actually installed
	// on the machine running the test.
	tools := []*config.Tool{{
		Name: "needs-brew", Category: config.CategoryGit, Method: config.InstallBrewFormula,
		CheckFn:      checkInstalled(false),
		Dependencies: []string{"homebrew"},
	}}
	m := modelWithSections([]Section{{Page: 0, Category: config.CategoryGit, Tools: tools}}, 0)
	m.checkCache["homebrew"] = config.CheckResult{Installed: false}

	updated, cmd := m.startInstall()
	if cmd != nil {
		t.Error("expected nil cmd when a dependency isn't installed")
	}
	mm := updated.(Model)
	if mm.lastErr == nil || mm.lastErr.Error() != "install homebrew first" {
		t.Errorf("lastErr = %v, want %q", mm.lastErr, "install homebrew first")
	}
}

func TestStartUninstallNoOpWhenNotInstalled(t *testing.T) {
	tools := []*config.Tool{{
		Name: "x", Category: config.CategoryGit, Method: config.InstallBrewFormula,
		CheckFn: checkInstalled(false),
	}}
	m := modelWithSections([]Section{{Page: 0, Category: config.CategoryGit, Tools: tools}}, 0)
	_, cmd := m.startUninstall()
	if cmd != nil {
		t.Error("expected nil cmd when not installed")
	}
}

func TestStartUninstallNoOpWhenNotUninstallable(t *testing.T) {
	tools := []*config.Tool{{
		Name: "xcode", Category: config.CategoryEditorsIDEs, Method: config.InstallMAS,
		CheckFn: checkInstalled(true),
	}}
	m := modelWithSections([]Section{{Page: 0, Category: config.CategoryEditorsIDEs, Tools: tools}}, 0)
	_, cmd := m.startUninstall()
	if cmd != nil {
		t.Error("expected nil cmd for a non-Uninstallable tool")
	}
}

func TestStartUninstallOpensConfirmModal(t *testing.T) {
	tools := []*config.Tool{{
		Name: "x", Category: config.CategoryGit, Method: config.InstallBrewFormula, Formula: "x",
		CheckFn: checkInstalled(true),
	}}
	m := modelWithSections([]Section{{Page: 0, Category: config.CategoryGit, Tools: tools}}, 0)
	updated, cmd := m.startUninstall()
	if cmd == nil {
		t.Fatal("expected modal init cmd")
	}
	mm := updated.(Model)
	if !mm.modalActive() {
		t.Error("modal should be active after startUninstall on an installed, uninstallable tool")
	}
}

func TestRebuildSectionsClampsCursorToPage(t *testing.T) {
	m := Model{
		tabs:     components.Tabs{Labels: PageLabels, Active: 0},
		sections: []Section{makeSection(0, config.CategoryGit, "a", "b", "c")},
		cursor:   cursorPos{section: 0, item: 5},
		expanded: map[string]bool{},
	}
	m.clampCursor()
	if m.cursor.item != 2 {
		t.Errorf("cursor.item = %d, want 2 (clamped to last)", m.cursor.item)
	}
}

func TestClampCursorJumpsToNewPage(t *testing.T) {
	m := Model{
		tabs: components.Tabs{Labels: PageLabels, Active: 0},
		sections: []Section{
			makeSection(0, config.CategoryGit, "a"),
			makeSection(1, config.CategoryDeploy, "b"),
		},
		cursor:   cursorPos{section: 0, item: 0},
		expanded: map[string]bool{},
	}
	m.tabs.Active = 1
	m.clampCursor()
	if m.cursor.section != 1 || m.cursor.item != 0 {
		t.Errorf("cursor = (%d, %d), want (1, 0) after switching to page 1", m.cursor.section, m.cursor.item)
	}
}
