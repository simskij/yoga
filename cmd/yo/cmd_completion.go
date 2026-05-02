package main

import (
	"os"

	"github.com/simskij/yo/internal/completion"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion <shell>",
	Short: "Generate shell completion scripts",
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenZshCompletion(os.Stdout)
	},
}

var completionNushellCmd = &cobra.Command{
	Use:   "nushell",
	Short: "Generate nushell completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return completion.GenerateNushell(rootCmd, os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(completionBashCmd, completionZshCmd, completionNushellCmd)
}
