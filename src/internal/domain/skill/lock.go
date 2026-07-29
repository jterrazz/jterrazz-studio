package skill

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LockEntry is one skill's record in ~/.agents/.skill-lock.json — enough to
// know where it came from and to compare it against its upstream folder for a
// currency check.
type LockEntry struct {
	Source     string // "<owner>/<repo>", e.g. "jterrazz/jterrazz-studio"
	SourceType string // "github" is the only type CheckUpToDate understands
	SkillPath  string // path within the source repo, e.g. "skills/<name>/SKILL.md"
	// FolderHash is the git tree SHA of the skill's folder in the source repo
	// as of the last install/update. This is the `skills` CLI's own record of
	// what it installed, and comparing it to the upstream tree SHA is exactly
	// how the CLI decides whether an update exists — see CheckUpToDate.
	FolderHash string
}

// lockFileSchema mirrors the on-disk shape of ~/.agents/.skill-lock.json
// (schema version 3). Only the fields ReadLock needs are declared; unknown
// fields (e.g. "dismissed", per-skill timestamps) are ignored.
type lockFileSchema struct {
	Version int                     `json:"version"`
	Skills  map[string]lockRawEntry `json:"skills"`
}

type lockRawEntry struct {
	Source     string `json:"source"`
	SourceType string `json:"sourceType"`
	SkillPath  string `json:"skillPath"`
	FolderHash string `json:"skillFolderHash"`
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
			FolderHash: e.FolderHash,
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
var httpClient = &http.Client{Timeout: 5 * time.Second}

// apiBaseURL is the api.github.com origin, overridable in tests to point at a
// local httptest server.
var apiBaseURL = "https://api.github.com"

// CheckUpToDate reports whether an installed skill's folder still matches the
// upstream folder it was installed from.
//
// The comparison is on git tree SHAs, not file contents: the lock file's
// skillFolderHash is the tree SHA of the skill's folder at install/update
// time, and the GitHub contents API reports the current tree SHA of that same
// folder. This is deliberately the *same* oracle the `skills` CLI uses to
// decide what to update.
//
// That equivalence is the whole point. An independent heuristic here (an
// earlier version diffed the local SKILL.md against raw.githubusercontent.com)
// can disagree with the CLI, and when it says "outdated" for a skill the CLI
// considers current, `skills update` does nothing, the row never clears, and
// the Skills tab offers an update that can never be applied. Comparing what
// the CLI compares keeps the badge and the action in agreement — and catches
// changes to supporting files (scripts, references), which a SKILL.md-only
// diff missed.
//
// Only sourceType "github" is supported. Any network error, non-200 response,
// missing upstream folder, or missing hash on either side reports
// StatusUnknown — staleness is never invented from an inconclusive check.
func CheckUpToDate(name string, entry LockEntry) UpdateStatus {
	return checkUpToDateIn(skillsDir(), name, entry)
}

// checkUpToDateIn is the testable variant of CheckUpToDate — takes the
// skills directory explicitly.
func checkUpToDateIn(dir, name string, entry LockEntry) UpdateStatus {
	if entry.SourceType != "github" || entry.Source == "" || entry.SkillPath == "" || entry.FolderHash == "" {
		return StatusUnknown
	}

	// The skill must actually be on disk. ListInstalled already guarantees
	// this for the TUI's callers, but a lock entry can outlive its folder.
	if _, err := os.Stat(filepath.Join(dir, name, "SKILL.md")); err != nil {
		return StatusUnknown
	}

	// skillPath points at the SKILL.md; the hash is over its containing folder.
	folder := path.Dir(entry.SkillPath)
	if folder == "." || folder == "/" {
		return StatusUnknown
	}

	upstream, err := upstreamFolderSHA(entry.Source, folder)
	if err != nil {
		return StatusUnknown
	}

	if upstream == entry.FolderHash {
		return StatusUpToDate
	}
	return StatusOutdated
}

// contentsEntry is one element of the GitHub contents API's directory
// listing. Only the fields needed to find a subfolder's tree SHA are declared.
type contentsEntry struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
	Type string `json:"type"`
}

// dirListing is one cached directory fetch — resolved at most once per
// (repo, parent directory) per process, however many goroutines race for it.
type dirListing struct {
	once sync.Once
	shas map[string]string // child name -> tree/blob SHA
	err  error
}

var (
	dirCacheMu sync.Mutex
	dirCache   = map[string]*dirListing{}
)

// upstreamFolderSHA returns the current git tree SHA of <source>/<folder>.
//
// It lists the folder's *parent* rather than the folder itself, so all skills
// from one repo collapse into a single API call — a `j config` launch with
// twenty installed skills spends one request per distinct source repo, not
// per skill. Results are cached for the process lifetime: upstream can't
// change mid-session, and the local side of the comparison (the lock file) is
// re-read on every check, so post-update rechecks still resolve correctly.
func upstreamFolderSHA(source, folder string) (string, error) {
	parent := path.Dir(folder)
	if parent == "." {
		parent = ""
	}

	key := source + "\x00" + parent
	dirCacheMu.Lock()
	listing, ok := dirCache[key]
	if !ok {
		listing = &dirListing{}
		dirCache[key] = listing
	}
	dirCacheMu.Unlock()

	listing.once.Do(func() {
		listing.shas, listing.err = fetchDirSHAs(source, parent)
	})
	if listing.err != nil {
		return "", listing.err
	}

	sha, ok := listing.shas[path.Base(folder)]
	if !ok {
		// Deleted upstream, or moved. Inconclusive, not stale.
		return "", fmt.Errorf("%s not found in %s/%s", path.Base(folder), source, parent)
	}
	return sha, nil
}

// resetDirCache drops the cached directory listings. Tests call it so one
// test's httptest server can't serve another's lookups.
func resetDirCache() {
	dirCacheMu.Lock()
	defer dirCacheMu.Unlock()
	dirCache = map[string]*dirListing{}
}

// fetchDirSHAs lists <source>/<dir> on the default branch and maps each child
// name to its SHA.
func fetchDirSHAs(source, dir string) (map[string]string, error) {
	url := fmt.Sprintf("%s/repos/%s/contents/%s", apiBaseURL, source, dir)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	// Unauthenticated api.github.com allows 60 requests/hour, which a few
	// `j config` launches can exhaust. A token (env, else the `gh` CLI's)
	// raises that to 5000; without one a rate-limited response is a non-200
	// and every affected skill settles to StatusUnknown rather than a wrong
	// verdict.
	if tok := githubToken(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries []contentsEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, err
	}

	shas := make(map[string]string, len(entries))
	for _, e := range entries {
		shas[e.Name] = e.SHA
	}
	return shas, nil
}

// githubToken is the token source, resolved once per process. A package var
// so tests can stub it out instead of shelling to `gh`.
var githubToken = sync.OnceValue(resolveGitHubToken)

// resolveGitHubToken reads a GitHub token from the environment, falling back
// to the `gh` CLI's stored credential. Returns "" when neither is available —
// the check degrades to unauthenticated rather than failing.
func resolveGitHubToken() string {
	for _, key := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
