package installview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/jterrazz-studio/src/internal/config"
	"github.com/jterrazz/jterrazz-studio/src/internal/presentation/print"
)

// State icons for each item row.
const (
	iconInstalled = "✓"
	iconMissing   = "✗"
	iconBusy      = "…"
)

// View implements tea.Model.
func (m Model) View() string {
	if m.modalActive() {
		return m.renderModal()
	}
	var b strings.Builder
	b.WriteString(m.renderHeader("j install"))
	b.WriteString(m.tabs.Render(m.contentWidth()))
	b.WriteString("\n\n")
	b.WriteString(m.renderBody())
	b.WriteString("\n")
	b.WriteString(m.renderDivider())
	b.WriteString("\n")
	b.WriteString(m.renderFooter())
	b.WriteString("\n")
	return b.String()
}

// renderModal renders the uninstall-confirmation form. Frames huh's output
// with the standard header (re-titled to formTitle) so the user keeps
// context, plus a footer hint.
func (m Model) renderModal() string {
	var b strings.Builder
	b.WriteString(m.renderHeader("j install"))
	if m.formHelp != "" {
		b.WriteString(detailTextStyle.Render(" " + wrapText(m.formHelp, m.contentWidth()-2)))
		b.WriteString("\n\n")
	}
	b.WriteString(m.form.View())
	b.WriteString("\n")
	b.WriteString(m.renderDivider())
	b.WriteString("\n")
	b.WriteString(footerLabelStyle.Render(" enter confirm   esc cancel"))
	b.WriteString("\n")
	return b.String()
}

// renderHeader builds the canonical j install header. The right-side context
// shows which machine the TUI is operating on (alias + colour-coded role).
func (m Model) renderHeader(command string) string {
	var context string
	if m.selfRole == "" {
		context = print.MutedText("(unregistered)")
	} else {
		context = m.selfAlias + " " + print.RenderRole(string(m.selfRole))
	}
	return print.RenderHeader(command, context, m.contentWidth())
}

func (m Model) renderDivider() string {
	w := m.contentWidth()
	if w <= 0 {
		w = 80
	}
	return dividerStyle.Render(strings.Repeat("─", w))
}

// renderBody renders the sections belonging to the active page only.
func (m Model) renderBody() string {
	page := m.tabs.Active
	var b strings.Builder
	any := false
	for sIdx, section := range m.sections {
		if section.Page != page {
			continue
		}
		if any {
			b.WriteString("\n")
		}
		any = true
		b.WriteString(m.renderSectionHeader(section))
		b.WriteString("\n")
		if section.Collapsed {
			continue
		}
		for iIdx, t := range section.Tools {
			b.WriteString(m.renderItem(t, sIdx, iIdx))
			b.WriteString("\n")
			if m.expanded[t.Name] {
				b.WriteString(m.renderDetail(t))
				b.WriteString("\n")
			}
		}
	}
	if !any {
		return contextStyle.Render("  No tools on this page.")
	}
	return b.String()
}

func (m Model) renderSectionHeader(s Section) string {
	caret := "▾"
	if s.Collapsed {
		caret = "▸"
	}
	installed, total := m.sectionInstalledCount(s)
	name := sectionHeaderStyle.Render(string(s.Category))
	count := sectionCountStyle.Render(fmt.Sprintf("%d/%d", installed, total))
	return fmt.Sprintf(" %s %s   %s", dividerStyle.Render(caret), name, count)
}

// renderItem renders one tool row: cursor + state icon + name + muted
// description on the left, right-aligned version + method label (plus a
// "(read-only)" marker for entries that are neither Installable nor
// Uninstallable).
func (m Model) renderItem(t *config.Tool, sectionIdx, itemIdx int) string {
	isCursor := m.cursor.section == sectionIdx && m.cursor.item == itemIdx
	installed := m.installed(t)

	icon := iconMissing
	iconStyle := stateMissingStyle
	if m.busy && isCursor {
		icon = iconBusy
	} else if installed {
		icon = iconInstalled
		iconStyle = stateInstalledStyle
	}

	cursorMark := "  "
	if isCursor {
		cursorMark = cursorStyle.Render("▶ ")
	}

	nameStyle := itemNameStyle
	if !installed {
		nameStyle = itemNameMutedStyle
	}

	width := m.contentWidth()
	right := m.rightLabel(t)
	rightW := lipgloss.Width(right)

	prefix := fmt.Sprintf(" %s%s %s", cursorMark, iconStyle.Render(icon), nameStyle.Render(t.Name))
	prefixW := lipgloss.Width(prefix)

	desc := ""
	descBudget := width - prefixW - rightW - 2
	if t.Description != "" && descBudget > 4 {
		text := " — " + t.Description
		if len([]rune(text)) > descBudget {
			text = string([]rune(text)[:descBudget-1]) + "…"
		}
		desc = detailTextStyle.Render(text)
	}

	left := prefix + desc
	leftW := lipgloss.Width(left)
	gap := width - leftW - rightW
	if gap < 1 {
		gap = 1
	}
	row := left + strings.Repeat(" ", gap) + right

	if isCursor {
		row = cursorRowStyle.Render(padToWidth(row, width))
	}
	return row
}

