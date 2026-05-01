package main

import (
	"fmt"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/repo"
	"github.com/simskij/yo/internal/ui"
	"github.com/spf13/cobra"
)

var initYoPathFlag string
var initDryRun bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize yo configuration and repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := runInit(initDryRun)
		return err
	},
}

// runInit initialises the yo config and git repo. In dry-run mode it prints
// what would happen without touching the filesystem. It returns the resolved
// yo.path so callers (e.g. dotsInitCmd) can use it without re-loading config.
func runInit(dryRun bool) (string, error) {
	exists, err := config.Exists()
	if err != nil {
		return "", err
	}
	if exists {
		fmt.Println(ui.Dim.Render("yo config already exists, skipping init."))
		cfg, err := config.Load()
		if err != nil {
			return "", err
		}
		return cfg.Yo.Path, nil
	}

	yoPath := initYoPathFlag
	if yoPath == "" {
		if dryRun {
			yoPath = "~/.yofiles"
		} else {
			yoPath, err = ui.Prompt(
				"Where should yo store its files?",
				"~/.yofiles",
				"~/.yofiles",
			)
			if err != nil {
				return "", err
			}
		}
	}

	if dryRun {
		fmt.Printf("%s Would write config to ~/.config/yo/config.yaml (dry run)\n", ui.Dim.Render("·"))
		expanded, err := config.ExpandPath(yoPath)
		if err != nil {
			return "", err
		}
		fmt.Printf("%s Would initialise repository at %s (dry run)\n", ui.Dim.Render("·"), expanded)
		return yoPath, nil
	}

	cfg := config.DefaultConfig()
	cfg.Yo.Path = yoPath

	if err := config.Save(cfg); err != nil {
		return "", err
	}

	expanded, err := config.ExpandPath(yoPath)
	if err != nil {
		return "", err
	}

	if err := repo.Init(expanded); err != nil {
		return "", err
	}

	fmt.Println(ui.Green.Render("✓") + " Config written to ~/.config/yo/config.yaml")
	fmt.Println(ui.Green.Render("✓") + " Repository initialised at " + expanded)
	return yoPath, nil
}

func init() {
	initCmd.Flags().StringVar(&initYoPathFlag, "path", "", "path for yo repository (default: ~/.yofiles)")
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "preview what yo init would do without touching the filesystem")
}
