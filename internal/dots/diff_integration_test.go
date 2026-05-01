//go:build integration

package dots_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/simskij/yo/internal/dots"
	"github.com/simskij/yo/internal/testutil"
)

func TestDiff_NoDifferences(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "zsh config")
	dst := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	var buf strings.Builder
	if err := dots.Diff(dotsRoot, "mymachine", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no differences message, got: %s", buf.String())
	}
}

func TestDiff_WithDifference(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "new content")
	dst := filepath.Join(home, ".zshrc")
	if err := os.WriteFile(dst, []byte("old content"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = src

	var buf strings.Builder
	if err := dots.Diff(dotsRoot, "mymachine", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "- old content") || !strings.Contains(out, "+ new content") {
		t.Errorf("expected diff output, got: %s", out)
	}
}

func TestDiff_SkipsMissingTargets(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "content")

	var buf strings.Builder
	if err := dots.Diff(dotsRoot, "mymachine", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("missing file should be skipped, got: %s", buf.String())
	}
}

func TestInit_DoesNotOverwriteReadme(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	readmePath := filepath.Join(dotsRoot, "README.md")
	if err := os.WriteFile(readmePath, []byte("custom readme"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := dots.Init(dotsRoot, "mymachine"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom readme" {
		t.Errorf("readme was overwritten, got: %s", data)
	}
}
