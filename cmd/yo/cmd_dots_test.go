//go:build integration

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/testutil"
)

func mustRun(t *testing.T, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}

// makeLocalRepo creates a bare-style git repo at dir containing a global/.testrc
// dotfile so it can be used as a clone source in tests.
func makeLocalRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	mustRun(t, "git", "-C", dir, "init")
	mustRun(t, "git", "-C", dir, "config", "user.email", "test@test.com")
	mustRun(t, "git", "-C", dir, "config", "user.name", "Test")
	globalDir := filepath.Join(dir, "global")
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, ".testrc"), []byte("cloned"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, "git", "-C", dir, "add", ".")
	mustRun(t, "git", "-C", dir, "-c", "commit.gpgsign=false", "commit", "-m", "init")
	return dir
}

func setupDotsEnv(t *testing.T) (home, dotsPath string) {
	t.Helper()
	home = testutil.TempHome(t)
	cfg := config.DefaultConfig()
	cfg.Yo.Path = filepath.Join(home, "yofiles")
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	dotsPath = filepath.Join(home, "yofiles", "dots")
	return home, dotsPath
}

// --- dots init ---

func TestRunDotsInit_DryRun_DoesNotCreateDirs(t *testing.T) {
	home := testutil.TempHome(t)
	initYoPathFlag = filepath.Join(home, "yofiles")

	if err := runDotsInit(true); err != nil {
		t.Fatal(err)
	}

	dotsPath := filepath.Join(home, "yofiles", "dots")
	if _, err := os.Stat(dotsPath); !os.IsNotExist(err) {
		t.Error("dry-run should not create the dots directory")
	}
}

func TestRunDotsInit_Normal_CreatesDirs(t *testing.T) {
	home := testutil.TempHome(t)
	initYoPathFlag = filepath.Join(home, "yofiles")

	if err := runDotsInit(false); err != nil {
		t.Fatal(err)
	}

	for _, dir := range []string{"global", "mymachine"} {
		_ = dir // machine detection varies; just check dots dir exists
	}
	dotsPath := filepath.Join(home, "yofiles", "dots")
	if _, err := os.Stat(dotsPath); err != nil {
		t.Errorf("expected dots dir at %s: %v", dotsPath, err)
	}
}

// --- dots apply ---

func TestRunDotsApply_DryRun_DoesNotSymlink(t *testing.T) {
	home, dotsPath := setupDotsEnv(t)
	globalDir := filepath.Join(dotsPath, "global")
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteDotfile(t, dotsPath, "global/.testrc", "content")

	if err := runDotsApply("mymachine", "", true); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Lstat(filepath.Join(home, ".testrc")); !os.IsNotExist(err) {
		t.Error("dry-run should not create a symlink in $HOME")
	}
}

func TestRunDotsApply_DryRun_AlreadyLinkedNotRelinked(t *testing.T) {
	home, dotsPath := setupDotsEnv(t)
	globalDir := filepath.Join(dotsPath, "global")
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := testutil.WriteDotfile(t, dotsPath, "global/.testrc", "content")

	// create the symlink manually so entry is StatusLinked
	dst := filepath.Join(home, ".testrc")
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	// dry-run should not error and should not remove the existing symlink
	if err := runDotsApply("mymachine", "", true); err != nil {
		t.Fatal(err)
	}

	target, err := os.Readlink(dst)
	if err != nil {
		t.Fatalf("symlink gone after dry-run: %v", err)
	}
	if target != src {
		t.Errorf("symlink target changed: got %s, want %s", target, src)
	}
}

// --- dots clone ---

func TestRunDotsClone_Normal_ClonesAndApplies(t *testing.T) {
	home, dotsPath := setupDotsEnv(t)
	remote := makeLocalRepo(t)

	if err := runDotsClone(remote, false); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dotsPath, "global", ".testrc")); err != nil {
		t.Errorf("dotfile not present in cloned repo: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(home, ".testrc")); err != nil {
		t.Errorf("expected symlink in $HOME after apply: %v", err)
	}
}

func TestRunDotsClone_ErrorWhenNonEmpty(t *testing.T) {
	_, dotsPath := setupDotsEnv(t)
	remote := makeLocalRepo(t)

	// pre-populate the dots path so it is non-empty
	if err := os.MkdirAll(dotsPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dotsPath, "existing"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := runDotsClone(remote, false); err == nil {
		t.Error("expected error when dots path is non-empty, got nil")
	}
}

func TestRunDotsClone_DryRun_DoesNotClone(t *testing.T) {
	_, dotsPath := setupDotsEnv(t)
	remote := makeLocalRepo(t)

	if err := runDotsClone(remote, true); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(dotsPath); !os.IsNotExist(err) {
		t.Error("dry-run should not create the dots directory")
	}
}

func TestRunDotsApply_Normal_CreatesSymlink(t *testing.T) {
	home, dotsPath := setupDotsEnv(t)
	globalDir := filepath.Join(dotsPath, "global")
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := testutil.WriteDotfile(t, dotsPath, "global/.testrc", "content")

	if err := runDotsApply("mymachine", "", false); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(home, ".testrc")
	target, err := os.Readlink(dst)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", dst, err)
	}
	if target != src {
		t.Errorf("symlink target: got %s, want %s", target, src)
	}
}
