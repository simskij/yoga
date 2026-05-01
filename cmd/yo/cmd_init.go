package main

import (
	"fmt"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize yo configuration",
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

	dotsPath, err := ui.Prompt(
		"Where is your dotfiles repo?",
		"~/.dotfiles",
		"~/.dotfiles",
	)
	if err != nil {
		return err
	}

	cfg := config.DefaultConfig()
	cfg.Dots.Path = dotsPath

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Println(ui.Green.Render("✓") + " Config written to ~/.config/yo/config.yaml")
	return nil
}
