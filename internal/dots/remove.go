package dots

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FindLayers returns which layers (e.g. "global", "mymachine") contain a copy
// of the file at filePath, based on what exists in dotsRoot.
func FindLayers(dotsRoot, filePath string) ([]string, error) {
	dotsRoot, err := expandPath(dotsRoot)
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(home, filePath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dotsRoot)
	if err != nil {
		return nil, err
	}

	var layers []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		candidate := filepath.Join(dotsRoot, e.Name(), rel)
		if _, err := os.Stat(candidate); err == nil {
			layers = append(layers, e.Name())
		}
	}
	return layers, nil
}

// Remove stops managing filePath: removes the repo copy from the given layer
// and restores the file contents to the original location.
func Remove(dotsRoot, filePath, layer string) error {
	dotsRoot, err := expandPath(dotsRoot)
	if err != nil {
		return err
	}

	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", filePath, err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s is not managed by yo", filePath)
	}

	target, err := os.Readlink(filePath)
	if err != nil {
		return err
	}

	layerRoot := filepath.Join(dotsRoot, layer)
	rel, err := filepath.Rel(layerRoot, target)
	if err != nil || rel == "" || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("%s is not managed by yo in layer %s", filePath, layer)
	}

	repoFile := filepath.Join(layerRoot, rel)

	src, err := os.Open(repoFile)
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.Remove(filePath); err != nil {
		return err
	}

	dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	src.Close()

	return os.Remove(repoFile)
}
