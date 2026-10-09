package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSkillMDIsInSyncWithConfig(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file location")
	}
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
	path := filepath.Join(repoRoot, defaultRelPath)

	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	rendered, err := render(string(current))
	if err != nil {
		t.Fatalf("rendering: %v", err)
	}

	if rendered != string(current) {
		t.Fatalf("%s is out of sync with src/internal/config — run `make skills` and commit the result", defaultRelPath)
	}
}
