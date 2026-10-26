package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jterrazz/studio-workstation/src/internal/domain/tool"
)

// ResourceCheck represents a system resource check (network, disk, cache)
type ResourceCheck struct {
	Name    string
	CheckFn func() ResourceResult
}

// ProcessInfo represents a single process entry
type ProcessInfo struct {
	Name  string
	Value string
	PID   string
}

// ProcessResult holds multiple processes
type ProcessResult struct {
	Processes []ProcessInfo
	Available bool
}

// RepoInfo describes a single git repository.
type RepoInfo struct {
	Name        string
	FullName    string
	Branch      string
	ChangeCount int
	Clean       bool
}

// ProjectGroup is a set of repos sharing a project prefix.
type ProjectGroup struct {
	Prefix string
	Repos  []RepoInfo
}

// DockerImageInfo describes a Docker image.
type DockerImageInfo struct {
	Repository string
	Tag        string
	Size       string
}

// DockerStatus holds the full Docker dashboard data.
type DockerStatus struct {
	Containers    []ProcessInfo
	Images        []DockerImageInfo
	DanglingCount int
	DanglingSize  string
	VolumesCount  int
	VolumesSize   string
	BuildCache    string
}

// DepRepoInfo describes dependency status for a single repo.
type DepRepoInfo struct {
	Name     string
	Manager  string
	Outdated int
}

// DepProjectGroup groups dependency info by project prefix.
type DepProjectGroup struct {
	Prefix string
	Repos  []DepRepoInfo
}

// ResourceResult holds the result of a resource check
type ResourceResult struct {
	Value     string
	Style     string
	Available bool
}

