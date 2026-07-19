package skill

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LockEntry is one skill's record in ~/.agents/.skill-lock.json — enough to
// know where it came from and to fetch its upstream copy for a currency
// check.
type LockEntry struct {
	Source     string // "<owner>/<repo>", e.g. "jterrazz/jterrazz-studio"
	SourceType string // "github" is the only type CheckUpToDate understands
	SkillPath  string // path within the source repo, e.g. "skills/<name>/SKILL.md"
}

// lockFileSchema mirrors the on-disk shape of ~/.agents/.skill-lock.json
// (schema version 3). Only the fields ReadLock needs are declared; unknown
// fields (e.g. "dismissed", per-skill hashes/timestamps) are ignored.
type lockFileSchema struct {
	Version int                     `json:"version"`
	Skills  map[string]lockRawEntry `json:"skills"`
}

type lockRawEntry struct {
	Source     string `json:"source"`
	SourceType string `json:"sourceType"`
	SkillPath  string `json:"skillPath"`
}

// lockPath returns ~/.agents/.skill-lock.json.
func lockPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".agents", ".skill-lock.json")
}

// ReadLock parses ~/.agents/.skill-lock.json. A missing or corrupt file
// returns an empty map rather than an error — the lock file is best-effort
// metadata for the update-currency check, never load-bearing for whether a
// skill counts as installed (that's ListInstalled's job).
func ReadLock() map[string]LockEntry {
	return readLockFrom(lockPath())
}

// readLockFrom is the testable variant of ReadLock — takes the lock file
// path explicitly.
func readLockFrom(path string) map[string]LockEntry {
	out := map[string]LockEntry{}
	if path == "" {
		return out
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return out
	}

	var parsed lockFileSchema
	if err := json.Unmarshal(data, &parsed); err != nil {
		return out
	}

	for name, e := range parsed.Skills {
		out[name] = LockEntry{
			Source:     e.Source,
			SourceType: e.SourceType,
			SkillPath:  e.SkillPath,
		}
	}
	return out
}

// UpdateStatus is the currency state of an installed skill relative to its
// upstream source.
type UpdateStatus int

const (
	StatusUnknown UpdateStatus = iota
	StatusUpToDate
	StatusOutdated
)

// httpClient is shared across currency checks; a short timeout keeps a slow
// or unreachable network from blocking the fan-out longer than necessary.
var httpClient = &http.Client{Timeout: 3 * time.Second}

// rawContentBaseURL is the raw.githubusercontent.com origin, overridable in
// tests to point at a local httptest server.
var rawContentBaseURL = "https://raw.githubusercontent.com"

// candidateRefs are tried in order when fetching a skill's upstream
// SKILL.md. "HEAD" resolves to the repo's default branch on
// raw.githubusercontent.com and works for every source repo checked during
// development; "main"/"master" are kept as a fallback for the rare repo
// where that doesn't hold.
var candidateRefs = []string{"HEAD", "main", "master"}

// CheckUpToDate reports whether an installed skill matches its upstream
// SKILL.md.
//
// This is a deliberate approximation: only SKILL.md is compared, so a skill
// whose supporting files (scripts, references, assets) changed upstream but
// whose SKILL.md didn't will still report StatusUpToDate. Doing better would
// mean diffing the whole skill folder, which isn't worth a 3-second-budget
// network check run for every installed skill on every `j config` launch.
//
// Only sourceType "github" is supported. Any network error, non-200
// response, or unreadable local file reports StatusUnknown — staleness is
// never invented from an inconclusive check.
func CheckUpToDate(name string, entry LockEntry) UpdateStatus {
	return checkUpToDateIn(skillsDir(), name, entry)
}

// checkUpToDateIn is the testable variant of CheckUpToDate — takes the
// skills directory explicitly.
func checkUpToDateIn(dir, name string, entry LockEntry) UpdateStatus {
	if entry.SourceType != "github" || entry.Source == "" || entry.SkillPath == "" {
		return StatusUnknown
	}

	installed, err := os.ReadFile(filepath.Join(dir, name, "SKILL.md"))
	if err != nil {
		return StatusUnknown
	}

	remote, err := fetchUpstreamSkillMd(entry.Source, entry.SkillPath)
	if err != nil {
		return StatusUnknown
	}

	if normalizeTrailingWhitespace(string(installed)) == normalizeTrailingWhitespace(string(remote)) {
		return StatusUpToDate
	}
	return StatusOutdated
}

// fetchUpstreamSkillMd fetches <source>/<ref>/<skillPath> from
// raw.githubusercontent.com, trying candidateRefs in order until one
// succeeds.
func fetchUpstreamSkillMd(source, skillPath string) ([]byte, error) {
	var lastErr error
	for _, ref := range candidateRefs {
		url := fmt.Sprintf("%s/%s/%s/%s", rawContentBaseURL, source, ref, skillPath)
		body, err := fetchURL(url)
		if err == nil {
			return body, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no candidate ref succeeded for %s/%s", source, skillPath)
	}
	return nil, lastErr
}

func fetchURL(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

// normalizeTrailingWhitespace trims trailing whitespace from each line and
// trailing blank lines from the end of the string, so files that differ only
// by trailing spaces/newlines compare equal.
func normalizeTrailingWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t\r")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}
