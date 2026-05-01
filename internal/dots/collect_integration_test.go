//go:build integration

package dots_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/dots"
	"github.com/simskij/yo/internal/testutil"
)

func TestCollect_GlobalOnly(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "zsh config")

	entries, err := dots.Collect(dotsRoot, "mymachine")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	want := filepath.Join(home, ".zshrc")
	if entries[0].Dst != want {
		t.Errorf("got dst %s, want %s", entries[0].Dst, want)
	}
	if entries[0].Status != dots.StatusMissing {
		t.Errorf("expected StatusMissing, got %v", entries[0].Status)
	}
}

func TestCollect_MachineOverridesGlobal(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "global zsh")
	testutil.WriteDotfile(t, dotsRoot, "mymachine/.zshrc", "machine zsh")

	entries, err := dots.Collect(dotsRoot, "mymachine")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (machine overrides global), got %d", len(entries))
	}

	machineSrc := filepath.Join(dotsRoot, "mymachine", ".zshrc")
	if entries[0].Src != machineSrc {
		t.Errorf("expected machine src %s, got %s", machineSrc, entries[0].Src)
	}
}

func TestCollect_BothLayers(t *testing.T) {
	testutil.TempHome(t)
	dotsRoot := t.TempDir()

	testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "global zsh")
	testutil.WriteDotfile(t, dotsRoot, "mymachine/.gitconfig", "machine git")

	entries, err := dots.Collect(dotsRoot, "mymachine")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestApplyEntry_CreatesSymlink(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "zsh config")
	dst := filepath.Join(home, ".zshrc")

	e := dots.Entry{Src: src, Dst: dst, Status: dots.StatusMissing}
	if err := dots.ApplyEntry(e); err != nil {
		t.Fatal(err)
	}

	target, err := os.Readlink(dst)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", dst, err)
	}
	if target != src {
		t.Errorf("symlink points to %s, want %s", target, src)
	}
}

func TestApplyEntry_BacksUpExisting(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.zshrc", "new zsh")
	dst := filepath.Join(home, ".zshrc")
	if err := os.WriteFile(dst, []byte("old zsh"), 0o644); err != nil {
		t.Fatal(err)
	}

	e := dots.Entry{Src: src, Dst: dst, Status: dots.StatusDiffers}
	if err := dots.ApplyEntry(e); err != nil {
		t.Fatal(err)
	}

	bak := dst + ".bak"
	data, err := os.ReadFile(bak)
	if err != nil {
		t.Fatalf("expected backup at %s: %v", bak, err)
	}
	if string(data) != "old zsh" {
		t.Errorf("backup content wrong: %s", data)
	}
}

func TestApplyEntry_CreatesParentDirs(t *testing.T) {
	home := testutil.TempHome(t)
	dotsRoot := t.TempDir()

	src := testutil.WriteDotfile(t, dotsRoot, "global/.config/nvim/init.lua", "nvim config")
	dst := filepath.Join(home, ".config", "nvim", "init.lua")

	e := dots.Entry{Src: src, Dst: dst, Status: dots.StatusMissing}
	if err := dots.ApplyEntry(e); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Lstat(dst); err != nil {
		t.Errorf("expected symlink at %s: %v", dst, err)
	}
}
