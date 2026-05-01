package main

import (
	"fmt"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/repo"
	"github.com/simskij/yo/internal/ui"
	"github.com/spf13/cobra"
)

var initYoPathFlag string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize yo configuration and repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if exists {
		fmt.Println(ui.Dim.Render("yo config already exists, skipping init."))
		return nil
	}

	yoPath := initYoPathFlag
	if yoPath == "" {
		yoPath, err = ui.Prompt(
			"Where should yo store its files?",
			"~/.yofiles",
			"~/.yofiles",
		)
		if err != nil {
			return err
		}
	}

	cfg := config.DefaultConfig()
	cfg.Yo.Path = yoPath

	if err := config.Save(cfg); err != nil {
		return err
	}

	expanded, err := config.ExpandPath(yoPath)
	if err != nil {
		return err
	}

	if err := repo.Init(expanded); err != nil {
		return err
	}

	fmt.Println(ui.Green.Render("✓") + " Config written to ~/.config/yo/config.yaml")
	fmt.Println(ui.Green.Render("✓") + " Repository initialised at " + expanded)
	return nil
}

func init() {
	initCmd.Flags().StringVar(&initYoPathFlag, "yo-path", "", "path for yo repository (default: ~/.yofiles)")
}
