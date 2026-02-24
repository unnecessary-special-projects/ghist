package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/coderstone/ghist/internal/models"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		asJSON, _ := cmd.Flags().GetBool("json")

		counts, err := s.TaskCountsByStatus()
		if err != nil {
			return err
		}

		milestones, err := s.MilestoneInfo()
		if err != nil {
			return err
		}

		events, err := s.ListEvents(5)
		if err != nil {
			return err
		}

		total := 0
		for _, c := range counts {
			total += c
		}

		summary := models.StatusSummary{
			TotalTasks:    total,
			TasksByStatus: counts,
			Milestones:    milestones,
			RecentEvents:  events,
		}

		if asJSON {
			data, err := json.MarshalIndent(summary, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Project Status\n")
		fmt.Printf("==============\n\n")

		fmt.Printf("Tasks: %d total", total)
		if total > 0 {
			fmt.Printf(" (")
			first := true
			for _, status := range []string{"todo", "in_progress", "done", "blocked"} {
				if c, ok := counts[status]; ok && c > 0 {
					if !first {
						fmt.Printf(", ")
					}
					fmt.Printf("%d %s", c, status)
					first = false
				}
			}
			fmt.Printf(")")
		}
		fmt.Println()

		if len(milestones) > 0 {
			fmt.Printf("\nMilestones:\n")
			for _, m := range milestones {
				pct := 0
				if m.Total > 0 {
					pct = (m.Done * 100) / m.Total
				}
				fmt.Printf("  %-20s %d/%d (%d%%)\n", m.Name, m.Done, m.Total, pct)
			}
		}

		if len(events) > 0 {
			fmt.Printf("\nRecent Events:\n")
			for _, e := range events {
				taskInfo := ""
				if e.TaskID != nil {
					taskInfo = fmt.Sprintf(" (task #%d)", *e.TaskID)
				}
				fmt.Printf("  [%s] %s%s\n", e.CreatedAt.Format("2006-01-02 15:04"), e.Message, taskInfo)
			}
		}

		return nil
	},
}

func init() {
	statusCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(statusCmd)
}