// NetworkChecks is the list of network resource checks
var NetworkChecks = []ResourceCheck{
	{
		Name: "local ip",
		CheckFn: func() ResourceResult {
			out, err := exec.Command("ipconfig", "getifaddr", "en0").Output()
			if err != nil {
				return ResourceResult{Available: false}
			}
			ip := strings.TrimSpace(string(out))
			if ip != "" {
				return ResourceResult{Value: ip, Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "public ip",
		CheckFn: func() ResourceResult {
			cmd := exec.Command("curl", "-s", "--max-time", "2", "-4", "ifconfig.me")
			out, err := cmd.Output()
			if err == nil {
				ip := strings.TrimSpace(string(out))
				if ip != "" {
					return ResourceResult{Value: ip, Style: "muted", Available: true}
				}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "vpn",
		CheckFn: func() ResourceResult {
			var vpnNames []string

			out, err := exec.Command("scutil", "--nc", "list").Output()
			_ = err // outside: scutil fails when no network config exists
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "(Connected)") {
					if idx := strings.LastIndex(line, `"`); idx > 0 {
						start := strings.LastIndex(line[:idx], `"`)
						if start >= 0 && start < idx {
							vpnNames = append(vpnNames, line[start+1:idx])
							continue
						}
					}
					vpnNames = append(vpnNames, "connected")
				}
			}

			if st, err := GetTailscaleFullStatus(); err == nil && st.ExitNode != "" {
				vpnNames = append(vpnNames, "Tailscale ("+st.ExitNode+")")
			}

			if len(vpnNames) > 0 {
				return ResourceResult{Value: strings.Join(vpnNames, ", "), Style: "success", Available: true}
			}
			return ResourceResult{Value: "none", Style: "muted", Available: true}
		},
	},
	{
		Name: "dns",
		CheckFn: func() ResourceResult {
			if IsDNSProfileInstalled() {
				return ResourceResult{Value: "Quad9 (encrypted)", Style: "success", Available: true}
			}
			out, err := exec.Command("scutil", "--dns").Output()
			_ = err // outside: the resolver may be unavailable
			var servers []string
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "nameserver[") {
					idx := strings.Index(line, "] : ")
					if idx == -1 {
						continue
					}
					server := strings.TrimSpace(line[idx+4:])
					if server == "" || server == "127.0.0.1" || server == "::1" {
						continue
					}
					if strings.Contains(server, ":") {
						continue
					}
					found := false
					for _, s := range servers {
						if s == server {
							found = true
							break
						}
					}
					if !found && len(servers) < 2 {
						servers = append(servers, server)
					}
				}
			}
			if len(servers) > 0 {
				return ResourceResult{Value: strings.Join(servers, ", "), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
}

// DiskCheck represents a disk usage check
type DiskCheck struct {
	Name    string
	Path    string
	Style   string
	CheckFn func() ResourceResult
}

// CacheChecks shows disk usage grouped by domain
var CacheChecks = []DiskCheck{
	{
		Name: "docker",
		CheckFn: func() ResourceResult {
			if !CommandExists("docker") {
				return ResourceResult{Available: false}
			}
			out, err := exec.Command("docker", "system", "df", "--format", "{{.Size}}").Output()
			if err != nil {
				return ResourceResult{Available: false}
			}
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			var total int64
			for _, line := range lines {
				total += parseDockerSize(strings.TrimSpace(line))
			}
			if total > 0 {
				return ResourceResult{Value: tool.FormatBytes(total), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "xcode",
		CheckFn: func() ResourceResult {
			paths := []string{
				"~/Library/Developer/Xcode/DerivedData",
				"~/Library/Developer/Xcode/Archives",
				"~/Library/Developer/Xcode/iOS DeviceSupport",
			}
			var total int64
			for _, p := range paths {
				total += GetDirSize(expandHome(p))
			}
			if total > 0 {
				return ResourceResult{Value: tool.FormatBytes(total), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{Name: "homebrew", Path: "~/Library/Caches/Homebrew", Style: "muted"},
	{
		Name: "packages",
		CheckFn: func() ResourceResult {
			paths := []string{
				"~/.npm",
				"~/Library/pnpm",
				"~/.bun/install/cache",
				"~/Library/Caches/Yarn",
				"~/Library/Caches/CocoaPods",
				"~/go/pkg/mod",
				"~/.gradle/caches",
			}
			var total int64
			for _, p := range paths {
				total += GetDirSize(expandHome(p))
			}
			if total > 0 {
				return ResourceResult{Value: tool.FormatBytes(total), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "multipass",
		CheckFn: func() ResourceResult {
			if !CommandExists("multipass") {
				return ResourceResult{Available: false}
			}
			path := expandHome("~/Library/Application Support/multipassd")
			if size := GetDirSize(path); size > 0 {
				return ResourceResult{Value: tool.FormatBytes(size), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{
		Name: "logs",
		CheckFn: func() ResourceResult {
			total := GetDirSize("/var/log") + GetDirSize(expandHome("~/Library/Logs"))
			if total > 0 {
				return ResourceResult{Value: tool.FormatBytes(total), Style: "muted", Available: true}
			}
			return ResourceResult{Available: false}
		},
	},
	{Name: "trash", Path: "~/.Trash", Style: "muted"},
}

func parseDockerSize(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "0B" {
		return 0
	}
	multipliers := map[string]float64{
		"B": 1, "kB": 1e3, "MB": 1e6, "GB": 1e9, "TB": 1e12,
	}
	for suffix, mult := range multipliers {
		if strings.HasSuffix(s, suffix) {
			numStr := strings.TrimSuffix(s, suffix)
			var val float64
			_, _ = fmt.Sscanf(numStr, "%f", &val)
			return int64(val * mult)
		}
	}
	return 0
}

// Check reports the disk path this check names, or what its own CheckFn answers.
func (d DiskCheck) Check() ResourceResult {
	if d.CheckFn != nil {
		return d.CheckFn()
	}

	path := expandHome(d.Path)
	if size := GetDirSize(path); size > 0 {
		return ResourceResult{Value: tool.FormatBytes(size), Style: d.Style, Available: true}
	}
	return ResourceResult{Available: false}
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}
	return path
}

// ProcessCheck represents a process resource check
type ProcessCheck struct {
	Name    string
	CheckFn func() []ProcessInfo
}

// ProcessChecks defines the process monitoring checks
var ProcessChecks = []ProcessCheck{
	{
		Name: "CPU",
		CheckFn: func() []ProcessInfo {
			out, err := exec.Command("ps", "-arcwwwxo", "pid,%cpu,comm").Output()
			if err != nil {
				return nil
			}
			return parseCPUOutput(out)
		},
	},
	{
		Name: "Memory",
		CheckFn: func() []ProcessInfo {
			out, err := exec.Command("ps", "-amcwwwxo", "pid,rss,comm").Output()
			if err != nil {
				return nil
			}
			return parseMemoryOutput(out)
		},
	},
	{
		Name: "Ports",
		CheckFn: func() []ProcessInfo {
			out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n", "-Fcn").Output()
			if err != nil {
				return nil
			}
			return parseListeningPortsFcn(out)
		},
	},
	{
		Name:    "Uptime",
		CheckFn: getUptimeInfo,
	},
}

var systemPortProcesses = map[string]bool{
	"rapportd":      true,
	"ControlCenter": true,
	"IPNExtension":  true, // value: Tailscale, already shown under Network
	"mDNSResponder": true,
	"launchd":       true,
	"systemd":       true,
}

func parseListeningPortsFcn(out []byte) []ProcessInfo {
	lines := strings.Split(string(out), "\n")

	type portEntry struct {
		port    string
		portNum int
		cmd     string
	}

	var currentCmd string
	seen := make(map[string]bool)
	var entries []portEntry

	for _, line := range lines {
		if line == "" {
			continue
		}
		switch line[0] {
		case 'c': // outside: lsof's command-name field
			currentCmd = line[1:]
		case 'n': // outside: lsof's network-address field, e.g. "*:8080"
			addr := line[1:]
			if currentCmd == "" {
				continue
			}
			if systemPortProcesses[currentCmd] {
				continue
			}
			idx := strings.LastIndex(addr, ":")
			if idx < 0 {
				continue
			}
			port := addr[idx+1:]
			key := port + "/" + currentCmd
			if seen[key] {
				continue
			}
			seen[key] = true

			num, _ := strconv.Atoi(port)
			entries = append(entries, portEntry{port: port, portNum: num, cmd: currentCmd})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].portNum < entries[j].portNum
	})

	result := make([]ProcessInfo, 0, len(entries))
	for _, e := range entries {
		result = append(result, ProcessInfo{
			Name:  e.cmd,
			Value: ":" + e.port,
		})
	}
	return result
}

func parseCPUOutput(out []byte) []ProcessInfo {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var processes []ProcessInfo

	for i := 1; i < len(lines) && len(processes) < 5; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) < 3 {
			continue
		}
		pid := fields[0]
		cpuPercent := fields[1]
		name := strings.Join(fields[2:], " ")

		processes = append(processes, ProcessInfo{
			Name:  name,
			Value: cpuPercent + "%",
			PID:   pid,
		})
	}

	return processes
}

func parseMemoryOutput(out []byte) []ProcessInfo {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var processes []ProcessInfo

	for i := 1; i < len(lines) && len(processes) < 5; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) < 3 {
			continue
		}
		pid := fields[0]
		rssKB := fields[1]
		name := strings.Join(fields[2:], " ")

		var formatted string
		if kb, err := strconv.ParseInt(rssKB, 10, 64); err == nil {
			mb := kb / 1024
			if mb >= 1024 {
				formatted = fmt.Sprintf("%.1fG", float64(mb)/1024)
			} else {
				formatted = fmt.Sprintf("%dM", mb)
			}
		} else {
			formatted = rssKB + "K"
		}

		processes = append(processes, ProcessInfo{
			Name:  name,
			Value: formatted,
			PID:   pid,
		})
	}

	return processes
}

func parseDockerContainers(out []byte) []ProcessInfo {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []ProcessInfo
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 2 {
			continue
		}
		name := parts[0]
		statusStr := parts[1]

		short := statusStr
		if strings.HasPrefix(statusStr, "Up ") {
			short = shortenUptime(strings.TrimPrefix(statusStr, "Up "))
		}

		portInfo := ""
		if len(parts) == 3 && parts[2] != "" {
			portInfo = extractHostPorts(parts[2])
		}

		value := short
		if portInfo != "" {
			value = portInfo + " " + short
		}

		result = append(result, ProcessInfo{Name: name, Value: value})
	}
	return result
}

func extractHostPorts(ports string) string {
	var hostPorts []string
	seen := make(map[string]bool)
	for _, mapping := range strings.Split(ports, ", ") {
		if idx := strings.Index(mapping, "->"); idx > 0 {
			hostPart := mapping[:idx]
			if colonIdx := strings.LastIndex(hostPart, ":"); colonIdx >= 0 {
				port := hostPart[colonIdx+1:]
				if !seen[port] {
					seen[port] = true
					hostPorts = append(hostPorts, ":"+port)
				}
			}
		}
	}
	if len(hostPorts) == 0 {
		return ""
	}
	return strings.Join(hostPorts, ",")
}

func shortenUptime(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, " ("); idx > 0 {
		s = s[:idx]
	}
	replacer := strings.NewReplacer(
		" seconds", "s", " second", "s",
		" minutes", "m", " minute", "m",
		" hours", "h", " hour", "h",
		" days", "d", " day", "d",
		" weeks", "w", " week", "w",
		" months", "mo", " month", "mo",
		"About ", "~",
	)
	return replacer.Replace(s)
}

// value: depth 2 covers ~/Developer/<org>/<repo> and keeps the scan fast on
// machines with deep node_modules or vendor trees.
const scanRepoMaxDepth = 2

func findRepoPaths(root string, maxDepth int) []string {
	var paths []string
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		if depth > maxDepth {
			return
		}
		if _, err := os.Stat(dir + "/.git"); err == nil {
			paths = append(paths, dir)
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			walk(dir+"/"+e.Name(), depth+1)
		}
	}
	walk(root, 0)
	return paths
}

func scanRepo(repoPath, relPath string) RepoInfo {
	branchCmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, _ := branchCmd.Output()
	branch := strings.TrimSpace(string(branchOut))
	if branch == "" {
		branch = "?"
	}

	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	out, _ := cmd.Output()
	status := strings.TrimSpace(string(out))
	changeCount := 0
	if status != "" {
		changeCount = len(strings.Split(status, "\n"))
	}

	return RepoInfo{
		FullName:    relPath,
		Branch:      branch,
		ChangeCount: changeCount,
		Clean:       changeCount == 0,
	}
}

// ScanAllRepos scans the git repos under ~/Developer, in parallel, grouped by project.
func ScanAllRepos() []ProjectGroup {
	devDir := os.Getenv("HOME") + "/Developer"
	paths := findRepoPaths(devDir, scanRepoMaxDepth)
	if len(paths) == 0 {
		return nil
	}

	repos := make([]RepoInfo, len(paths))
	var wg sync.WaitGroup
	wg.Add(len(paths))
	for i, repoPath := range paths {
		go func(i int, repoPath string) {
			defer wg.Done()
			rel := strings.TrimPrefix(repoPath, devDir+"/")
			repos[i] = scanRepo(repoPath, rel)
		}(i, repoPath)
	}
	wg.Wait()

	groupMap := make(map[string][]RepoInfo)
	var groupOrder []string
	for _, repo := range repos {
		var prefix, name string
		if idx := strings.LastIndex(repo.FullName, "/"); idx > 0 {
			prefix = repo.FullName[:idx]
			name = repo.FullName[idx+1:]
		} else if idx := strings.Index(repo.FullName, "-"); idx > 0 {
			prefix = repo.FullName[:idx]
			name = repo.FullName[idx+1:]
		} else {
			prefix = repo.FullName
			name = repo.FullName
		}
		repo.Name = name
		if _, exists := groupMap[prefix]; !exists {
			groupOrder = append(groupOrder, prefix)
		}
		groupMap[prefix] = append(groupMap[prefix], repo)
	}

	sort.Strings(groupOrder)

	var result []ProjectGroup
	for _, prefix := range groupOrder {
		group := ProjectGroup{Prefix: prefix, Repos: groupMap[prefix]}
		sort.Slice(group.Repos, func(i, j int) bool {
			return group.Repos[i].Name < group.Repos[j].Name
		})
		result = append(result, group)
	}
	return result
}

// ScanDependencies counts the outdated dependencies of every repo ScanAllRepos finds.
func ScanDependencies() []DepProjectGroup {
	devDir := os.Getenv("HOME") + "/Developer"
	paths := findRepoPaths(devDir, scanRepoMaxDepth)
	if len(paths) == 0 {
		return nil
	}

	type depResult struct {
		fullName string
		manager  string
		outdated int
	}

	var jobs []depResult
	for _, repoPath := range paths {
		var manager string
		if _, err := os.Stat(repoPath + "/pnpm-lock.yaml"); err == nil {
			manager = "pnpm"
		} else if _, err := os.Stat(repoPath + "/go.sum"); err == nil {
			manager = "go"
		} else if _, err := os.Stat(repoPath + "/bun.lock"); err == nil {
			manager = "bun"
		} else if _, err := os.Stat(repoPath + "/package-lock.json"); err == nil {
			manager = "npm"
		}
		if manager == "" {
			continue
		}
		rel := strings.TrimPrefix(repoPath, devDir+"/")
		jobs = append(jobs, depResult{fullName: rel, manager: manager, outdated: -1})
	}

	type indexedResult struct {
		idx      int
		outdated int
	}
	ch := make(chan indexedResult, len(jobs))
	for i, job := range jobs {
		go func(idx int, j depResult) {
			repoPath := devDir + "/" + j.fullName
			count := countOutdated(repoPath, j.manager)
			ch <- indexedResult{idx: idx, outdated: count}
		}(i, job)
	}
	for range jobs {
		r := <-ch
		jobs[r.idx].outdated = r.outdated
	}

	groupMap := make(map[string][]DepRepoInfo)
	var groupOrder []string
	for _, job := range jobs {
		var prefix, name string
		if idx := strings.LastIndex(job.fullName, "/"); idx > 0 {
			prefix = job.fullName[:idx]
			name = job.fullName[idx+1:]
		} else if idx := strings.Index(job.fullName, "-"); idx > 0 {
			prefix = job.fullName[:idx]
			name = job.fullName[idx+1:]
		} else {
			prefix = job.fullName
			name = job.fullName
		}
		if _, exists := groupMap[prefix]; !exists {
			groupOrder = append(groupOrder, prefix)
		}
		groupMap[prefix] = append(groupMap[prefix], DepRepoInfo{
			Name:     name,
			Manager:  job.manager,
			Outdated: job.outdated,
		})
	}

	sort.Strings(groupOrder)
	var result []DepProjectGroup
	for _, prefix := range groupOrder {
		repos := groupMap[prefix]
		sort.Slice(repos, func(i, j int) bool { return repos[i].Name < repos[j].Name })
		result = append(result, DepProjectGroup{Prefix: prefix, Repos: repos})
	}
	return result
}

func countOutdated(repoPath, manager string) int {
	switch manager {
	case "pnpm":
		out, err := exec.Command("pnpm", "outdated", "--dir", repoPath, "--format", "json").Output()
		if err != nil {
			// outside: pnpm outdated exits 1 when there are outdated deps.
			if len(out) == 0 {
				return 0
			}
		}
		outStr := strings.TrimSpace(string(out))
		if outStr == "" || outStr == "{}" || outStr == "[]" {
			return 0
		}
		// outside: each key of pnpm's JSON object is an outdated package.
		count := 0
		for _, line := range strings.Split(outStr, "\n") {
			if strings.Contains(line, `"current"`) {
				count++
			}
		}
		return count
	case "go":
		out, _ := exec.Command("go", "list", "-m", "-u", "-json", "all").
			Output()
		count := 0
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, `"Update"`) {
				count++
			}
		}
		return count
	case "bun":
		// outside: bun has no usable outdated command yet.
		return 0
	case "npm":
		out, _ := exec.Command("npm", "outdated", "--json", "--prefix", repoPath).Output()
		outStr := strings.TrimSpace(string(out))
		if outStr == "" || outStr == "{}" {
			return 0
		}
		count := 0
		for _, line := range strings.Split(outStr, "\n") {
			if strings.Contains(line, `"current"`) {
				count++
			}
		}
		return count
	}
	return 0
}

// GetDockerStatus returns full Docker dashboard data.
func GetDockerStatus() (DockerStatus, error) {
	if !CommandExists("docker") {
		return DockerStatus{}, errors.New("docker not found")
	}

	var ds DockerStatus

	out, err := exec.Command("docker", "ps", "--format", "{{.Names}}\t{{.Status}}\t{{.Ports}}").Output()
	if err == nil {
		ds.Containers = parseDockerContainers(out)
	}

	imgOut, err := exec.Command("docker", "images", "--format", "{{.Repository}}\t{{.Tag}}\t{{.Size}}").Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(imgOut)), "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 3)
			if len(parts) < 3 {
				continue
			}
			if parts[0] == "<none>" {
				continue
			}
			ds.Images = append(ds.Images, DockerImageInfo{
				Repository: parts[0],
				Tag:        parts[1],
				Size:       parts[2],
			})
		}
	}

	dangOut, err := exec.Command("docker", "images", "-f", "dangling=true", "--format", "{{.Size}}").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(dangOut)), "\n")
		for _, line := range lines {
			if line != "" {
				ds.DanglingCount++
			}
		}
		if ds.DanglingCount > 0 {
			sizeOut, _ := exec.Command("docker", "images", "-f", "dangling=true", "--format", "{{.Size}}").Output()
			var total int64
			for _, line := range strings.Split(strings.TrimSpace(string(sizeOut)), "\n") {
				total += parseDockerSize(strings.TrimSpace(line))
			}
			if total > 0 {
				ds.DanglingSize = tool.FormatBytes(total)
			}
		}
	}

	volOut, err := exec.Command("docker", "volume", "ls", "-q").Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(volOut)), "\n") {
			if line != "" {
				ds.VolumesCount++
			}
		}
	}

	dfOut, err := exec.Command("docker", "system", "df", "--format", "{{.Type}}\t{{.Size}}").Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(dfOut)), "\n") {
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 && strings.Contains(parts[0], "Build") {
				ds.BuildCache = strings.TrimSpace(parts[1])
			}
		}
	}

	return ds, nil
}

