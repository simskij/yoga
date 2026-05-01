package dots

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/testutil"
)

func TestRemove_RejectsUnmanagedFile(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	target := filepath.Join(home, ".zshrc")
	if err := os.WriteFile(target, []byte("zsh"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := FindLayers(dotsRoot, target)
	if err != nil {
		t.Fatal(err)
	}

	err = Remove(dotsRoot, target, "global")
	if err == nil {
		t.Error("expected error for unmanaged file, got nil")
	}
}

func TestRemove_RestoresFileFromGlobal(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "zsh config")
	dst := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	if err := Remove(dotsRoot, dst, "global"); err != nil {
		t.Fatal(err)
	}

	// symlink should be gone, regular file should be in place
	info, err := os.Lstat(dst)
	if err != nil {
		t.Fatalf("expected file at %s: %v", dst, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("expected regular file, got symlink")
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "zsh config" {
		t.Errorf("expected 'zsh config', got %s", data)
	}

	// repo copy should be gone
	if _, err := os.Stat(src); err == nil {
		t.Error("expected repo copy to be removed")
	}
}

func TestRemove_RestoresFileFromMachineLayer(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "mymachine/.zshrc", "machine zsh")
	dst := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	if err := Remove(dotsRoot, dst, "mymachine"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "machine zsh" {
		t.Errorf("expected 'machine zsh', got %s", data)
	}
}

func TestFindLayers_SingleLayer(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "zsh config")
	dst := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	layers, err := FindLayers(dotsRoot, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(layers) != 1 || layers[0] != "global" {
		t.Errorf("expected [global], got %v", layers)
	}
}

func TestFindLayers_BothLayers(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "global zsh")
	testutil.WriteDotfile(t, dotsRoot, "mymachine/.zshrc", "machine zsh")

	// symlink points to machine layer (machine wins)
	src := filepath.Join(dotsRoot, "mymachine", ".zshrc")
	dst := filepath.Join(home, ".zshrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	layers, err := FindLayers(dotsRoot, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(layers) != 2 {
		t.Errorf("expected 2 layers, got %v", layers)
	}
}

func TestFindLayers_UnmanagedFile(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	dst := filepath.Join(home, ".zshrc")
	if err := os.WriteFile(dst, []byte("zsh"), 0o644); err != nil {
		t.Fatal(err)
	}

	layers, err := FindLayers(dotsRoot, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(layers) != 0 {
		t.Errorf("expected no layers for unmanaged file, got %v", layers)
	}
}
