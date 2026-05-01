package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/dots"
	"github.com/simskij/yo/internal/ui"
	"github.com/spf13/cobra"
)

var machineFlag string

var dotsCmd = &cobra.Command{
	Use:   "dots",
	Short: "Manage dotfiles",
}

var dotsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold dotfiles repo structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runInit(); err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		if err := dots.Init(cfg.Dots.Path, machine); err != nil {
			return err
		}
		expanded, err := config.ExpandPath(cfg.Dots.Path)
		if err != nil {
			return err
		}
		fmt.Printf("%s Scaffolded dotfiles repo at %s\n", ui.Green.Render("✓"), expanded)
		fmt.Printf("  %s global/\n", ui.Dim.Render("·"))
		fmt.Printf("  %s %s/\n", ui.Dim.Render("·"), machine)
		return nil
	},
}

var dotsApplyCmd = &cobra.Command{
	Use:   "apply [subpath]",
	Short: "Symlink dotfiles into $HOME",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		subpath := ""
		if len(args) > 0 {
			subpath = args[0]
		}
		entries, err := dots.Collect(cfg.Dots.Path, machine)
		if err != nil {
			return err
		}
		home, _ := os.UserHomeDir()
		applied := 0
		for _, e := range entries {
			if subpath != "" {
				if !isSubpath(e.Dst, filepath.Join(home, subpath)) {
					continue
				}
			}
			if e.Status == dots.StatusLinked {
				continue
			}
			if err := dots.ApplyEntry(e); err != nil {
				return fmt.Errorf("apply %s: %w", e.Dst, err)
			}
			fmt.Printf("%s %s\n", ui.Green.Render("✓"), e.Dst)
			applied++
		}
		if applied == 0 {
			fmt.Println(ui.Dim.Render("Nothing to apply."))
		}
		return nil
	},
}

var dotsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show dotfile link status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		entries, err := dots.Collect(cfg.Dots.Path, machine)
		if err != nil {
			return err
		}
		for _, e := range entries {
			switch e.Status {
			case dots.StatusLinked:
				fmt.Printf("%s %s\n", ui.Green.Render("✓"), e.Dst)
			case dots.StatusDiffers:
				fmt.Printf("%s %s\n", ui.Yellow.Render("~"), e.Dst)
			case dots.StatusMissing:
				fmt.Printf("%s %s\n", ui.Red.Render("✗"), e.Dst)
			}
		}
		return nil
	},
}

var dotsDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show content diff for files that differ",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		return dots.Diff(cfg.Dots.Path, machine, os.Stdout)
	},
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("yo is not initialized — run `yo init` first")
	}
	return cfg, nil
}

func hostname(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	return os.Hostname()
}

func isSubpath(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel == "." || len(rel) > 0 && rel[0] != '.'
}

var dotsAddCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Move a file into the dotfiles repo and symlink it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		layer := "global"
		if machineFlag != "" {
			layer = machineFlag
		}

		filePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		if err := dots.Add(cfg.Dots.Path, layer, filePath); err != nil {
			return err
		}

		fmt.Printf("%s %s → %s/%s\n", ui.Green.Render("✓"), filePath, layer, args[0])
		return nil
	},
}

var dotsRemoveCmd = &cobra.Command{
	Use:   "remove <file>",
	Short: "Stop managing a file and restore it in place",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		filePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		layers, err := dots.FindLayers(cfg.Dots.Path, filePath)
		if err != nil {
			return err
		}
		if len(layers) == 0 {
			return fmt.Errorf("%s is not managed by yo", filePath)
		}

		layer := layers[0]
		if len(layers) > 1 {
			layer, err = ui.Select("Which layer do you want to remove from?", layers)
			if err != nil {
				return err
			}
		}

		if err := dots.Remove(cfg.Dots.Path, filePath, layer); err != nil {
			return err
		}

		fmt.Printf("%s %s restored from %s layer\n", ui.Green.Render("✓"), filePath, layer)
		return nil
	},
}

func init() {
	dotsCmd.PersistentFlags().StringVar(&machineFlag, "machine", "", "override machine name (default: hostname)")
	dotsCmd.AddCommand(dotsInitCmd, dotsApplyCmd, dotsStatusCmd, dotsDiffCmd, dotsAddCmd, dotsRemoveCmd)
}
