package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestSkillMDIsInSyncWithConfig keeps the committed projection honest: if
// src/internal/config changes without a `make skills` run, this fails and
// tells the developer what to do about it.
func TestSkillMDIsInSyncWithConfig(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file location")
	}
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
	path := filepath.Join(repoRoot, defaultRelPath)

	// Given: the committed projection
	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	// When: re-rendering the generated blocks from the live registries
	rendered, err := render(string(current))
	if err != nil {
		t.Fatalf("rendering: %v", err)
	}

	// Then: the committed file must match, or the projection has drifted
	if rendered != string(current) {
		t.Fatalf("%s is out of sync with src/internal/config — run `make skills` and commit the result", defaultRelPath)
	}
}
