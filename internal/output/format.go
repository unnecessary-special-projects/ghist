package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/coderstone/ghist/internal/models"
)

func PrintTaskTable(tasks []models.Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "REF\tTITLE\tSTATUS\tPRIORITY\tTYPE\tMILESTONE\tPLAN")
	fmt.Fprintln(w, "---\t-----\t------\t--------\t----\t---------\t----")
	for _, t := range tasks {
		plan := ""
		if t.Plan != "" {
			plan = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", t.RefID, t.Title, StatusLabel(t.Status), t.Priority, t.Type, t.Milestone, plan)
	}
	w.Flush()
}

func PrintTaskDetail(t *models.Task, events []models.Event) {
	fmt.Printf("Task %s\n", t.RefID)
	fmt.Printf("  Title:       %s\n", t.Title)
	fmt.Printf("  Status:      %s\n", StatusLabel(t.Status))
	if t.Priority != "" {
		fmt.Printf("  Priority:    %s\n", t.Priority)
	}
	if t.Type != "" {
		fmt.Printf("  Type:        %s\n", t.Type)
	}
	if t.Description != "" {
		fmt.Printf("  Description: %s\n", t.Description)
	}
	if t.Milestone != "" {
		fmt.Printf("  Milestone:   %s\n", t.Milestone)
	}
	if t.Plan != "" {
		fmt.Printf("  Plan:\n")
		for _, line := range strings.Split(t.Plan, "\n") {
			fmt.Printf("    %s\n", line)
		}
	}
	if t.CommitHash != "" {
		fmt.Printf("  Commit:      %s\n", t.CommitHash)
	}
	if t.LegacyID != "" {
		fmt.Printf("  Legacy ID:   %s\n", t.LegacyID)
	}
	fmt.Printf("  Created:     %s\n", t.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("  Updated:     %s\n", t.UpdatedAt.Format("2006-01-02 15:04"))

	if len(events) > 0 {
		fmt.Println()
		fmt.Println("  Events:")
		for _, e := range events {
			fmt.Printf("    [%s] %s (%s)\n", e.CreatedAt.Format("2006-01-02 15:04"), e.Message, e.Type)
		}
	}
}

func StatusLabel(status string) string {
	switch status {
	case "todo":
		return "todo"
	case "in_planning":
		return "in_planning"
	case "in_progress":
		return "in_progress"
	case "done":
		return "done"
	case "blocked":
		return "blocked"
	default:
		return status
	}
}
