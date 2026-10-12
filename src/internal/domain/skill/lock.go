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

// LockEntry is one skill's record in ~/.agents/.skill-lock.json.
type LockEntry struct {
	Source     string
	SourceType string // invariant: github is the only type CheckUpToDate understands
	SkillPath  string
	// FolderHash is the git tree SHA of the skill's folder at its last install or update.
	FolderHash string
}

// outside: mirrors ~/.agents/.skill-lock.json (schema version 3); only the fields
// ReadLock needs are declared.
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

func lockPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".agents", ".skill-lock.json")
}

// ReadLock parses ~/.agents/.skill-lock.json; a missing or corrupt file gives an empty map.
func ReadLock() map[string]LockEntry {
	return readLockFrom(lockPath())
}

func readLockFrom(lockPath string) map[string]LockEntry {
	out := map[string]LockEntry{}
	if lockPath == "" {
		return out
	}

	data, err := os.ReadFile(lockPath)
	if err != nil {
		return out
	}

	var parsed lockFileSchema
	if err := json.Unmarshal(data, &parsed); err != nil {
		return out
	}

	for name, e := range parsed.Skills {
		out[name] = LockEntry(e)
	}
	return out
}

// UpdateStatus is an installed skill's currency against its upstream source.
type UpdateStatus int

const (
	StatusUnknown UpdateStatus = iota
	StatusUpToDate
	StatusOutdated
)

// value: a 5s timeout keeps an unreachable network from holding the currency fan-out.
var httpClient = &http.Client{Timeout: 5 * time.Second}

// value: a var so the tests point it at an httptest server.
var apiBaseURL = "https://api.github.com"

// CheckUpToDate compares an installed skill's folder tree SHA with its upstream folder's.
func CheckUpToDate(name string, entry LockEntry) UpdateStatus {
	return checkUpToDateIn(skillsDir(), name, entry)
}

func checkUpToDateIn(dir, name string, entry LockEntry) UpdateStatus {
	// invariant: the verdict compares git tree SHAs, the oracle the skills CLI updates by,
	// so the badge never offers an update `skills update` would not apply; an inconclusive
	// check (network, non-200, missing folder or hash) is StatusUnknown, never stale.
	if entry.SourceType != "github" || entry.Source == "" || entry.SkillPath == "" || entry.FolderHash == "" {
		return StatusUnknown
	}

	// hazard: a lock entry can outlive its folder.
	if _, err := os.Stat(filepath.Join(dir, name, "SKILL.md")); err != nil {
		return StatusUnknown
	}

	// outside: skillPath names the SKILL.md; the hash is over its containing folder.
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

type contentsEntry struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
	Type string `json:"type"`
}

type dirListing struct {
	once sync.Once
	shas map[string]string
	err  error
}

var (
	dirCacheMu sync.Mutex
	dirCache   = map[string]*dirListing{}
)

// value: the parent folder is listed, so every skill of one repo costs one request; the
// cache lasts the process, since the lock side is re-read on every check.
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
		// invariant: deleted or moved upstream is inconclusive, not stale.
		return "", fmt.Errorf("%s not found in %s/%s", path.Base(folder), source, parent)
	}
	return sha, nil
}

func resetDirCache() {
	dirCacheMu.Lock()
	defer dirCacheMu.Unlock()
	dirCache = map[string]*dirListing{}
}

func fetchDirSHAs(source, dir string) (map[string]string, error) {
	url := fmt.Sprintf("%s/repos/%s/contents/%s", apiBaseURL, source, dir)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	// outside: unauthenticated api.github.com allows 60 requests an hour, a token 5000;
	// a rate-limited response is a non-200, so the skill settles to StatusUnknown.
	if tok := githubToken(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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

// value: a package var so the tests stub it instead of shelling out to gh.
var githubToken = sync.OnceValue(resolveGitHubToken)

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
