package crypt

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	"filippo.io/age/agessh"
)

var candidatePaths = []string{
	"~/.ssh/id_ed25519",
	"~/.ssh/id_rsa",
}

// LoadIdentities loads age identities from the given paths, falling back to
// well-known SSH key locations. Extra paths are tried first.
func LoadIdentities(extraPaths ...string) ([]age.Identity, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var paths []string
	paths = append(paths, extraPaths...)
	for _, p := range candidatePaths {
		if len(p) > 1 && p[:2] == "~/" {
			p = filepath.Join(home, p[2:])
		}
		paths = append(paths, p)
	}

	var ids []age.Identity
	for _, p := range paths {
		loaded, err := loadIdentityFile(p)
		if err != nil {
			continue
		}
		ids = append(ids, loaded...)
	}

	if len(ids) == 0 {
		return nil, errors.New("no age identities found; set dots.encryption.identity in ~/.config/yo/config.yaml or ensure ~/.ssh/id_ed25519 or ~/.ssh/id_rsa exists")
	}
	return ids, nil
}

func loadIdentityFile(path string) ([]age.Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Try native age identity first
	if ids, err := age.ParseIdentities(bytes.NewReader(data)); err == nil && len(ids) > 0 {
		return ids, nil
	}

	// Try SSH identity
	id, err := agessh.ParseIdentity(data)
	if err != nil {
		return nil, err
	}
	return []age.Identity{id}, nil
}

// Encrypt encrypts plaintext to dst using recipients derived from ids.
func Encrypt(dst io.Writer, plaintext []byte, ids []age.Identity) error {
	recipients, err := identitiesToRecipients(ids)
	if err != nil {
		return err
	}
	w, err := age.Encrypt(dst, recipients...)
	if err != nil {
		return err
	}
	if _, err := w.Write(plaintext); err != nil {
		return err
	}
	return w.Close()
}

// Decrypt decrypts src using ids and returns the plaintext.
func Decrypt(src io.Reader, ids []age.Identity) ([]byte, error) {
	r, err := age.Decrypt(src, ids...)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

// EncryptFile encrypts the file at src and writes the result to dst.
func EncryptFile(src, dst string, ids []age.Identity) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return Encrypt(f, data, ids)
}

// DecryptFile decrypts the file at src and writes plaintext to dst.
func DecryptFile(src, dst string, ids []age.Identity) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	defer f.Close()

	data, err := Decrypt(f, ids)
	if err != nil {
		return fmt.Errorf("decrypt %s: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o600)
}

func identitiesToRecipients(ids []age.Identity) ([]age.Recipient, error) {
	recipients := make([]age.Recipient, 0, len(ids))
	for _, id := range ids {
		switch id := id.(type) {
		case *age.X25519Identity:
			recipients = append(recipients, id.Recipient())
		case *agessh.Ed25519Identity:
			recipients = append(recipients, id.Recipient())
		case *agessh.RSAIdentity:
			recipients = append(recipients, id.Recipient())
		default:
			return nil, fmt.Errorf("unsupported identity type %T", id)
		}
	}
	if len(recipients) == 0 {
		return nil, errors.New("no recipients derived from identities")
	}
	return recipients, nil
}
