package project

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/unnecessary-special-projects/ghist/internal/models"
	"github.com/unnecessary-special-projects/ghist/internal/store"
)

// UpdateContext reads current state from the store and writes current_context.json.
func UpdateContext(root string, s *store.Store) error {
	tasks, err := s.ListTasks("", "", "", "")
	if err != nil {
		return fmt.Errorf("listing tasks: %w", err)
	}

	events, err := s.ListEvents(10)
	if err != nil {
		return fmt.Errorf("listing events: %w", err)
	}

	counts, err := s.TaskCountsByStatus()
	if err != nil {
		return fmt.Errorf("counting tasks: %w", err)
	}

	milestones, err := s.MilestoneInfo()
	if err != nil {
		return fmt.Errorf("querying milestones: %w", err)
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	ctx := models.ProjectContext{
		Tasks:        tasks,
		RecentEvents: events,
		Summary: models.StatusSummary{
			TotalTasks:    total,
			TasksByStatus: counts,
			Milestones:    milestones,
			RecentEvents:  events,
		},
	}

	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling context: %w", err)
	}

	path := ContextPath(root)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing context file: %w", err)
	}

	return nil
}
