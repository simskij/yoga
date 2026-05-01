package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/testutil"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Yo.Path != "~/.yofiles" {
		t.Errorf("expected ~/.yofiles, got %s", cfg.Yo.Path)
	}
}

func TestExpandPath_Tilde(t *testing.T) {
	home := testutil.TempHome(t)
	got, err := config.ExpandPath("~/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "foo/bar")
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestExpandPath_Absolute(t *testing.T) {
	got, err := config.ExpandPath("/absolute/path")
	if err != nil {
		t.Fatal(err)
	}
	if got != "/absolute/path" {
		t.Errorf("got %s, want /absolute/path", got)
	}
}

func TestExpandPath_NoTilde(t *testing.T) {
	got, err := config.ExpandPath("relative/path")
	if err != nil {
		t.Fatal(err)
	}
	if got != "relative/path" {
		t.Errorf("got %s, want relative/path", got)
	}
}

func TestLoad_NotExist(t *testing.T) {
	testutil.TempHome(t)
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg != nil {
		t.Errorf("expected nil config when file does not exist, got %+v", cfg)
	}
}

func TestSaveAndLoad(t *testing.T) {
	testutil.TempHome(t)

	want := &config.Config{}
	want.Yo.Path = "/tmp/my-yofiles"

	if err := config.Save(want); err != nil {
		t.Fatal(err)
	}

	got, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected config, got nil")
	}
	if got.Yo.Path != want.Yo.Path {
		t.Errorf("got %s, want %s", got.Yo.Path, want.Yo.Path)
	}
}

func TestExists_False(t *testing.T) {
	testutil.TempHome(t)
	exists, err := config.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("expected false when config does not exist")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	home := testutil.TempHome(t)
	cfgDir := filepath.Join(home, ".config", "yo")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("key: {unclosed"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load()
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestExists_True(t *testing.T) {
	home := testutil.TempHome(t)
	cfgDir := filepath.Join(home, ".config", "yo")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("yo:\n  path: ~/.yofiles\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	exists, err := config.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Error("expected true when config exists")
	}
}
