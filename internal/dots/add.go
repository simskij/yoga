package dots

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Add(dotsRoot, layer, filePath string) error {
	dotsRoot, err := expandPath(dotsRoot)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", filePath, err)
	}

	if info.IsDir() {
		return errors.New("directories are not supported; add files individually")
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(filePath)
		if err != nil {
			return err
		}
		layerRoot := filepath.Join(dotsRoot, layer)
		rel, err := filepath.Rel(layerRoot, target)
		if err == nil && rel != "" && !strings.HasPrefix(rel, "..") {
			return fmt.Errorf("%s is already managed by yo", filePath)
		}
	}

	if !filepath.IsAbs(filePath) || !strings.HasPrefix(filePath, home+string(filepath.Separator)) {
		return fmt.Errorf("file must be under $HOME")
	}

	rel, err := filepath.Rel(home, filePath)
	if err != nil {
		return err
	}

	dst := filepath.Join(dotsRoot, layer, rel)

	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("%s already exists in the %s layer", rel, layer)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	if err := os.Rename(filePath, dst); err != nil {
		return err
	}

	return os.Symlink(dst, filePath)
}
