package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempHome creates a temporary directory, sets $HOME to it, and returns the path.
// The original $HOME is restored when the test completes.
func TempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	orig := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	t.Cleanup(func() { os.Setenv("HOME", orig) })
	return dir
}

// WriteDotfile creates a file at the given path relative to the dotfiles root,
// creating parent directories as needed.
func WriteDotfile(t *testing.T, dotsRoot, rel, content string) string {
	t.Helper()
	full := filepath.Join(dotsRoot, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("WriteDotfile mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteDotfile write: %v", err)
	}
	return full
}
