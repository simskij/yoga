package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"filippo.io/age"
	"github.com/simskij/yo/internal/config"
	"github.com/simskij/yo/internal/crypt"
	"github.com/simskij/yo/internal/dots"
	"github.com/simskij/yo/internal/ui"
	"github.com/spf13/cobra"
)

var machineFlag string
var dotsInitDryRun bool
var applyDryRun bool
var cloneDryRun bool

var dotsCmd = &cobra.Command{
	Use:   "dots",
	Short: "Manage dotfiles",
}

var dotsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold dots directory structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDotsInit(dotsInitDryRun)
	},
}

func runDotsInit(dryRun bool) error {
	yoPath, err := runInit(dryRun)
	if err != nil {
		return err
	}

	machine, err := hostname(machineFlag)
	if err != nil {
		return err
	}

	expanded, err := config.ExpandPath(yoPath)
	if err != nil {
		return err
	}
	dp := filepath.Join(expanded, "dots")

	if dryRun {
		fmt.Printf("%s Would scaffold dots at %s (dry run)\n", ui.Dim.Render("·"), dp)
		fmt.Printf("  %s global/\n", ui.Dim.Render("·"))
		fmt.Printf("  %s %s/\n", ui.Dim.Render("·"), machine)
		return nil
	}

	if err := dots.Init(dp, machine); err != nil {
		return err
	}
	fmt.Printf("%s Scaffolded dots at %s\n", ui.Green.Render("✓"), dp)
	fmt.Printf("  %s global/\n", ui.Dim.Render("·"))
	fmt.Printf("  %s %s/\n", ui.Dim.Render("·"), machine)
	return nil
}

var dotsApplyCmd = &cobra.Command{
	Use:   "apply [subpath]",
	Short: "Symlink dotfiles into $HOME",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subpath := ""
		if len(args) > 0 {
			subpath = args[0]
		}
		return runDotsApply(machineFlag, subpath, applyDryRun)
	},
}

func runDotsApply(machine, subpath string, dryRun bool) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	dp, err := cfg.DotsPath()
	if err != nil {
		return err
	}
	if machine == "" {
		machine, err = hostname("")
		if err != nil {
			return err
		}
	}
	entries, err := dots.Collect(dp, machine)
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
		if e.Encrypted {
			if dryRun {
				fmt.Printf("%s %s %s (dry run)\n", ui.Dim.Render("·"), ui.Dim.Render("🔒"), e.Dst)
				applied++
				continue
			}
			ids, err := loadIdentities(cfg)
			if err != nil {
				return err
			}
			if err := dots.ApplyEncrypted(e.Src, e.Dst, ids); err != nil {
				return fmt.Errorf("apply %s: %w", e.Dst, err)
			}
			fmt.Printf("%s %s %s\n", ui.Green.Render("✓"), ui.Dim.Render("🔒"), e.Dst)
			applied++
			continue
		}
		if e.Status == dots.StatusLinked {
			if dryRun {
				fmt.Printf("%s %s\n", ui.Green.Render("✓"), e.Dst)
			}
			continue
		}
		if dryRun {
			fmt.Printf("%s %s (dry run)\n", ui.Dim.Render("·"), e.Dst)
			applied++
			continue
		}
		if err := dots.ApplyEntry(e); err != nil {
			return fmt.Errorf("apply %s: %w", e.Dst, err)
		}
		fmt.Printf("%s %s\n", ui.Green.Render("✓"), e.Dst)
		applied++
	}
	if applied == 0 && !dryRun {
		fmt.Println(ui.Dim.Render("Nothing to apply."))
	}
	return nil
}

var dotsCloneCmd = &cobra.Command{
	Use:   "clone <repo>",
	Short: "Clone an existing dotfiles repo and apply it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDotsClone(args[0], cloneDryRun)
	},
}

func runDotsClone(repo string, dryRun bool) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	dp, err := cfg.DotsPath()
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("%s Would clone %s into %s (dry run)\n", ui.Dim.Render("·"), repo, dp)
		fmt.Printf("%s Would run dots apply (dry run)\n", ui.Dim.Render("·"))
		return nil
	}

	if entries, err := os.ReadDir(dp); err == nil && len(entries) > 0 {
		return fmt.Errorf("dots path %s already exists and is non-empty; remove it first", dp)
	}

	if err := os.MkdirAll(filepath.Dir(dp), 0o755); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", repo, dp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone: %w", err)
	}

	return runDotsApply(machineFlag, "", false)
}

var dotsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show dotfile link status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		dp, err := cfg.DotsPath()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		entries, err := dots.Collect(dp, machine)
		if err != nil {
			return err
		}
		for _, e := range entries {
			enc := ""
			if e.Encrypted {
				enc = " " + ui.Dim.Render("🔒")
			}
			switch e.Status {
			case dots.StatusLinked:
				fmt.Printf("%s%s %s\n", ui.Green.Render("✓"), enc, e.Dst)
			case dots.StatusDiffers:
				fmt.Printf("%s%s %s\n", ui.Yellow.Render("~"), enc, e.Dst)
			case dots.StatusMissing:
				fmt.Printf("%s%s %s\n", ui.Red.Render("✗"), enc, e.Dst)
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
		dp, err := cfg.DotsPath()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		return dots.Diff(dp, machine, os.Stdout)
	},
}

var encryptFlag bool

var dotsAddCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Move a file into the dotfiles repo and symlink it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		dp, err := cfg.DotsPath()
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
		if encryptFlag {
			ids, err := loadIdentities(cfg)
			if err != nil {
				return err
			}
			if err := dots.AddEncrypted(dp, layer, filePath, ids); err != nil {
				return err
			}
			fmt.Printf("%s %s %s → %s/%s.age\n", ui.Green.Render("✓"), ui.Dim.Render("🔒"), filePath, layer, args[0])
			return nil
		}
		if err := dots.Add(dp, layer, filePath); err != nil {
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
		dp, err := cfg.DotsPath()
		if err != nil {
			return err
		}
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		layers, err := dots.FindLayers(dp, filePath)
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
		if err := dots.Remove(dp, filePath, layer); err != nil {
			return err
		}
		fmt.Printf("%s %s restored from %s layer\n", ui.Green.Render("✓"), filePath, layer)
		return nil
	},
}

var dotsEditCmd = &cobra.Command{
	Use:   "edit <file>",
	Short: "Decrypt, edit, and re-encrypt a managed encrypted file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		dp, err := cfg.DotsPath()
		if err != nil {
			return err
		}
		machine, err := hostname(machineFlag)
		if err != nil {
			return err
		}
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		entries, err := dots.Collect(dp, machine)
		if err != nil {
			return err
		}
		var srcPath string
		for _, e := range entries {
			if e.Dst == filePath && e.Encrypted {
				srcPath = e.Src
				break
			}
		}
		if srcPath == "" {
			return fmt.Errorf("%s is not a managed encrypted file", filePath)
		}
		ids, err := loadIdentities(cfg)
		if err != nil {
			return err
		}
		return dots.EditEncrypted(srcPath, ids)
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

func loadIdentities(cfg *config.Config) ([]age.Identity, error) {
	var extra []string
	if cfg.Dots.Encryption.Identity != "" {
		expanded, err := config.ExpandPath(cfg.Dots.Encryption.Identity)
		if err != nil {
			return nil, err
		}
		extra = append(extra, expanded)
	}
	return crypt.LoadIdentities(extra...)
}

func init() {
	dotsCmd.PersistentFlags().StringVar(&machineFlag, "machine", "", "override machine name (default: hostname)")
	dotsAddCmd.Flags().BoolVar(&encryptFlag, "encrypt", false, "encrypt the file before adding to the repo")
	dotsInitCmd.Flags().BoolVar(&dotsInitDryRun, "dry-run", false, "preview what dots init would do without touching the filesystem")
	dotsApplyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "preview what would be applied without touching the filesystem")
	dotsCloneCmd.Flags().BoolVar(&cloneDryRun, "dry-run", false, "preview what clone would do without touching the filesystem")
	dotsCmd.AddCommand(dotsInitCmd, dotsApplyCmd, dotsStatusCmd, dotsDiffCmd, dotsAddCmd, dotsRemoveCmd, dotsEditCmd, dotsCloneCmd)
}
