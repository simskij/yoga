package crypt_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/simskij/yo/internal/crypt"
	"github.com/simskij/yo/internal/testutil"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	_, identityFile := testutil.GenerateAgeIdentity(t)

	ids, err := crypt.LoadIdentities(identityFile)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("super secret content")

	var buf bytes.Buffer
	if err := crypt.Encrypt(&buf, plaintext, ids); err != nil {
		t.Fatal(err)
	}

	got, err := crypt.Decrypt(&buf, ids)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, plaintext) {
		t.Errorf("got %q, want %q", got, plaintext)
	}
}

func TestLoadIdentities_NoneFound(t *testing.T) {
	testutil.TempHome(t)
	_, err := crypt.LoadIdentities()
	if err == nil {
		t.Error("expected error when no identities found, got nil")
	}
}

func TestLoadIdentities_SSHEd25519(t *testing.T) {
	home := testutil.TempHome(t)
	testutil.GenerateSSHEd25519Key(t, filepath.Join(home, ".ssh", "id_ed25519"))

	ids, err := crypt.LoadIdentities()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 0 {
		t.Error("expected at least one identity from ~/.ssh/id_ed25519")
	}
}

func TestLoadIdentities_ExplicitPath(t *testing.T) {
	testutil.TempHome(t)
	_, keyFile := testutil.GenerateAgeIdentity(t)

	ids, err := crypt.LoadIdentities(keyFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 0 {
		t.Error("expected identity from explicit path")
	}
}

func TestLoadIdentities_SkipsBadKey(t *testing.T) {
	home := testutil.TempHome(t)
	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sshDir, "id_ed25519"), []byte("not a key"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, keyFile := testutil.GenerateAgeIdentity(t)

	ids, err := crypt.LoadIdentities(keyFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 0 {
		t.Error("expected identity from fallback age key")
	}
}

func TestEncryptFile_DecryptFile(t *testing.T) {
	_, identityFile := testutil.GenerateAgeIdentity(t)

	ids, err := crypt.LoadIdentities(identityFile)
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	src := filepath.Join(dir, "secret.txt")
	enc := filepath.Join(dir, "secret.txt.age")
	dst := filepath.Join(dir, "decrypted.txt")

	if err := os.WriteFile(src, []byte("my secret"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := crypt.EncryptFile(src, enc, ids); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(enc); err != nil {
		t.Fatalf("expected encrypted file at %s", enc)
	}

	if err := crypt.DecryptFile(enc, dst, ids); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "my secret" {
		t.Errorf("got %q, want %q", data, "my secret")
	}
}
