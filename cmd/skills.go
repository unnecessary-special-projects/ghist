package cmd

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// skillsFS is set from main.go via SetSkillsFS
var skillsFS embed.FS

func SetSkillsFS(fs embed.FS) {
	skillsFS = fs
}

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage skill definitions",
}

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := skillsFS.ReadDir("skills")
		if err != nil {
			return fmt.Errorf("reading skills: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), ".md")
			content, err := skillsFS.ReadFile(filepath.Join("skills", entry.Name()))
			if err != nil {
				continue
			}

			// Extract title from first line (# Title) and description from first non-empty line after
			lines := strings.Split(string(content), "\n")
			title := name
			description := ""
			if len(lines) > 0 {
				title = strings.TrimPrefix(strings.TrimSpace(lines[0]), "# ")
			}
			for _, line := range lines[1:] {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
					description = trimmed
					break
				}
			}

			fmt.Printf("  %-20s %s\n", name, title)
			if description != "" {
				fmt.Printf("  %-20s %s\n", "", description)
			}
			fmt.Println()
		}

		return nil
	},
}

var skillsShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show a skill's full instructions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		filename := name + ".md"

		content, err := skillsFS.ReadFile(filepath.Join("skills", filename))
		if err != nil {
			return fmt.Errorf("skill %q not found (run 'ghist skills list' to see available skills)", name)
		}

		fmt.Print(string(content))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(skillsCmd)
	skillsCmd.AddCommand(skillsListCmd)
	skillsCmd.AddCommand(skillsShowCmd)
}
