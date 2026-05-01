package dots

import (
	"os"
	"path/filepath"
)

type FileStatus int

const (
	StatusLinked    FileStatus = iota // correctly symlinked
	StatusDiffers                     // exists but points elsewhere or is a regular file
	StatusMissing                     // target does not exist
	StatusEncrypted                   // encrypted file, managed by copy not symlink
)

type Entry struct {
	Src       string
	Dst       string
	Status    FileStatus
	Encrypted bool
}

func Collect(dotsPath, machine string) ([]Entry, error) {
	dotsPath, err := expandPath(dotsPath)
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// global wins first, then machine overrides
	layers := []string{
		filepath.Join(dotsPath, "global"),
		filepath.Join(dotsPath, machine),
	}

	seen := map[string]string{} // dst -> src

	for _, layer := range layers {
		if _, err := os.Stat(layer); os.IsNotExist(err) {
			continue
		}
		err := filepath.WalkDir(layer, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(layer, path)
			if err != nil {
				return err
			}
			// encrypted files map to dst without the .age extension
			dst := filepath.Join(home, rel)
			if isEncrypted(path) {
				dst = dst[:len(dst)-len(encExt)]
			}
			seen[dst] = path
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	entries := make([]Entry, 0, len(seen))
	for dst, src := range seen {
		enc := isEncrypted(src)
		entries = append(entries, Entry{
			Src:       src,
			Dst:       dst,
			Encrypted: enc,
			Status:    fileStatus(src, dst),
		})
	}
	return entries, nil
}

func ApplyEntry(e Entry) error {
	if err := os.MkdirAll(filepath.Dir(e.Dst), 0o755); err != nil {
		return err
	}

	if _, err := os.Lstat(e.Dst); err == nil {
		if err := os.Rename(e.Dst, e.Dst+".bak"); err != nil {
			return err
		}
	}

	return os.Symlink(e.Src, e.Dst)
}

func fileStatus(src, dst string) FileStatus {
	info, err := os.Lstat(dst)
	if err != nil {
		return StatusMissing
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return StatusDiffers
	}
	target, err := os.Readlink(dst)
	if err != nil || target != src {
		return StatusDiffers
	}
	return StatusLinked
}

func expandPath(path string) (string, error) {
	if len(path) > 1 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
