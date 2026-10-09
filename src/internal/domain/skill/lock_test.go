package skill

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func TestReadLockFromMissingFile(t *testing.T) {
	got := readLockFrom(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestReadLockFromEmptyPath(t *testing.T) {
	got := readLockFrom("")
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestReadLockFromCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".skill-lock.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	got := readLockFrom(path)
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestReadLockFromParsesFixture(t *testing.T) {
	fixture := `{
  "version": 3,
  "skills": {
    "jterrazz-stack": {
      "source": "jterrazz/jterrazz-studio",
      "sourceType": "github",
      "sourceUrl": "https://github.com/jterrazz/jterrazz-studio.git",
      "skillPath": "skills/jterrazz-stack/SKILL.md",
      "skillFolderHash": "abc123",
      "installedAt": "2026-03-27T18:46:44.709Z",
      "updatedAt": "2026-07-19T12:40:35.009Z"
    }
  },
  "dismissed": []
}`
	path := filepath.Join(t.TempDir(), ".skill-lock.json")
	if err := os.WriteFile(path, []byte(fixture), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := readLockFrom(path)
	entry, ok := got["jterrazz-stack"]
	if !ok {
		t.Fatalf("got %v, missing jterrazz-stack", got)
	}
	want := LockEntry{
		Source:     "jterrazz/jterrazz-studio",
		SourceType: "github",
		SkillPath:  "skills/jterrazz-stack/SKILL.md",
		FolderHash: "abc123",
	}
	if entry != want {
		t.Errorf("entry = %+v, want %+v", entry, want)
	}
}

// invariant: CheckUpToDate never reaches the real network in a test: the token source
// is stubbed and the directory cache cleared on both sides.
func withTestServer(t *testing.T, handler http.HandlerFunc, fn func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	defer srv.Close()

	prevBase, prevToken := apiBaseURL, githubToken
	apiBaseURL = srv.URL
	githubToken = func() string { return "" }
	resetDirCache()
	defer func() {
		apiBaseURL, githubToken = prevBase, prevToken
		resetDirCache()
	}()

	fn()
}

func contentsHandler(sha string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/contents/skills" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"s","sha":"` + sha + `","type":"dir"}]`))
	}
}

func mkSkill(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, name), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name, "SKILL.md"), []byte("# Skill\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestCheckUpToDateInSkipsNonGithubSource(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	entry := LockEntry{Source: "s", SourceType: "local", SkillPath: "skills/s/SKILL.md", FolderHash: "abc"}
	if got := checkUpToDateIn(dir, "s", entry); got != StatusUnknown {
		t.Errorf("got %v, want StatusUnknown", got)
	}
}

// outside: a lock entry without skillFolderHash predates the field; nothing to compare,
// so never a false badge.
func TestCheckUpToDateInSkipsMissingFolderHash(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	if got := checkUpToDateIn(dir, "s", entry); got != StatusUnknown {
		t.Errorf("got %v, want StatusUnknown", got)
	}
}

func TestCheckUpToDateInSkipsMissingLocalFile(t *testing.T) {
	dir := t.TempDir()
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "abc"}
	if got := checkUpToDateIn(dir, "missing", entry); got != StatusUnknown {
		t.Errorf("got %v, want StatusUnknown", got)
	}
}

func TestCheckUpToDateInMatchingFolderHashIsUpToDate(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "deadbeef"}
	withTestServer(t, contentsHandler("deadbeef"), func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUpToDate {
			t.Errorf("got %v, want StatusUpToDate", got)
		}
	})
}

func TestCheckUpToDateInDetectsOutdated(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "oldhash"}
	withTestServer(t, contentsHandler("newhash"), func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusOutdated {
			t.Errorf("got %v, want StatusOutdated", got)
		}
	})
}

// invariant: local SKILL.md drift is irrelevant; only the recorded folder hash counts.
func TestCheckUpToDateInIgnoresLocalContentDrift(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	if err := os.WriteFile(filepath.Join(dir, "s", "SKILL.md"), []byte("locally edited\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "samehash"}
	withTestServer(t, contentsHandler("samehash"), func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUpToDate {
			t.Errorf("got %v, want StatusUpToDate (folder hash is the oracle, not file contents)", got)
		}
	})
}

func TestCheckUpToDateInUnknownOn404(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "abc"}
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUnknown {
			t.Errorf("got %v, want StatusUnknown", got)
		}
	})
}

func TestCheckUpToDateInUnknownWhenFolderGoneUpstream(t *testing.T) {
	dir := t.TempDir()
	mkSkill(t, dir, "s")
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "abc"}
	withTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"other","sha":"x","type":"dir"}]`))
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUnknown {
			t.Errorf("got %v, want StatusUnknown", got)
		}
	})
}

func TestUpstreamFolderSHACachesPerParentDirectory(t *testing.T) {
	dir := t.TempDir()
	var hits int32
	handler := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		if r.URL.Path != "/repos/owner/repo/contents/skills" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"a","sha":"h1","type":"dir"},{"name":"b","sha":"h2","type":"dir"}]`))
	}

	withTestServer(t, handler, func() {
		mkSkill(t, dir, "a")
		mkSkill(t, dir, "b")
		a := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/a/SKILL.md", FolderHash: "h1"}
		b := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/b/SKILL.md", FolderHash: "stale"}

		if got := checkUpToDateIn(dir, "a", a); got != StatusUpToDate {
			t.Errorf("a: got %v, want StatusUpToDate", got)
		}
		if got := checkUpToDateIn(dir, "b", b); got != StatusOutdated {
			t.Errorf("b: got %v, want StatusOutdated", got)
		}
		if n := atomic.LoadInt32(&hits); n != 1 {
			t.Errorf("server hits = %d, want 1 (both skills share one listing)", n)
		}
	})
}

// invariant: concurrent callers share one fetch; run with -race.
func TestUpstreamFolderSHADedupesConcurrentFetches(t *testing.T) {
	dir := t.TempDir()
	var hits int32
	handler := func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"s","sha":"h","type":"dir"}]`))
	}

	withTestServer(t, handler, func() {
		mkSkill(t, dir, "s")
		entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md", FolderHash: "h"}

		var wg sync.WaitGroup
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if got := checkUpToDateIn(dir, "s", entry); got != StatusUpToDate {
					t.Errorf("got %v, want StatusUpToDate", got)
				}
			}()
		}
		wg.Wait()

		if n := atomic.LoadInt32(&hits); n != 1 {
			t.Errorf("server hits = %d, want 1 (concurrent callers share one fetch)", n)
		}
	})
}
