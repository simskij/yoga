package main

import (
	"fmt"

	"github.com/simskij/yo/internal/ui"
	"github.com/simskij/yo/internal/upgrade"
	"github.com/spf13/cobra"
)

const defaultRepo = "simskij/yoga"

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade yo to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s Fetching latest release from %s...\n", ui.Dim.Render("·"), defaultRepo)
		upgraded, err := upgrade.Run(defaultRepo, version)
		if err != nil {
			return err
		}
		if upgraded {
			fmt.Printf("%s yo upgraded successfully\n", ui.Green.Render("✓"))
		} else {
			fmt.Println(ui.Dim.Render("Already up to date."))
		}
		return nil
	},
}
