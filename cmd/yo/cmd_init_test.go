//go:build integration

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/testutil"
)

func TestRunInit_DryRun_DoesNotWriteConfig(t *testing.T) {
	testutil.TempHome(t)
	initYoPathFlag = ""

	if _, err := runInit(true); err != nil {
		t.Fatal(err)
	}

	exists, err := config.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("dry-run should not write config file")
	}
}

func TestRunInit_DryRun_DoesNotCreateRepo(t *testing.T) {
	home := testutil.TempHome(t)
	initYoPathFlag = ""

	if _, err := runInit(true); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(home, ".yofiles")); !os.IsNotExist(err) {
		t.Error("dry-run should not create the yofiles directory")
	}
}

func TestRunInit_DryRun_WhenConfigExists_ReturnsExistingPath(t *testing.T) {
	testutil.TempHome(t)
	initYoPathFlag = ""

	cfg := config.DefaultConfig()
	cfg.Yo.Path = "~/.my-yofiles"
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}

	yoPath, err := runInit(true)
	if err != nil {
		t.Fatal(err)
	}
	if yoPath != "~/.my-yofiles" {
		t.Errorf("got %q, want ~/.my-yofiles", yoPath)
	}
}

func TestRunInit_Normal_WritesConfig(t *testing.T) {
	testutil.TempHome(t)
	initYoPathFlag = t.TempDir()

	if _, err := runInit(false); err != nil {
		t.Fatal(err)
	}

	exists, err := config.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Error("expected config to be written")
	}
}

func TestRunInit_Normal_InitialisesRepo(t *testing.T) {
	yoPath := t.TempDir()
	testutil.TempHome(t)
	initYoPathFlag = yoPath

	if _, err := runInit(false); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(yoPath, ".git")); err != nil {
		t.Errorf("expected .git dir at %s: %v", yoPath, err)
	}
}
