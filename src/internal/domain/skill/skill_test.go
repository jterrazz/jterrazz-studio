package skill

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// mkSkillDir creates dir/name (optionally as a symlink to a real directory
// elsewhere) with a SKILL.md unless withSkillMd is false.
func mkSkillDir(t *testing.T, base, name string, withSkillMd bool) {
	t.Helper()
	path := filepath.Join(base, name)
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if withSkillMd {
		if err := os.WriteFile(filepath.Join(path, "SKILL.md"), []byte("# skill\n"), 0o644); err != nil {
			t.Fatalf("write SKILL.md: %v", err)
		}
	}
}

func TestListInstalledInMissingDir(t *testing.T) {
	got := listInstalledIn(filepath.Join(t.TempDir(), "does-not-exist"))
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestListInstalledInEmptyPath(t *testing.T) {
	got := listInstalledIn("")
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestListInstalledInFindsSkillDirs(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "alpha", true)
	mkSkillDir(t, dir, "beta", true)

	got := listInstalledIn(dir)
	sort.Strings(got)
	want := []string{"alpha", "beta"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestListInstalledInSkipsDirsWithoutSkillMd(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "has-skill", true)
	mkSkillDir(t, dir, "no-skill", false)

	got := listInstalledIn(dir)
	if len(got) != 1 || got[0] != "has-skill" {
		t.Errorf("got %v, want [has-skill]", got)
	}
}

func TestListInstalledInSkipsFilesAndDotfiles(t *testing.T) {
	dir := t.TempDir()
	mkSkillDir(t, dir, "real-skill", true)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("not a skill"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, ".hidden"), 0o755); err != nil {
		t.Fatalf("mkdir hidden: %v", err)
	}

	got := listInstalledIn(dir)
	if len(got) != 1 || got[0] != "real-skill" {
		t.Errorf("got %v, want [real-skill]", got)
	}
}

func TestListInstalledInFollowsSymlinkedSkillDir(t *testing.T) {
	dir := t.TempDir()
	realDir := t.TempDir()
	mkSkillDir(t, realDir, "linked-skill", true)

	if err := os.Symlink(filepath.Join(realDir, "linked-skill"), filepath.Join(dir, "linked-skill")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	got := listInstalledIn(dir)
	if len(got) != 1 || got[0] != "linked-skill" {
		t.Errorf("got %v, want [linked-skill]", got)
	}
}
