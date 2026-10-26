package configview

import (
	"sort"

	"github.com/jterrazz/studio-workstation/src/internal/config"
)

// Section is one category of scripts under a collapsible header of the workstation config TUI.
type Section struct {
	Category  config.ScriptCategory
	Scripts   []*config.Script
	Collapsed bool
}

// invariant: a category not listed here is appended after, alphabetically.
var sectionOrder = []config.ScriptCategory{
	config.ScriptCategoryTerminal,
	config.ScriptCategorySecurity,
	config.ScriptCategoryEditor,
	config.ScriptCategorySystem,
	config.ScriptCategoryServer,
}

func buildSections(role config.Role) []Section {
	return buildSectionsFrom(config.Scripts, role)
}

// invariant: a script without a Category is dropped: every config item declares one.
func buildSectionsFrom(scripts []config.Script, role config.Role) []Section {
	groups := map[config.ScriptCategory][]*config.Script{}
	for i := range scripts {
		s := &scripts[i]
		if !s.MatchesRole(role) {
			continue
		}
		if s.Category == "" {
			continue
		}
		groups[s.Category] = append(groups[s.Category], s)
	}

	sortByName := func(scripts []*config.Script) {
		sort.Slice(scripts, func(i, j int) bool {
			return scripts[i].Name < scripts[j].Name
		})
	}

	var sections []Section
	seen := map[config.ScriptCategory]bool{}
	for _, cat := range sectionOrder {
		if scripts, ok := groups[cat]; ok && len(scripts) > 0 {
			sortByName(scripts)
			sections = append(sections, Section{Category: cat, Scripts: scripts})
			seen[cat] = true
		}
	}

	var extras []config.ScriptCategory
	for cat := range groups {
		if !seen[cat] {
			extras = append(extras, cat)
		}
	}
	sort.Slice(extras, func(i, j int) bool { return extras[i] < extras[j] })
	for _, cat := range extras {
		scripts := groups[cat]
		sortByName(scripts)
		sections = append(sections, Section{Category: cat, Scripts: scripts})
	}

	return sections
}

// invariant: a script without a CheckFn has no state and is not counted.
func (s Section) installedCount() (installed, total int) {
	for _, sc := range s.Scripts {
		if sc.CheckFn == nil {
			continue
		}
		total++
		if sc.CheckFn().Installed {
			installed++
		}
	}
	return
}