func getUptimeInfo() []ProcessInfo {
	var result []ProcessInfo

	uptimeOut, err := exec.Command("uptime").Output()
	_ = err // outside: uptime is not on every platform
	uptimeStr := strings.TrimSpace(string(uptimeOut))
	if uptimeStr != "" {
		if idx := strings.Index(uptimeStr, "up "); idx >= 0 {
			rest := uptimeStr[idx+3:]
			if uIdx := strings.Index(rest, " user"); uIdx > 0 {
				upPart := rest[:uIdx]
				if cIdx := strings.LastIndex(upPart, ","); cIdx > 0 {
					upPart = strings.TrimSpace(upPart[:cIdx])
				}
				result = append(result, ProcessInfo{Name: "uptime", Value: strings.TrimRight(upPart, ",")})
			}
		}
		if idx := strings.Index(uptimeStr, "load averages: "); idx >= 0 {
			loads := strings.TrimSpace(uptimeStr[idx+len("load averages: "):])
			result = append(result, ProcessInfo{Name: "load", Value: loads})
		} else if idx := strings.Index(uptimeStr, "load average: "); idx >= 0 {
			loads := strings.TrimSpace(uptimeStr[idx+len("load average: "):])
			result = append(result, ProcessInfo{Name: "load", Value: loads})
		}
	}

	battOut, err := exec.Command("pmset", "-g", "batt").Output()
	if err == nil {
		for _, line := range strings.Split(string(battOut), "\n") {
			if strings.Contains(line, "%") {
				line = strings.TrimSpace(line)
				if tabIdx := strings.Index(line, "\t"); tabIdx >= 0 {
					info := strings.TrimSpace(line[tabIdx:])
					parts := strings.SplitN(info, ";", 3)
					if len(parts) >= 2 {
						pct := strings.TrimSpace(parts[0])
						state := strings.TrimSpace(parts[1])
						result = append(result, ProcessInfo{Name: "battery", Value: pct + " " + state})
					}
				}
			}
		}
	}

	return result
}
