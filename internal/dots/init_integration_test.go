//go:build integration

package dots_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/dots"
	"github.com/simskij/yo/internal/testutil"
)

func TestInit_CreatesDirectories(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	if err := dots.Init(dotsRoot, "mymachine"); err != nil {
		t.Fatal(err)
	}

	for _, dir := range []string{"global", "mymachine"} {
		path := filepath.Join(dotsRoot, dir)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected dir %s to exist: %v", path, err)
		}
	}
}

func TestInit_WritesReadme(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	if err := dots.Init(dotsRoot, "mymachine"); err != nil {
		t.Fatal(err)
	}

	readme := filepath.Join(dotsRoot, "README.md")
	if _, err := os.Stat(readme); err != nil {
		t.Errorf("expected README.md at %s: %v", readme, err)
	}
}

func TestInit_Idempotent(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	if err := dots.Init(dotsRoot, "mymachine"); err != nil {
		t.Fatal(err)
	}
	if err := dots.Init(dotsRoot, "mymachine"); err != nil {
		t.Errorf("second init should be a no-op, got error: %v", err)
	}
}
