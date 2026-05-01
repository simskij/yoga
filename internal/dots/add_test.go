package dots

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/testutil"
)

func TestAdd_RejectsDirectory(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	dir := filepath.Join(home, "mydir")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := Add(dotsRoot, "global", dir)
	if err == nil {
		t.Error("expected error for directory, got nil")
	}
}

func TestAdd_RejectsAlreadyManagedSymlink(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "zsh config")
	dst := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	err := Add(dotsRoot, "global", dst)
	if err == nil {
		t.Error("expected error for already-managed symlink, got nil")
	}
}

func TestAdd_RejectsIfAlreadyInRepo(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "existing")
	target := filepath.Join(home, ".zshrc")
	if err := os.WriteFile(target, []byte("local"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Add(dotsRoot, "global", target)
	if err == nil {
		t.Error("expected error when file already exists in repo, got nil")
	}
}

func TestAdd_MovesFileAndSymlinks(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	target := filepath.Join(home, ".zshrc")
	if err := os.WriteFile(target, []byte("zsh config"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Add(dotsRoot, "global", target); err != nil {
		t.Fatal(err)
	}

	// original location should now be a symlink
	linkTarget, err := os.Readlink(target)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", target, err)
	}

	expectedSrc := filepath.Join(dotsRoot, "global", ".zshrc")
	if linkTarget != expectedSrc {
		t.Errorf("symlink points to %s, want %s", linkTarget, expectedSrc)
	}

	// content should be preserved in the repo
	data, err := os.ReadFile(linkTarget)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "zsh config" {
		t.Errorf("expected content 'zsh config', got %s", data)
	}
}

func TestAdd_RejectsFileOutsideHome(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	outsideFile := filepath.Join(dotsRoot, "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Add(dotsRoot, "global", outsideFile)
	if err == nil {
		t.Error("expected error for file outside $HOME, got nil")
	}
}

func TestAdd_RejectsRelativePath(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	err := Add(dotsRoot, "global", "relative/path/file.txt")
	if err == nil {
		t.Error("expected error for relative path, got nil")
	}
}

func TestAdd_RejectsPathWithHomePrefixButOutside(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	// e.g. home=/tmp/abc, sibling=/tmp/abcmalicious/file
	sibling := home + "malicious"
	if err := os.MkdirAll(sibling, 0o755); err != nil {
		t.Fatal(err)
	}
	outsideFile := filepath.Join(sibling, "file.txt")
	if err := os.WriteFile(outsideFile, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Add(dotsRoot, "global", outsideFile)
	if err == nil {
		t.Errorf("expected error for path outside $HOME, got nil (home=%s, file=%s)", home, outsideFile)
	}
}

func TestAdd_HandlesNestedPath(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	target := filepath.Join(home, ".config", "nvim", "init.lua")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("nvim config"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Add(dotsRoot, "global", target); err != nil {
		t.Fatal(err)
	}

	expectedSrc := filepath.Join(dotsRoot, "global", ".config", "nvim", "init.lua")
	linkTarget, err := os.Readlink(target)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", target, err)
	}
	if linkTarget != expectedSrc {
		t.Errorf("symlink points to %s, want %s", linkTarget, expectedSrc)
	}
}
