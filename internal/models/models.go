package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTaskID accepts a bare numeric ID ("19") or a prefixed ref ("GHST-19")
// and returns the numeric task ID.
func ParseTaskID(raw string) (int64, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.ToUpper(s), "GHST-")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid task id: %s", raw)
	}
	return id, nil
}

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Plan        string    `json:"plan"`
	Status      string    `json:"status"`
	Milestone   string    `json:"milestone"`
	CommitHash  string    `json:"commit_hash"`
	Priority    string    `json:"priority"`
	Type        string    `json:"type"`
	RefID       string    `json:"ref_id"`
	LegacyID    string    `json:"legacy_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Event struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Metadata  string    `json:"metadata"`
	TaskID    *int64    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Opportunity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectContext struct {
	Tasks      []Task  `json:"tasks"`
	RecentEvents []Event `json:"recent_events"`
	Summary    StatusSummary `json:"summary"`
}

type StatusSummary struct {
	TotalTasks     int              `json:"total_tasks"`
	TasksByStatus  map[string]int   `json:"tasks_by_status"`
	Milestones     []MilestoneInfo  `json:"milestones"`
	RecentEvents   []Event          `json:"recent_events"`
}

type MilestoneInfo struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
	Done  int    `json:"done"`
}
