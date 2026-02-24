package cmd

import (
	"fmt"
	"os"

	"github.com/unnecessary-special-projects/ghist/internal/project"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Update ghist configuration after an upgrade",
	Long:  "Re-runs database migrations, updates CLAUDE.md injection, and refreshes context. Use this after upgrading ghist to apply new changes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}

		fmt.Println("Refreshing ghist...")

		if err := project.Refresh(cwd, os.Stdin); err != nil {
			return err
		}

		fmt.Println("ghist refreshed successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
