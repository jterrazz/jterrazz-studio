package status

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jterrazz/studio-workstation/src/internal/config"
	"github.com/jterrazz/studio-workstation/src/internal/domain/status"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/components"
	"github.com/jterrazz/studio-workstation/src/internal/presentation/theme"
)

const minColumnWidth = 44

// TabLabels lists the workstation status tabs in display order; an index here is Model.activeTab.
var TabLabels = []string{"System", "Workspace", "Applications", "Configuration"}

func (m Model) renderContent() string {
	w := m.width
	if w < minColumnWidth {
		w = minColumnWidth
	}

	sections := m.groupBySection()

	switch m.activeTab {
	case 0:
		return m.renderActivity(sections, w) + "\n" + m.renderEnvironment(sections, w)
	case 1:
		return m.renderWorkspace(sections, w)
	case 2:
		return m.renderTools(sections, w)
	case 3:
		return m.renderConfig(sections, w)
	}
	return ""
}

func (m Model) renderTabBar(width int) string {
	return components.Tabs{Labels: TabLabels, Active: m.activeTab}.Render(width)
}

func (m Model) renderActivity(sections sectionMap, w int) string {
	var b strings.Builder

	innerW := w - 4

	sysItems := m.getSubsectionItems(sections, "Environment", "Health")
	var sysRows []string
	for _, item := range sysItems {
		if item.Kind == status.KindProcess {
			for _, p := range item.Processes {
				sysRows = append(sysRows, " "+padName(p.Name, 12)+theme.Muted.Render(p.Value))
			}
		}
	}

	sysColW := innerW / 4
	if sysColW < 24 {
		sysColW = 24
	}
	sparkColW := innerW - sysColW - 1

	// value: a sparkline row is " LABEL " (7) + padding + graph + " VALUE " (8).
	const slLabelW = 7
	const slValueW = 8
	graphW := sparkColW * 40 / 100
	if graphW < 15 {
		graphW = 15
	}
	if graphW > sparkHistorySize {
		graphW = sparkHistorySize
	}
	slPad := sparkColW - slLabelW - graphW - slValueW
	if slPad < 0 {
		slPad = 0
	}

	buildSparkRow := func(label string, history []float64, colored bool) string {
		row := fmt.Sprintf(" %-6s", label)
		if len(history) > 0 {
			var spark string
			if colored {
				spark = renderColoredSparkline(history, graphW)
			} else {
				spark = renderAutoSparkline(history, graphW)
			}
			current := history[len(history)-1]
			var val string
			if colored {
				val = fmt.Sprintf("%5.0f%% ", current)
			} else {
				val = formatBytesPerSec(current) + " "
			}
			row += strings.Repeat(" ", slPad) + spark + theme.Muted.Render(" "+val)
		} else {
			row += strings.Repeat(" ", slPad+graphW) + theme.Muted.Render("      - ")
		}
		return row
	}

	sparkLines := []string{
		buildSparkRow("CPU", m.cpuHistory, true),
		buildSparkRow("GPU", m.gpuHistory, true),
		buildSparkRow("NET ↓", m.netRxHistory, false),
		buildSparkRow("NET ↑", m.netTxHistory, false),
	}

	var activityRows []string
	activityRows = append(activityRows, "")

	if innerW >= minColWidthResponsive*2 {
		divider := theme.SectionBorder.Render("│")

		activityRows = append(activityRows,
			padTo(" "+theme.SubSection.Render("HEALTH"), sysColW)+divider+" "+theme.SubSection.Render("GRAPHS"),
			padTo("", sysColW)+divider)

		maxR := len(sysRows)
		if len(sparkLines) > maxR {
			maxR = len(sparkLines)
		}
		for i := range maxR {
			left := getOr(sysRows, i, "")
			right := getOr(sparkLines, i, "")
			activityRows = append(activityRows, padTo(left, sysColW)+divider+right)
		}
	} else {
		activityRows = append(activityRows, " "+theme.SubSection.Render("HEALTH"), "")
		activityRows = append(activityRows, sysRows...)
		activityRows = append(activityRows, "", " "+theme.SubSection.Render("GRAPHS"), "")
		activityRows = append(activityRows, sparkLines...)
	}

	activityRows = append(activityRows, "")

	b.WriteString(components.SubsectionBox("System", activityRows, w))
	b.WriteString("\n")

	cpuItems := m.getSubsectionItems(sections, "System", "CPU")
	memItems := m.getSubsectionItems(sections, "System", "Memory")

	barColW := innerW / 2
	cpuRows := m.renderProcessBars(cpuItems, barColW-2)
	memRows := m.renderProcessBars(memItems, barColW-2)

	procContent := renderColumnsRaw([]namedColumn{
		{"CPU", cpuRows},
		{"Memory", memRows},
	}, innerW)

	procLines := strings.Split(strings.TrimRight(procContent, "\n"), "\n")
	boxRows := make([]string, 0, len(procLines)+2)
	boxRows = append(boxRows, "")
	boxRows = append(boxRows, procLines...)
	boxRows = append(boxRows, "")

	b.WriteString(components.SubsectionBox("Processes", boxRows, w))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderEnvironment(sections sectionMap, w int) string {
	numCols := 2
	if w < minColWidthResponsive*2 {
		numCols = 1
	}
	colWidth := w / numCols

	netItems := m.getSubsectionItems(sections, "Environment", "Network")
	netInnerW := colWidth - 4
	var netRows []string
	for _, item := range netItems {
		if item.Kind == status.KindNetwork && item.Loaded && item.Available {
			icon := ""
			if (item.Name == "vpn" || item.Name == "dns") && item.Style == "success" {
				icon = components.Badge(true) + " "
			} else if item.Name == "vpn" || item.Name == "dns" {
				icon = theme.Warning.Render("⚠") + " "
			}
			left := " " + item.Name
			right := icon + theme.Muted.Render(item.Value)
			gap := netInnerW - components.VisibleLen(left) - components.VisibleLen(right)
			if gap < 1 {
				gap = 1
			}
			netRows = append(netRows, left+strings.Repeat(" ", gap)+right)
		} else if item.Kind == status.KindNetwork && !item.Loaded {
			netRows = append(netRows, " "+item.Name+" "+m.loadingIndicator())
		}
	}

	tsItems := m.getSubsectionItems(sections, "Environment", "Tailscale")
	tsInnerW := colWidth - 4
	var tsBox string
	for _, item := range tsItems {
		switch {
		case item.Kind == status.KindTailscale && item.Loaded && item.Available && item.TailscaleStatus != nil:
			st := item.TailscaleStatus

			var tsTopRows []string

			// value: a connected tailnet reads neutral: it is a mesh, not privacy.
			statusVal := strings.ToLower(st.BackendState)
			if st.Daemon != "" {
				statusVal += " (" + string(st.Daemon) + ")"
			}
			left := " status"
			right := theme.Muted.Render(statusVal)
			gap := tsInnerW - components.VisibleLen(left) - components.VisibleLen(right)
			if gap < 1 {
				gap = 1
			}
			tsTopRows = append(tsTopRows, left+strings.Repeat(" ", gap)+right)

			if st.IP != "" {
				left = " ip"
				right = theme.Muted.Render(st.IP)
				gap = tsInnerW - components.VisibleLen(left) - components.VisibleLen(right)
				if gap < 1 {
					gap = 1
				}
				tsTopRows = append(tsTopRows, left+strings.Repeat(" ", gap)+right)
			}

			// value: the exit node reads ✓ only when active (a real VPN), ⚠ when none (traffic exposed).
			exitVal := "none"
			exitIcon := theme.Warning.Render("⚠") + " "
			if st.ExitNode != "" {
				exitVal = st.ExitNode
				exitIcon = components.Badge(true) + " "
			}
			left = " exit node"
			right = exitIcon + theme.Muted.Render(exitVal)
			gap = tsInnerW - components.VisibleLen(left) - components.VisibleLen(right)
			if gap < 1 {
				gap = 1
			}
			tsTopRows = append(tsTopRows, left+strings.Repeat(" ", gap)+right)

			var tsPeerRows []string
			if len(st.Peers) > 0 {
				for _, peer := range st.Peers {
					peerName := " " + peer.Name
					peerIP := peer.IP
					onlineLabel := "offline"
					onlineIcon := ""
					if peer.Online {
						onlineLabel = "online"
						onlineIcon = " " + components.Badge(true)
					}

					rightPart := theme.Muted.Render(peerIP) + onlineIcon + " " + theme.Muted.Render(onlineLabel)
					gap = tsInnerW - components.VisibleLen(peerName) - components.VisibleLen(rightPart)
					if gap < 1 {
						gap = 1
					}
					tsPeerRows = append(tsPeerRows, peerName+strings.Repeat(" ", gap)+rightPart)
				}
			} else {
				tsPeerRows = append(tsPeerRows, " "+theme.Muted.Render("no peers"))
			}

			tsBox = components.SubsectionBoxWithSeparator("Tailscale", tsTopRows, tsPeerRows, colWidth)
		case item.Kind == status.KindTailscale && !item.Loaded:
			tsBox = components.SubsectionBox("Tailscale", []string{" " + m.loadingIndicator()}, colWidth)
		case item.Kind == status.KindTailscale && item.Loaded && !item.Available:
			tsBox = components.SubsectionBox("Tailscale", []string{" " + theme.Muted.Render("not available")}, colWidth)
		}
	}

	svcItems := m.getSubsectionItems(sections, "Environment", "Services")
	svcInnerW := colWidth - 4
	var svcRows []string
	for _, item := range svcItems {
		if item.Kind == status.KindProcess {
			for _, p := range item.Processes {
				left := " " + p.Name
				right := theme.Muted.Render(p.Value)
				leftLen := components.VisibleLen(left)
				rightLen := components.VisibleLen(right)
				gap := svcInnerW - leftLen - rightLen
				if gap < 1 {
					gap = 1
				}
				svcRows = append(svcRows, left+strings.Repeat(" ", gap)+right)
			}
		}
	}
	if len(svcRows) == 0 {
		svcRows = append(svcRows, " "+theme.Muted.Render("none"))
	}

	var boxes []string
	if len(netRows) > 0 {
		boxes = append(boxes, components.SubsectionBox("Network", netRows, colWidth))
	}
	if tsBox != "" {
		boxes = append(boxes, tsBox)
	}
	if len(svcRows) > 0 {
		boxes = append(boxes, components.SubsectionBox("Services", svcRows, colWidth))
	}

	diskItems := m.getSubsectionItems(sections, "Environment", "Disk")
	diskReady := len(m.diskSortOrder) > 0 && m.diskMaxSize > 0
	if diskReady {
		orderMap := make(map[string]int, len(m.diskSortOrder))
		for i, id := range m.diskSortOrder {
			orderMap[id] = i
		}
		sort.SliceStable(diskItems, func(i, j int) bool {
			return orderMap[diskItems[i].ID] < orderMap[diskItems[j].ID]
		})
	}
	diskInnerW := colWidth - 4
	barWidth := 15
	valueWidth := 9
	nameWidth := diskInnerW - barWidth - valueWidth - 3
	if nameWidth < 10 {
		nameWidth = 10
	}
	simpleNameWidth := diskInnerW - valueWidth - 2
	if simpleNameWidth < 10 {
		simpleNameWidth = 10
	}
	var diskRows []string
	for _, item := range diskItems {
		if item.Kind != status.KindCache {
			continue
		}
		if !item.Loaded {
			diskRows = append(diskRows, " "+padName(item.Name, simpleNameWidth)+m.loadingIndicator())
		} else if item.Available {
			if diskReady {
				size := parseDisplaySize(item.Value)
				ratio := size / m.diskMaxSize
				if ratio > 1.0 {
					ratio = 1.0
				}
				bar := renderBar(ratio, barWidth)
				diskRows = append(diskRows, " "+fmt.Sprintf("%-*s", nameWidth, item.Name)+" "+bar+" "+theme.Muted.Render(fmt.Sprintf("%*s", valueWidth, item.Value)))
			} else {
				diskRows = append(diskRows, " "+fmt.Sprintf("%-*s", simpleNameWidth, item.Name)+" "+theme.Muted.Render(fmt.Sprintf("%*s", valueWidth, item.Value)))
			}
		}
	}
	if len(diskRows) > 0 {
		boxes = append(boxes, components.SubsectionBox("Disk Usage", diskRows, colWidth))
	}

	if updRows := m.buildCheckRows(m.getSubsectionItems(sections, "Environment", "Updates"), colWidth-4); len(updRows) > 0 {
		boxes = append(boxes, components.SubsectionBox("Updates", updRows, colWidth))
	}

	if len(boxes) == 0 {
		return ""
	}

	if numCols > len(boxes) {
		numCols = len(boxes)
	}

	columns := make([][]string, numCols)
	colHeights := make([]int, numCols)
	for _, box := range boxes {
		shortest := 0
		for i := 1; i < numCols; i++ {
			if colHeights[i] < colHeights[shortest] {
				shortest = i
			}
		}
		columns[shortest] = append(columns[shortest], box)
		colHeights[shortest] += strings.Count(box, "\n") + 1
	}

	var renderedCols []string
	for _, col := range columns {
		renderedCols = append(renderedCols, lipgloss.JoinVertical(lipgloss.Left, col...))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedCols...) + "\n"
}

func (m Model) renderWorkspace(sections sectionMap, w int) string {
	numCols := 2
	if w < minColWidthResponsive*2 {
		numCols = 1
	}
	colWidth := w / numCols
	innerW := colWidth - 4

	repoItems := m.getSubsectionItems(sections, "Workspace", "Repositories")
	repoInnerW := colWidth - 4
	var repoRows []string
	for _, item := range repoItems {
		if item.Kind == status.KindRepository {
			switch {
			case !item.Loaded:
				repoRows = append(repoRows, " "+m.loadingIndicator())
			case !item.Available:
				repoRows = append(repoRows, " "+theme.Muted.Render("no repositories found"))
			default:
				maxNameW := 0
				for _, group := range item.ProjectGroups {
					for _, repo := range group.Repos {
						if len(repo.Name) > maxNameW {
							maxNameW = len(repo.Name)
						}
					}
				}

				for gi, group := range item.ProjectGroups {
					if gi > 0 {
						repoRows = append(repoRows, "")
					}
					repoRows = append(repoRows, " "+theme.SubSection.Render(group.Prefix+"/"))

					for ri, repo := range group.Repos {
						connector := "├─"
						if ri == len(group.Repos)-1 {
							connector = "└─"
						}

						branchStyle := theme.Muted

						var statusStr string
						if repo.Clean {
							statusStr = theme.Success.Render("clean")
						} else {
							statusStr = theme.Warning.Render(fmt.Sprintf("%d changed", repo.ChangeCount))
						}

						left := "   " + theme.Muted.Render(connector) + " " + fmt.Sprintf("%-*s", maxNameW, repo.Name)
						right := statusStr + "  " + branchStyle.Render(repo.Branch)
						gap := repoInnerW - components.VisibleLen(left) - components.VisibleLen(right)
						if gap < 1 {
							gap = 1
						}
						repoRows = append(repoRows, left+strings.Repeat(" ", gap)+right)
					}
				}
			}
		}
	}

	dockerItems := m.getSubsectionItems(sections, "Workspace", "Docker")
	var dockerBox string
	for _, item := range dockerItems {
		if item.Kind == status.KindDocker && item.Loaded && item.Available && item.DockerStatus != nil {
			ds := item.DockerStatus

			var containerRows []string
			if len(ds.Containers) > 0 {
				for _, c := range ds.Containers {
					left := " " + theme.ServiceRunning.Render("●") + " " + c.Name
					right := theme.Muted.Render(c.Value)
					gap := innerW - components.VisibleLen(left) - components.VisibleLen(right)
					if gap < 1 {
						gap = 1
					}
					containerRows = append(containerRows, left+strings.Repeat(" ", gap)+right)
				}
			} else {
				containerRows = append(containerRows, " "+theme.Muted.Render("no containers running"))
			}

			var imageRows []string
			for _, img := range ds.Images {
				name := img.Repository
				if len(name) > 20 {
					name = name[:17] + "..."
				}
				left := " " + name
				right := theme.Muted.Render(img.Tag) + "  " + theme.Muted.Render(img.Size)
				gap := innerW - components.VisibleLen(left) - components.VisibleLen(right)
				if gap < 1 {
					gap = 1
				}
				imageRows = append(imageRows, left+strings.Repeat(" ", gap)+right)
			}
			if ds.DanglingCount > 0 {
				left := " " + theme.Warning.Render(fmt.Sprintf("+ %d dangling", ds.DanglingCount))
				right := theme.Muted.Render(ds.DanglingSize)
				gap := innerW - components.VisibleLen(left) - components.VisibleLen(right)
				if gap < 1 {
					gap = 1
				}
				imageRows = append(imageRows, left+strings.Repeat(" ", gap)+right)
			}
			if len(imageRows) == 0 {
				imageRows = append(imageRows, " "+theme.Muted.Render("no images"))
			}

			var summaryRow string
			parts := []string{
				fmt.Sprintf("images %d", len(ds.Images)),
			}
			if ds.VolumesCount > 0 {
				parts = append(parts, fmt.Sprintf("volumes %d", ds.VolumesCount))
			}
			if ds.BuildCache != "" && ds.BuildCache != "0B" {
				parts = append(parts, "cache "+ds.BuildCache)
			}
			summaryRow = " " + theme.Muted.Render(strings.Join(parts, "    "))

			dockerBox = components.SubsectionBoxWithSections("Docker",
				[][]string{containerRows, imageRows, {summaryRow}}, colWidth)
		} else if item.Kind == status.KindDocker && !item.Loaded {
			dockerBox = components.SubsectionBox("Docker", []string{" " + m.loadingIndicator()}, colWidth)
		}
	}

	depItems := m.getSubsectionItems(sections, "Workspace", "Dependencies")
	var depRows []string
	for _, item := range depItems {
		if item.Kind == status.KindDependency && item.Loaded && item.Available {
			maxNameW := 0
			for _, group := range item.DepGroups {
				for _, repo := range group.Repos {
					if len(repo.Name) > maxNameW {
						maxNameW = len(repo.Name)
					}
				}
			}

			for gi, group := range item.DepGroups {
				if gi > 0 {
					depRows = append(depRows, "")
				}
				depRows = append(depRows, " "+theme.SubSection.Render(group.Prefix+"/"))
				for ri, repo := range group.Repos {
					connector := "├─"
					if ri == len(group.Repos)-1 {
						connector = "└─"
					}

					var statusStr string
					switch {
					case repo.Outdated == 0:
						statusStr = theme.Success.Render("✓ ok")
					case repo.Outdated > 0:
						statusStr = theme.Warning.Render(fmt.Sprintf("%d outdated", repo.Outdated))
					default:
						statusStr = theme.Muted.Render("?")
					}

					left := "   " + theme.Muted.Render(connector) + " " + fmt.Sprintf("%-*s", maxNameW, repo.Name)
					right := statusStr + "  " + theme.Muted.Render(repo.Manager)
					gap := innerW - components.VisibleLen(left) - components.VisibleLen(right)
					if gap < 1 {
						gap = 1
					}
					depRows = append(depRows, left+strings.Repeat(" ", gap)+right)
				}
			}
		} else if item.Kind == status.KindDependency && !item.Loaded {
			depRows = append(depRows, " "+m.loadingIndicator())
		}
	}

	var boxes []string
	if len(repoRows) > 0 {
		boxes = append(boxes, components.SubsectionBox("Repositories", repoRows, colWidth))
	}
	if len(depRows) > 0 {
		boxes = append(boxes, components.SubsectionBox("Dependencies", depRows, colWidth))
	}
	if dockerBox != "" {
		boxes = append(boxes, dockerBox)
	}

	if len(boxes) == 0 {
		return ""
	}

	if numCols > len(boxes) {
		numCols = len(boxes)
	}

	columns := make([][]string, numCols)
	colHeights := make([]int, numCols)
	for _, box := range boxes {
		shortest := 0
		for i := 1; i < numCols; i++ {
			if colHeights[i] < colHeights[shortest] {
				shortest = i
			}
		}
		columns[shortest] = append(columns[shortest], box)
		colHeights[shortest] += strings.Count(box, "\n") + 1
	}

	var renderedCols []string
	for _, col := range columns {
		renderedCols = append(renderedCols, lipgloss.JoinVertical(lipgloss.Left, col...))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedCols...) + "\n"
}

// invariant: the script-to-box mapping has no fallback bucket on purpose: a script
// added to the registry and not listed here shows up missing in review, never
// silently in a Misc box.
func (m Model) renderConfig(sections sectionMap, w int) string {
	numCols := 3
	if w < minColWidthResponsive*2 {
		numCols = 1
	} else if w < minColWidthResponsive*3 {
		numCols = 2
	}
	colWidth := w / numCols
	innerW := colWidth - 4

	buildCheckRows := func(items []status.Item) []string {
		return m.buildCheckRows(items, innerW)
	}

	configSub := func(name string) []status.Item {
		return m.getSubsectionItems(sections, "Config", name)
	}
	var systemSecurity []status.Item
	for _, item := range m.getSubsectionItems(sections, "Environment", "Health") {
		if item.Kind == status.KindSecurity {
			systemSecurity = append(systemSecurity, item)
		}
	}
	pick := func(items []status.Item, names ...string) []status.Item {
		want := make(map[string]bool, len(names))
		for _, n := range names {
			want[n] = true
		}
		var out []status.Item
		for _, item := range items {
			if want[item.Name] {
				out = append(out, item)
			}
		}
		return out
	}

	securityScripts := configSub(string(config.ScriptCategorySecurity))

	type configBox struct {
		title string
		items []status.Item
	}
	type configGroup struct {
		label string
		boxes []configBox
	}

	groups := []configGroup{
		{"Identity", []configBox{
			{"Identity", configSub("Identity")},
			{"Auth", pick(securityScripts, "gpg", "ssh", "gh")},
		}},
		{"Privacy & Security", []configBox{
			{"Hardening", systemSecurity},
			{"Privacy", pick(securityScripts, "spotlight-exclude", "dns")},
		}},
		{"Server", []configBox{
			{"Server", configSub(string(config.ScriptCategoryServer))},
			{"Network", configSub("Network")},
		}},
		{"Workspace", []configBox{
			{"Terminal", configSub(string(config.ScriptCategoryTerminal))},
			{"Editor", configSub(string(config.ScriptCategoryEditor))},
			{"System", configSub(string(config.ScriptCategorySystem))},
		}},
		{"Daemons", []configBox{
			{"LaunchAgents", configSub("Daemons")},
		}},
	}

	var allItems []status.Item
	for _, group := range groups {
		for _, bx := range group.boxes {
			allItems = append(allItems, bx.items...)
		}
	}

	var b strings.Builder
	if rollup := configHealthRollup(allItems, w); rollup != "" {
		b.WriteString(rollup + "\n\n")
	}
	for _, group := range groups {
		var rendered []string
		for _, bx := range group.boxes {
			rows := buildCheckRows(bx.items)
			if len(rows) == 0 {
				continue
			}
			rendered = append(rendered, components.SubsectionBox(bx.title, rows, colWidth))
		}
		if len(rendered) == 0 {
			continue
		}
		b.WriteString(groupSeparator(group.label, w) + "\n")
		b.WriteString(masonry(rendered, numCols) + "\n")
	}
	return b.String()
}

// invariant: empty while items still load, so no 0% bar paints before the probes finish.
func configHealthRollup(items []status.Item, w int) string {
	var loaded, healthy int
	for _, item := range items {
		if !item.Loaded {
			continue
		}
		loaded++
		ok := item.Installed
		if item.Kind == status.KindSecurity || item.Kind == status.KindIdentity {
			ok = item.Installed == item.GoodWhen
		}
		if ok {
			healthy++
		}
	}
	if loaded == 0 {
		return ""
	}
	drift := loaded - healthy

	const barWidth = 24
	ratio := float64(healthy) / float64(loaded)
	filled := int(ratio*float64(barWidth) + 0.5)
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	pct := fmt.Sprintf("%d%%", int(ratio*100+0.5))
	left := fmt.Sprintf(" %s   %s",
		theme.Success.Render(fmt.Sprintf("%d healthy", healthy)),
		theme.Warning.Render(fmt.Sprintf("%d drift", drift)),
	)
	if drift == 0 {
		left = " " + theme.Success.Render(fmt.Sprintf("%d healthy", healthy))
	}
	right := theme.Muted.Render(bar) + " " + theme.Muted.Render(pct)
	gap := w - components.VisibleLen(left) - components.VisibleLen(right) - 1
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

func (m Model) renderTools(sections sectionMap, w int) string {
	subsections, ok := sections["Tools"]
	if !ok {
		return ""
	}

	type toolGroup struct {
		label      string
		categories []string
	}
	groups := []toolGroup{
		{"Core", []string{"Package Managers", "Runtimes"}},
		{"CLI", []string{"Terminal", "Git", "System", "Deploy", "AI Agents", "AI Tooling"}},
		{"Apps", []string{"Development", "Creative", "Communication", "Productivity", "AI Apps", "Browse", "Security", "Entertainment", "Utilities"}},
	}

	numCols := 3
	if w < minColWidthResponsive*2 {
		numCols = 1
	} else if w < minColWidthResponsive*3 {
		numCols = 2
	}

	colWidth := w / numCols

	buildBoxes := func(categories []string) []string {
		var boxes []string
		for _, sub := range categories {
			items, ok := subsections[sub]
			if !ok || len(items) == 0 {
				continue
			}
			nameW := 0
			for _, item := range items {
				if len(item.Name) > nameW {
					nameW = len(item.Name)
				}
			}
			innerW := colWidth - 4
			var rows []string
			for _, item := range items {
				rows = append(rows, m.renderToolRowCompact(item, nameW, innerW))
			}
			boxes = append(boxes, components.SubsectionBox(sub, rows, colWidth))
		}
		return boxes
	}

	var b strings.Builder
	for _, group := range groups {
		boxes := buildBoxes(group.categories)
		if len(boxes) == 0 {
			continue
		}
		b.WriteString(groupSeparator(group.label, w) + "\n")
		b.WriteString(masonry(boxes, numCols) + "\n")
	}

	return b.String()
}

func (m Model) buildCheckRows(items []status.Item, innerW int) []string {
	var rows []string
	for _, item := range items {
		if !item.Loaded {
			rows = append(rows, " "+m.loadingIndicator()+" "+item.Name)
			continue
		}
		ok := item.Installed
		if item.Kind == status.KindSecurity || item.Kind == status.KindIdentity {
			ok = item.Installed == item.GoodWhen
		}
		left := " " + components.Badge(ok) + " " + item.Name
		if item.Detail == "" {
			rows = append(rows, left)
			continue
		}
		right := theme.Muted.Render(item.Detail)
		gap := innerW - components.VisibleLen(left) - components.VisibleLen(right)
		if gap < 1 {
			gap = 1
		}
		rows = append(rows, left+strings.Repeat(" ", gap)+right)
	}
	return rows
}

func groupSeparator(label string, width int) string {
	text := " " + label + " "
	remaining := width - len(text)
	if remaining < 2 {
		remaining = 2
	}
	dots := remaining / 2
	return theme.Muted.Render("·") + theme.SubSection.Render(text) +
		theme.Muted.Render(strings.Repeat(" ·", dots))
}

func masonry(boxes []string, numCols int) string {
	if len(boxes) == 0 {
		return ""
	}
	nc := numCols
	if nc > len(boxes) {
		nc = len(boxes)
	}
	columns := make([][]string, nc)
	colHeights := make([]int, nc)
	for _, box := range boxes {
		shortest := 0
		for i := 1; i < nc; i++ {
			if colHeights[i] < colHeights[shortest] {
				shortest = i
			}
		}
		columns[shortest] = append(columns[shortest], box)
		colHeights[shortest] += strings.Count(box, "\n") + 1
	}
	rendered := make([]string, len(columns))
	for i, col := range columns {
		rendered[i] = lipgloss.JoinVertical(lipgloss.Left, col...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m Model) renderToolRowCompact(item status.Item, nameW, rowW int) string {
	if !item.Loaded {
		return " " + m.loadingIndicator() + " " + fmt.Sprintf("%-*s", nameW, item.Name)
	}

	badge := components.Badge(item.Installed)
	name := fmt.Sprintf("%-*s", nameW, item.Name)

	version := item.Version
	if strings.HasPrefix(version, item.Name+" ") {
		version = strings.TrimPrefix(version, item.Name+" ")
	}
	if strings.HasPrefix(version, item.Name+"/") {
		version = strings.TrimPrefix(version, item.Name+"/")
	}
	if len(version) > 16 {
		version = version[:13] + "..."
	}

	left := " " + badge + " " + name

	right := ""
	switch item.Status {
	case "running":
		right += theme.ServiceRunning.Render(theme.IconServiceOn) + " "
	case "stopped":
		right += theme.Muted.Render("○") + " "
	}

	if version != "" {
		right += theme.Muted.Render(version)
	}

	method := item.Method
	if method != "" {
		minNeeded := components.VisibleLen(left) + 1 + len(method) + 1 + components.VisibleLen(right)
		if minNeeded <= rowW {
			right += " " + theme.Muted.Render(method)
		}
	}

	if right == "" {
		return left
	}

	leftLen := components.VisibleLen(left)
	rightLen := components.VisibleLen(right)
	gap := rowW - leftLen - rightLen
	if gap < 1 {
		gap = 1
	}

	return left + strings.Repeat(" ", gap) + right
}

var totalMemoryMB float64

func getTotalMemoryMB() float64 {
	if totalMemoryMB > 0 {
		return totalMemoryMB
	}
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return 16 * 1024
	}
	var bytes int64
	_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &bytes)
	if bytes > 0 {
		totalMemoryMB = float64(bytes) / (1024 * 1024)
	} else {
		totalMemoryMB = 16 * 1024
	}
	return totalMemoryMB
}

func (m Model) renderProcessBars(items []status.Item, width int) []string {
	var rows []string
	for _, item := range items {
		if item.Kind != status.KindProcess {
			continue
		}
		if !item.Loaded {
			rows = append(rows, " "+m.loadingIndicator())
			continue
		}

		maxItems := 5
		barWidth := 15
		valueWidth := 7
		nameWidth := width - barWidth - valueWidth - 3
		if nameWidth < 10 {
			nameWidth = 10
		}

		var maxVal float64
		isPercent := len(item.Processes) > 0 && strings.HasSuffix(item.Processes[0].Value, "%")
		if isPercent {
			maxVal = 100.0
		} else {
			for _, p := range item.Processes {
				v := parseProcessValue(p.Value)
				if v > maxVal {
					maxVal = v
				}
			}
		}
		if maxVal == 0 {
			maxVal = 1
		}

		// value: memory colors by its share of the machine's RAM, not relative to the top
		// process; CPU by its own value.
		totalMem := getTotalMemoryMB()

		for i, p := range item.Processes {
			if i >= maxItems {
				break
			}
			name := p.Name
			if nameWidth > 3 && len(name) > nameWidth {
				name = name[:nameWidth-3] + "..."
			}

			v := parseProcessValue(p.Value)
			sizeRatio := v / maxVal
			if sizeRatio > 1.0 {
				sizeRatio = 1.0
			}

			colorRatio := sizeRatio
			if !isPercent {
				colorRatio = v / totalMem
				if colorRatio > 1.0 {
					colorRatio = 1.0
				}
			}

			bar := renderBarWithColor(sizeRatio, colorRatio, barWidth)
			rows = append(rows, " "+fmt.Sprintf("%-*s", nameWidth, name)+" "+bar+" "+theme.Muted.Render(fmt.Sprintf("%*s", valueWidth, p.Value)))
		}
	}
	return rows
}

func renderBar(ratio float64, width int) string {
	return renderBarWithColor(ratio, ratio, width)
}

func renderBarWithColor(sizeRatio, colorRatio float64, width int) string {
	if sizeRatio < 0 {
		sizeRatio = 0
	}
	if sizeRatio > 1 {
		sizeRatio = 1
	}
	filled := int(sizeRatio * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	var filledStyle lipgloss.Style
	switch {
	case colorRatio < 0.5:
		filledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSuccess))
	case colorRatio < 0.8:
		filledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorWarning))
	default:
		filledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorDanger))
	}

	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorBorder))
	return filledStyle.Render(strings.Repeat("█", filled)) + emptyStyle.Render(strings.Repeat("░", empty))
}

