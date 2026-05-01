package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yo",
	Short: "Personal kitchen-sink CLI",
}

func main() {
	rootCmd.AddCommand(initCmd, dotsCmd, upgradeCmd, versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
