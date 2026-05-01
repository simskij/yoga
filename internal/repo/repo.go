package repo

import (
	"os"
	"os/exec"
	"path/filepath"
)

const readme = `# .yofiles

Managed by yo. Each domain owns a subdirectory:

  dots/     — dotfiles (symlinked into $HOME)

Run ` + "`yo --help`" + ` for available commands.
`

// Init creates the yofiles root directory, initialises a git repo,
// and writes a README if none exists.
func Init(yoPath string) error {
	if err := os.MkdirAll(yoPath, 0o755); err != nil {
		return err
	}

	readmePath := filepath.Join(yoPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil {
			return err
		}
	}

	gitDir := filepath.Join(yoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", yoPath)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