func renderColoredSparkline(values []float64, width int) string {
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	data := values
	if len(data) > width {
		data = data[len(data)-width:]
	}

	maxVal := 100.0
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSuccess))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorWarning))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorDanger))

	var result strings.Builder

	pad := width - len(data)
	result.WriteString(strings.Repeat(" ", pad))

	for _, v := range data {
		ratio := v / maxVal
		idx := int(ratio * float64(len(blocks)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(blocks) {
			idx = len(blocks) - 1
		}

		ch := string(blocks[idx])
		switch {
		case ratio < 0.4:
			result.WriteString(green.Render(ch))
		case ratio < 0.7:
			result.WriteString(yellow.Render(ch))
		default:
			result.WriteString(red.Render(ch))
		}
	}

	return result.String()
}

func renderAutoSparkline(values []float64, width int) string {
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	data := values
	if len(data) > width {
		data = data[len(data)-width:]
	}

	maxVal := 1.0
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorSpecial))

	var result strings.Builder
	pad := width - len(data)
	result.WriteString(strings.Repeat(" ", pad))

	for _, v := range data {
		ratio := v / maxVal
		idx := int(ratio * float64(len(blocks)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(blocks) {
			idx = len(blocks) - 1
		}
		result.WriteString(green.Render(string(blocks[idx])))
	}

	return result.String()
}

func formatBytesPerSec(b float64) string {
	switch {
	case b >= 1e9:
		return fmt.Sprintf("%5.1fG", b/1e9)
	case b >= 1e6:
		return fmt.Sprintf("%5.1fM", b/1e6)
	case b >= 1e3:
		return fmt.Sprintf("%5.1fK", b/1e3)
	default:
		return fmt.Sprintf("%5.0fB", b)
	}
}

type sectionMap = map[string]map[string][]status.Item

func (m Model) groupBySection() sectionMap {
	sections := make(sectionMap)

	for _, baseItem := range m.itemOrder {
		item := m.items[baseItem.ID]
		if item.Kind == status.KindHeader || item.Kind == status.KindSystemInfo {
			continue
		}
		if sections[item.Section] == nil {
			sections[item.Section] = make(map[string][]status.Item)
		}
		sections[item.Section][item.SubSection] = append(sections[item.Section][item.SubSection], item)
	}

	if toolSubs, ok := sections["Tools"]; ok {
		for sub, items := range toolSubs {
			sort.Slice(items, func(i, j int) bool {
				return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
			})
			toolSubs[sub] = items
		}
	}

	return sections
}

func (m Model) getSubsectionItems(sections sectionMap, section, subsection string) []status.Item {
	if subs, ok := sections[section]; ok {
		if items, ok := subs[subsection]; ok {
			return items
		}
	}
	return nil
}

func parseDisplaySize(s string) float64 {
	s = strings.TrimSpace(s)
	upper := strings.ToUpper(s)
	// hazard: longest suffixes first, or "B" matches "GB".
	suffixes := []struct {
		suffix string
		mult   float64
	}{
		{"TB", 1e12},
		{"GB", 1e9},
		{"MB", 1e6},
		{"KB", 1e3},
		{"B", 1},
	}
	for _, entry := range suffixes {
		if strings.HasSuffix(upper, entry.suffix) {
			numStr := strings.TrimSpace(s[:len(s)-len(entry.suffix)])
			var val float64
			_, _ = fmt.Sscanf(numStr, "%f", &val)
			return val * entry.mult
		}
	}
	return 0
}

type namedColumn struct {
	title string
	rows  []string
}

const minColWidthResponsive = 30

func renderColumnsImpl(cols []namedColumn, totalWidth int, outerPadding bool) string {
	if len(cols) == 0 {
		return ""
	}

	numCols := totalWidth / minColWidthResponsive
	if numCols < 1 {
		numCols = 1
	}
	if numCols > 3 {
		numCols = 3
	}
	if numCols > len(cols) {
		numCols = len(cols)
	}

	var b strings.Builder
	divider := theme.SectionBorder.Render("│")

	for start := 0; start < len(cols); start += numCols {
		end := start + numCols
		if end > len(cols) {
			end = len(cols)
		}
		group := cols[start:end]
		colW := totalWidth / len(group)

		if outerPadding {
			b.WriteString("\n")
		}

		for i, col := range group {
			header := " " + theme.SubSection.Render(strings.ToUpper(col.title))
			if i < len(group)-1 {
				b.WriteString(padTo(header, colW) + divider)
			} else {
				b.WriteString(header)
			}
		}
		b.WriteString("\n")

		for i := range group {
			if i < len(group)-1 {
				b.WriteString(padTo("", colW) + divider)
			}
		}
		b.WriteString("\n")

		maxRows := 0
		for _, col := range group {
			if len(col.rows) > maxRows {
				maxRows = len(col.rows)
			}
		}

		for r := range maxRows {
			for i, col := range group {
				row := getOr(col.rows, r, "")
				if i < len(group)-1 {
					b.WriteString(padTo(row, colW) + divider)
				} else {
					b.WriteString(row)
				}
			}
			b.WriteString("\n")
		}
	}

	if outerPadding {
		b.WriteString("\n")
	}
	return b.String()
}

func renderColumnsRaw(cols []namedColumn, totalWidth int) string {
	return renderColumnsImpl(cols, totalWidth, false)
}

func parseProcessValue(s string) float64 {
	var v float64
	switch {
	case strings.HasSuffix(s, "%"):
		_, _ = fmt.Sscanf(strings.TrimSuffix(s, "%"), "%f", &v)
	case strings.HasSuffix(s, "G"):
		_, _ = fmt.Sscanf(strings.TrimSuffix(s, "G"), "%f", &v)
		v *= 1024
	case strings.HasSuffix(s, "M"):
		_, _ = fmt.Sscanf(strings.TrimSuffix(s, "M"), "%f", &v)
	case strings.HasSuffix(s, "K"):
		_, _ = fmt.Sscanf(strings.TrimSuffix(s, "K"), "%f", &v)
		v /= 1024
	}
	return v
}

func padTo(s string, width int) string {
	pad := width - components.VisibleLen(s)
	if pad < 0 {
		pad = 0
	}
	return s + strings.Repeat(" ", pad)
}

func padName(name string, width int) string {
	if len(name) > width {
		name = name[:width-1] + "…"
	}
	return fmt.Sprintf("%-*s", width, name)
}

func getOr(slice []string, i int, fallback string) string {
	if i < len(slice) {
		return slice[i]
	}
	return fallback
}
