package installview

import (
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
)

// PageLabels enumerates the j install pages in display order. Index in this
// slice maps to Model.tabs.Active and to the index into pageFamilies below.
var PageLabels = []string{"Dev", "Infra", "AI", "Apps"}

// pageFamilies groups the 18 ToolCategory families under each page, in the
// order they render within that page. Every config.ToolCategory must appear
// in exactly one page's list — see TestBuildSectionsCoversAllFamiliesExactlyOnce.
var pageFamilies = [][]config.ToolCategory{
	{ // Dev
		config.CategoryPackageManager,
		config.CategoryRuntimes,
		config.CategoryShellTerminal,
		config.CategoryCLITools,
		config.CategoryGit,
		config.CategoryEditorsIDEs,
	},
	{ // Infra
		config.CategoryContainersVMs,
		config.CategoryDeploy,
		config.CategoryRemoteAccess,
		config.CategorySecurity,
	},
	{ // AI
		config.CategoryAIAgents,
		config.CategoryAITooling,
		config.CategoryAIApps,
	},
	{ // Apps
		config.CategoryBrowsers,
		config.CategoryCommunication,
		config.CategoryProductivity,
		config.CategoryMedia,
		config.CategorySystemUtilities,
	},
}

// Section is a family of tools grouped under one collapsible header on a
// page. Order within a page mirrors pageFamilies; order within a section
// mirrors the registry (config.Tools) — no re-sorting, unlike the config
// TUI's alphabetical sections.
type Section struct {
	Page      int
	Category  config.ToolCategory
	Tools     []*config.Tool
	Collapsed bool
}

// buildSections groups config.Tools by family into pages, in canonical order.
func buildSections() []Section {
	return buildSectionsFrom(config.Tools)
}

// buildSectionsFrom is the testable variant: takes the tool list explicitly
// instead of pulling from the package-level config.Tools.
func buildSectionsFrom(tools []config.Tool) []Section {
	groups := map[config.ToolCategory][]*config.Tool{}
	for i := range tools {
		t := &tools[i]
		groups[t.Category] = append(groups[t.Category], t)
	}

	var sections []Section
	for page, cats := range pageFamilies {
		for _, cat := range cats {
			items, ok := groups[cat]
			if !ok || len(items) == 0 {
				continue
			}
			sections = append(sections, Section{Page: page, Category: cat, Tools: items})
		}
	}
	return sections
}

// installedCount returns how many of the section's tools are currently
// installed, running each Tool's live Check(). Used by tests; the running
// TUI reads from Model's checkCache instead (see Model.sectionInstalledCount)
// to avoid re-forking a subprocess on every render.
func (s Section) installedCount() (installed, total int) {
	total = len(s.Tools)
	for _, t := range s.Tools {
		if t.Check().Installed {
			installed++
		}
	}
	return
}

// firstItemCursorForPage returns the position of the first item on the given
// page, or {-1, -1} if the page has no items.
func firstItemCursorForPage(sections []Section, page int) cursorPos {
	for i, sec := range sections {
		if sec.Page != page {
			continue
		}
		if len(sec.Tools) > 0 {
			return cursorPos{section: i, item: 0}
		}
	}
	return cursorPos{section: -1, item: -1}
}
