package dots

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"github.com/simskij/yo/internal/crypt"
)

const encExt = ".age"

func isEncrypted(src string) bool {
	return strings.HasSuffix(src, encExt)
}

// AddEncrypted moves filePath into the dotfiles repo layer, encrypts it, and
// removes the plaintext original. The destination file in the repo gets .age extension.
func AddEncrypted(dotsRoot, layer, filePath string, ids []age.Identity) error {
	dotsRoot, err := expandPath(dotsRoot)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if !filepath.IsAbs(filePath) || !strings.HasPrefix(filePath, home+string(filepath.Separator)) {
		return fmt.Errorf("file must be under $HOME")
	}

	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", filePath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("directories are not supported; add files individually")
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%s is already managed by yo", filePath)
	}

	rel, err := filepath.Rel(home, filePath)
	if err != nil {
		return err
	}

	dst := filepath.Join(dotsRoot, layer, rel+encExt)
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("%s already exists in the %s layer", rel, layer)
	}

	if err := crypt.EncryptFile(filePath, dst, ids); err != nil {
		return err
	}

	return os.Remove(filePath)
}

// ApplyEncrypted decrypts src into dst (copy, not symlink).
func ApplyEncrypted(src, dst string, ids []age.Identity) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return crypt.DecryptFile(src, dst, ids)
}

// EditEncrypted decrypts src to a temp file, opens $EDITOR, then re-encrypts.
func EditEncrypted(src string, ids []age.Identity) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "yo-edit-*")
	if err != nil {
		return err
	}
	tmpPath := f.Name()
	defer os.Remove(tmpPath)

	decrypted, err := crypt.Decrypt(f, ids)
	_ = f.Close()
	if err != nil {
		// src may not be encrypted yet — treat raw
		decrypted = data
	}

	if err := os.WriteFile(tmpPath, decrypted, 0o600); err != nil {
		return err
	}

	if err := openEditor(tmpPath); err != nil {
		return err
	}

	return crypt.EncryptFile(tmpPath, src, ids)
}
