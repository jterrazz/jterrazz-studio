package skill

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
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
	if err := os.WriteFile(path, []byte(fixture), 0o644); err != nil {
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
	}
	if entry != want {
		t.Errorf("entry = %+v, want %+v", entry, want)
	}
}

// withTestServer points rawContentBaseURL at a local httptest server for the
// duration of fn, restoring it afterwards — CheckUpToDate must never hit the
// real network in tests.
func withTestServer(t *testing.T, handler http.HandlerFunc, fn func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	defer srv.Close()

	prevBase := rawContentBaseURL
	rawContentBaseURL = srv.URL
	defer func() { rawContentBaseURL = prevBase }()

	fn()
}

func TestCheckUpToDateInSkipsNonGithubSource(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "s", true)
	entry := LockEntry{Source: "s", SourceType: "local", SkillPath: "skills/s/SKILL.md"}
	if got := checkUpToDateIn(dir, "s", entry); got != StatusUnknown {
		t.Errorf("got %v, want StatusUnknown", got)
	}
}

func TestCheckUpToDateInSkipsMissingLocalFile(t *testing.T) {
	dir := t.TempDir()
	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	if got := checkUpToDateIn(dir, "missing", entry); got != StatusUnknown {
		t.Errorf("got %v, want StatusUnknown", got)
	}
}

func TestCheckUpToDateInMatchesUpToDate(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "s", false)
	content := "# Skill\n\nSame content.\n"
	if err := os.WriteFile(filepath.Join(dir, "s", "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(content))
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUpToDate {
			t.Errorf("got %v, want StatusUpToDate", got)
		}
	})
}

func TestCheckUpToDateInIgnoresTrailingWhitespaceDiffs(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "s", false)
	if err := os.WriteFile(filepath.Join(dir, "s", "SKILL.md"), []byte("line one \nline two\n\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("line one\nline two"))
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUpToDate {
			t.Errorf("got %v, want StatusUpToDate (trailing whitespace should be ignored)", got)
		}
	})
}

func TestCheckUpToDateInDetectsOutdated(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "s", false)
	if err := os.WriteFile(filepath.Join(dir, "s", "SKILL.md"), []byte("old content\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("new content\n"))
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusOutdated {
			t.Errorf("got %v, want StatusOutdated", got)
		}
	})
}

func TestCheckUpToDateInUnknownOn404(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "s", false)
	if err := os.WriteFile(filepath.Join(dir, "s", "SKILL.md"), []byte("content\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUnknown {
			t.Errorf("got %v, want StatusUnknown", got)
		}
	})
}

func TestCheckUpToDateInFallsBackToMainRef(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "s", false)
	content := "content\n"
	if err := os.WriteFile(filepath.Join(dir, "s", "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	entry := LockEntry{Source: "owner/repo", SourceType: "github", SkillPath: "skills/s/SKILL.md"}
	withTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Only succeed for the "main" ref, simulating a source where HEAD
		// doesn't resolve on raw.githubusercontent.com.
		if r.URL.Path == "/owner/repo/main/skills/s/SKILL.md" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(content))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}, func() {
		if got := checkUpToDateIn(dir, "s", entry); got != StatusUpToDate {
			t.Errorf("got %v, want StatusUpToDate via main-ref fallback", got)
		}
	})
}
