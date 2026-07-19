package installview

import (
	"testing"

	"github.com/jterrazz/jterrazz-studio/src/internal/config"
)

func checkInstalled(b bool) func() config.CheckResult {
	return func() config.CheckResult { return config.CheckResult{Installed: b} }
}

// TestPageFamiliesCoverAllCategoriesExactlyOnce guards the routing table
// itself: every config.ToolCategory must appear in exactly one page's family
// list, with no duplicates and no gaps.
func TestPageFamiliesCoverAllCategoriesExactlyOnce(t *testing.T) {
	want := map[config.ToolCategory]bool{}
	for _, cat := range config.ToolCategories {
		want[cat] = true
	}

	seen := map[config.ToolCategory]int{}
	for _, cats := range pageFamilies {
		for _, cat := range cats {
			seen[cat]++
		}
	}

	if len(want) != 18 {
		t.Fatalf("config.ToolCategories has %d entries, want 18 (test assumption)", len(want))
	}
	for cat := range want {
		if seen[cat] != 1 {
			t.Errorf("category %q appears %d times across pageFamilies, want exactly 1", cat, seen[cat])
		}
	}
	for cat, n := range seen {
		if !want[cat] {
			t.Errorf("pageFamilies references unknown category %q (n=%d)", cat, n)
		}
	}
}

// TestBuildSectionsCoversAllFamiliesExactlyOnce runs the grouping against the
// real registry (config.Tools) and checks every one of the 18 families
// produces exactly one section — i.e. the catalog has at least one tool per
// family and buildSectionsFrom doesn't duplicate or drop any.
func TestBuildSectionsCoversAllFamiliesExactlyOnce(t *testing.T) {
	sections := buildSections()
	seen := map[config.ToolCategory]int{}
	for _, s := range sections {
		seen[s.Category]++
	}
	for _, cat := range config.ToolCategories {
		if seen[cat] != 1 {
			t.Errorf("category %q produced %d sections, want 1", cat, seen[cat])
		}
	}
}

func TestBuildSectionsFromAssignsPage(t *testing.T) {
	tools := []config.Tool{
		{Name: "a", Category: config.CategoryGit},           // Dev
		{Name: "b", Category: config.CategoryDeploy},        // Infra
		{Name: "c", Category: config.CategoryAIAgents},      // AI
		{Name: "d", Category: config.CategoryCommunication}, // Apps
	}
	sections := buildSectionsFrom(tools)
	wantPage := map[config.ToolCategory]int{
		config.CategoryGit:           0,
		config.CategoryDeploy:        1,
		config.CategoryAIAgents:      2,
		config.CategoryCommunication: 3,
	}
	if len(sections) != len(wantPage) {
		t.Fatalf("got %d sections, want %d", len(sections), len(wantPage))
	}
	for _, s := range sections {
		if s.Page != wantPage[s.Category] {
			t.Errorf("section %s has Page=%d, want %d", s.Category, s.Page, wantPage[s.Category])
		}
	}
}

func TestBuildSectionsFromPreservesRegistryOrder(t *testing.T) {
	tools := []config.Tool{
		{Name: "zeta", Category: config.CategoryCLITools},
		{Name: "alpha", Category: config.CategoryCLITools},
		{Name: "mid", Category: config.CategoryCLITools},
	}
	sections := buildSectionsFrom(tools)
	if len(sections) != 1 {
		t.Fatalf("got %d sections, want 1", len(sections))
	}
	got := []string{sections[0].Tools[0].Name, sections[0].Tools[1].Name, sections[0].Tools[2].Name}
	want := []string{"zeta", "alpha", "mid"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("tools[%d] = %q, want %q (registry order, not sorted)", i, got[i], want[i])
		}
	}
}

func TestSectionInstalledCount(t *testing.T) {
	tools := []config.Tool{
		{Name: "a", Category: config.CategoryGit, CheckFn: checkInstalled(true)},
		{Name: "b", Category: config.CategoryGit, CheckFn: checkInstalled(false)},
	}
	sections := buildSectionsFrom(tools)
	installed, total := sections[0].installedCount()
	if installed != 1 || total != 2 {
		t.Errorf("got %d/%d, want 1/2", installed, total)
	}
}

func TestFirstItemCursorForPage(t *testing.T) {
	sections := []Section{
		{Page: 0, Category: config.CategoryGit, Tools: nil},
		{Page: 1, Category: config.CategoryDeploy, Tools: []*config.Tool{{Name: "x"}}},
	}
	got := firstItemCursorForPage(sections, 1)
	if got.section != 1 || got.item != 0 {
		t.Errorf("cursor = (%d,%d), want (1,0)", got.section, got.item)
	}
	if empty := firstItemCursorForPage(sections, 0); empty.section != -1 {
		t.Errorf("page 0 has no tools, want (-1,-1), got (%d,%d)", empty.section, empty.item)
	}
}
