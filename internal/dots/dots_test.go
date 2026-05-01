package dots

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/testutil"
)

func TestExpandPath_Tilde(t *testing.T) {
	home := testutil.TempHome(t)
	got, err := expandPath("~/foo")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "foo")
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestExpandPath_NoTilde(t *testing.T) {
	got, err := expandPath("/absolute")
	if err != nil {
		t.Fatal(err)
	}
	if got != "/absolute" {
		t.Errorf("got %s, want /absolute", got)
	}
}

func TestFileStatus_Missing(t *testing.T) {
	home := testutil.TempHome(t)
	src := filepath.Join(home, "src")
	dst := filepath.Join(home, "dst")
	if got := fileStatus(src, dst); got != StatusMissing {
		t.Errorf("got %v, want StatusMissing", got)
	}
}

func TestFileStatus_Linked(t *testing.T) {
	home := testutil.TempHome(t)
	src := filepath.Join(home, "src")
	dst := filepath.Join(home, "dst")

	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}

	if got := fileStatus(src, dst); got != StatusLinked {
		t.Errorf("got %v, want StatusLinked", got)
	}
}

func TestFileStatus_Differs_WrongTarget(t *testing.T) {
	home := testutil.TempHome(t)
	src := filepath.Join(home, "src")
	other := filepath.Join(home, "other")
	dst := filepath.Join(home, "dst")

	for _, f := range []string{src, other} {
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(other, dst); err != nil {
		t.Fatal(err)
	}

	if got := fileStatus(src, dst); got != StatusDiffers {
		t.Errorf("got %v, want StatusDiffers", got)
	}
}

func TestFileStatus_Differs_RegularFile(t *testing.T) {
	home := testutil.TempHome(t)
	src := filepath.Join(home, "src")
	dst := filepath.Join(home, "dst")

	for _, f := range []string{src, dst} {
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if got := fileStatus(src, dst); got != StatusDiffers {
		t.Errorf("got %v, want StatusDiffers", got)
	}
}
