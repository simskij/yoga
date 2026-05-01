package dots

import (
	"os"
	"path/filepath"
)

const readme = `# dots

Managed by yo. Structure:

  global/       — files applied to all machines (mini $HOME)
  <hostname>/   — machine-specific overrides (mini $HOME)

Files are symlinked into $HOME mirroring their path within each layer.
Machine-specific files take precedence over global ones.

Run ` + "`yo dots apply`" + ` to apply.
`

func Init(dotsPath, machine string) error {
	dotsPath, err := expandPath(dotsPath)
	if err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(dotsPath, "global"),
		filepath.Join(dotsPath, machine),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	readmePath := filepath.Join(dotsPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil {
			return err
		}
	}

	return nil
}
