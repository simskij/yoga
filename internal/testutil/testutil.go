package testutil

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
	"golang.org/x/crypto/ssh"
)

// TempHome creates a temporary directory, sets $HOME to it, and returns the path.
// The original $HOME is restored when the test completes.
func TempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	orig := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	t.Cleanup(func() { os.Setenv("HOME", orig) })
	return dir
}

// GenerateAgeIdentity creates a temporary age identity file and returns the
// identity and its path.
func GenerateAgeIdentity(t *testing.T) (*age.X25519Identity, string) {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("GenerateAgeIdentity: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "age.key")
	if err := os.WriteFile(path, []byte(id.String()+"\n"), 0o600); err != nil {
		t.Fatalf("GenerateAgeIdentity write: %v", err)
	}
	return id, path
}

// GenerateSSHEd25519Key generates an SSH ed25519 key pair and writes the
// private key to path (creating parent dirs as needed).
func GenerateSSHEd25519Key(t *testing.T, path string) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateSSHEd25519Key: %v", err)
	}
	_ = pub

	sshPriv, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatalf("GenerateSSHEd25519Key signer: %v", err)
	}
	_ = sshPriv

	privPEM, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		t.Fatalf("GenerateSSHEd25519Key marshal: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("GenerateSSHEd25519Key mkdir: %v", err)
	}

	var buf strings.Builder
	if err := pem.Encode(&buf, privPEM); err != nil {
		t.Fatalf("GenerateSSHEd25519Key pem: %v", err)
	}

	if err := os.WriteFile(path, []byte(buf.String()), 0o600); err != nil {
		t.Fatalf("GenerateSSHEd25519Key write: %v", err)
	}
}

// WriteDotfile creates a file at the given path relative to the dotfiles root,
// creating parent directories as needed.
func WriteDotfile(t *testing.T, dotsRoot, rel, content string) string {
	t.Helper()
	full := filepath.Join(dotsRoot, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("WriteDotfile mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteDotfile write: %v", err)
	}
	return full
}