// rightLabel builds the right-aligned "version  method  (read-only)" segment
// of an item row.
func (m Model) rightLabel(t *config.Tool) string {
	var parts []string
	if v := m.displayVersion(t); v != "" {
		parts = append(parts, v)
	}
	parts = append(parts, t.Method.String())
	if !t.Installable() && !t.Uninstallable() {
		parts = append(parts, "(read-only)")
	}
	return sectionCountStyle.Render(strings.Join(parts, "  "))
}

// displayVersion reads the cached version for t and trims it the same way
// the status TUI's renderToolRowCompact does: strip a leading "name " or
// "name/" prefix, then cap the length.
func (m Model) displayVersion(t *config.Tool) string {
	version := m.cachedCheck(t).Version
	if version == "" {
		return ""
	}
	if strings.HasPrefix(version, t.Name+" ") {
		version = strings.TrimPrefix(version, t.Name+" ")
	}
	if strings.HasPrefix(version, t.Name+"/") {
		version = strings.TrimPrefix(version, t.Name+"/")
	}
	if len(version) > 16 {
		version = version[:13] + "..."
	}
	return version
}

// renderDetail renders the expanded detail panel for a tool: description,
// Replaces (if set), method, formula, dependencies, post-install scripts.
func (m Model) renderDetail(t *config.Tool) string {
	var lines []string
	lines = append(lines, "│")
	if t.Description != "" {
		for _, line := range strings.Split(wrapText(t.Description, m.contentWidth()-6), "\n") {
			lines = append(lines, "│ "+detailTextStyle.Render(line))
		}
	}
	if t.Replaces != "" {
		lines = append(lines, "│ "+detailTextStyle.Render("replaces: "+t.Replaces))
	}
	lines = append(lines, "│ "+detailTextStyle.Render("method: "+t.Method.String()))
	if t.Formula != "" {
		lines = append(lines, "│ "+detailTextStyle.Render("formula: "+t.Formula))
	}
	if len(t.Dependencies) > 0 {
		lines = append(lines, "│ "+detailTextStyle.Render("depends on: "+strings.Join(t.Dependencies, ", ")))
	}
	if len(t.Scripts) > 0 {
		lines = append(lines, "│ "+detailTextStyle.Render("post-install: "+strings.Join(t.Scripts, ", ")))
	}
	lines = append(lines, "│")
	return detailFrameStyle.Render(strings.Join(lines, "\n"))
}

// renderFooter shows contextual i/u/space/tab hints for the highlighted row,
// a busy line while an action runs, and the last result/error banner.
func (m Model) renderFooter() string {
	if m.busy {
		return contextStyle.Render(" " + iconBusy + " " + m.busyAction + "…")
	}

	t := m.currentTool()
	var footer string
	if t == nil {
		footer = contextStyle.Render(" no item selected")
	} else {
		var hints []string
		installed := m.installed(t)
		if !installed && t.Installable() {
			hints = append(hints, footerKey("i", "install"))
		}
		if installed && t.Uninstallable() {
			hints = append(hints, footerKey("u", "uninstall"))
		}
		detailLabel := "details"
		if m.expanded[t.Name] {
			detailLabel = "close"
		}
		hints = append(hints, footerKey("space", detailLabel))
		hints = append(hints, footerKey("tab", "fold"))
		hints = append(hints, footerKey("q", "quit"))

		prefix := footerLabelStyle.Render(" ▶ " + t.Name + "  ")
		footer = prefix + strings.Join(hints, footerSepStyle.Render("   "))
	}

	if m.lastResult != "" {
		footer = m.renderResult() + "\n" + footer
	}
	return footer
}

func (m Model) renderResult() string {
	if m.lastErr != nil {
		return resultErrStyle.Render(" ✗ " + m.lastResult)
	}
	return resultOkStyle.Render(" ✓ " + m.lastResult)
}

func footerKey(k, label string) string {
	return footerKeyStyle.Render(k) + " " + footerLabelStyle.Render(label)
}

func (m Model) contentWidth() int {
	if m.width <= 0 {
		return 80
	}
	return m.width
}

// padToWidth right-pads text with spaces so the visible width matches `w`.
func padToWidth(text string, w int) string {
	current := lipgloss.Width(text)
	if current >= w {
		return text
	}
	return text + strings.Repeat(" ", w-current)
}

// wrapText word-wraps `s` to lines no wider than `width`. Preserves any
// pre-existing newlines.
func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	var out []string
	for _, paragraph := range strings.Split(s, "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
		var line strings.Builder
		for _, w := range words {
			if line.Len() == 0 {
				line.WriteString(w)
				continue
			}
			if line.Len()+1+len(w) > width {
				out = append(out, line.String())
				line.Reset()
				line.WriteString(w)
				continue
			}
			line.WriteString(" ")
			line.WriteString(w)
		}
		if line.Len() > 0 {
			out = append(out, line.String())
		}
	}
	return strings.Join(out, "\n")
}
