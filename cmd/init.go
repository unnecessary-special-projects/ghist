package cmd

import (
	"fmt"
	"os"

	"github.com/coderstone/ghist/internal/project"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ghist in the current directory",
	Long:  "Creates .ghist/ directory, initializes the SQLite database, injects into CLAUDE.md, and optionally adds .ghist/ to .gitignore.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}

		fmt.Println("Initializing ghist...")

		if err := project.Init(cwd, os.Stdin); err != nil {
			return err
		}

		fmt.Println("ghist initialized successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
