//go:build integration

package repo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/repo"
	"github.com/simskij/yo/internal/testutil"
)

func TestInit_CreatesDirectory(t *testing.T) {
	testutil.TempHome(t)
	yoPath := filepath.Join(t.TempDir(), "yofiles")

	if err := repo.Init(yoPath); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(yoPath); err != nil {
		t.Errorf("expected yoPath to exist: %v", err)
	}
}

func TestInit_WritesReadme(t *testing.T) {
	testutil.TempHome(t)
	yoPath := t.TempDir()

	if err := repo.Init(yoPath); err != nil {
		t.Fatal(err)
	}

	readme := filepath.Join(yoPath, "README.md")
	if _, err := os.Stat(readme); err != nil {
		t.Errorf("expected README.md at %s: %v", readme, err)
	}
}

func TestInit_InitializesGitRepo(t *testing.T) {
	testutil.TempHome(t)
	yoPath := t.TempDir()

	if err := repo.Init(yoPath); err != nil {
		t.Fatal(err)
	}

	gitDir := filepath.Join(yoPath, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("expected .git dir at %s: %v", gitDir, err)
	}
}

func TestInit_Idempotent(t *testing.T) {
	testutil.TempHome(t)
	yoPath := t.TempDir()

	if err := repo.Init(yoPath); err != nil {
		t.Fatal(err)
	}
	if err := repo.Init(yoPath); err != nil {
		t.Errorf("second init should be a no-op, got error: %v", err)
	}
}
