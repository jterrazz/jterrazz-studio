package installview

import (
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
)

// PageLabels lists the j install pages in display order; an index here is Model.tabs.Active and pageFamilies'.
var PageLabels = []string{"Dev", "Infra", "AI", "Apps"}

// invariant: every config.ToolCategory appears in exactly one page's list
// (TestPageFamiliesCoverAllCategoriesExactlyOnce).
var pageFamilies = [][]config.ToolCategory{
	{
		config.CategoryPackageManager,
		config.CategoryRuntimes,
		config.CategoryShellTerminal,
		config.CategoryCLITools,
		config.CategoryGit,
		config.CategoryEditorsIDEs,
	},
	{
		config.CategoryContainersVMs,
		config.CategoryDeploy,
		config.CategoryRemoteAccess,
		config.CategorySecurity,
	},
	{
		config.CategoryAIAgents,
		config.CategoryAITooling,
		config.CategoryAIApps,
	},
	{
		config.CategoryBrowsers,
		config.CategoryCommunication,
		config.CategoryProductivity,
		config.CategoryMedia,
		config.CategorySystemUtilities,
	},
}

// Section is one family of tools under a collapsible header, in registry order, never re-sorted.
type Section struct {
	Page      int
	Category  config.ToolCategory
	Tools     []*config.Tool
	Collapsed bool
}

func buildSections() []Section {
	return buildSectionsFrom(config.Tools)
}

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

// value: the running TUI reads Model's checkCache instead; this live count is the tests'.
func (s Section) installedCount() (installed, total int) {
	total = len(s.Tools)
	for _, t := range s.Tools {
		if t.Check().Installed {
			installed++
		}
	}
	return
}

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
