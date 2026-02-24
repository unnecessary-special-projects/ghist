package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "ghist",
	Short:   "Project memory layer for AI agents",
	Long:    "Ghist maintains persistent project state so AI agents never lose context between sessions.",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
