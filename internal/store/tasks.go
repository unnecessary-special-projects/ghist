package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/unnecessary-special-projects/ghist/internal/models"
)

// CreateTaskInput holds the fields needed to create a new task.
type CreateTaskInput struct {
	Title, Description, Status, Milestone, Priority, Type, LegacyID string
}

// TaskUpdate holds optional fields to update on an existing task.
type TaskUpdate struct {
	Title       *string
	Description *string
	Plan        *string
	Status      *string
	Milestone   *string
	CommitHash  *string
	Priority    *string
	Type        *string
	LegacyID    *string
}

func (s *Store) tasksDir() string {
	return filepath.Join(s.root, "tasks")
}

func (s *Store) taskPath(id int64) string {
	return filepath.Join(s.tasksDir(), fmt.Sprintf("%d.json", id))
}

func (s *Store) CreateTask(in CreateTaskInput) (*models.Task, error) {
	if in.Status == "" {
		in.Status = "todo"
	}
	id, err := nextID(s.tasksDir())
	if err != nil {
		return nil, fmt.Errorf("getting next id: %w", err)
	}
	now := time.Now().UTC()
	t := models.Task{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		Milestone:   in.Milestone,
		Priority:    in.Priority,
		Type:        in.Type,
		LegacyID:    in.LegacyID,
		RefID:       fmt.Sprintf("GHST-%d", id),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.writeTask(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) GetTask(id int64) (*models.Task, error) {
	data, err := os.ReadFile(s.taskPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("reading task %d: %w", id, err)
	}
	var t models.Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing task %d: %w", id, err)
	}
	return &t, nil
}

func (s *Store) ListTasks(status, milestone, priority, taskType string) ([]models.Task, error) {
	entries, err := os.ReadDir(s.tasksDir())
	if err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}

	var tasks []models.Task
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.tasksDir(), e.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading task file %s: %w", e.Name(), err)
		}
		var t models.Task
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("parsing task file %s: %w", e.Name(), err)
		}
		if status != "" && t.Status != status {
			continue
		}
		if milestone != "" && t.Milestone != milestone {
			continue
		}
		if priority != "" && t.Priority != priority {
			continue
		}
		if taskType != "" && t.Type != taskType {
			continue
		}
		tasks = append(tasks, t)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	return tasks, nil
}

func (s *Store) UpdateTask(id int64, u TaskUpdate) (*models.Task, error) {
	t, err := s.GetTask(id)
	if err != nil {
		return nil, fmt.Errorf("task %d not found", id)
	}

	if u.Title != nil {
		t.Title = *u.Title
	}
	if u.Description != nil {
		t.Description = *u.Description
	}
	if u.Plan != nil {
		t.Plan = *u.Plan
	}
	if u.Status != nil {
		t.Status = *u.Status
	}
	if u.Milestone != nil {
		t.Milestone = *u.Milestone
	}
	if u.CommitHash != nil {
		t.CommitHash = *u.CommitHash
	}
	if u.Priority != nil {
		t.Priority = *u.Priority
	}
	if u.Type != nil {
		t.Type = *u.Type
	}
	if u.LegacyID != nil {
		t.LegacyID = *u.LegacyID
	}
	t.UpdatedAt = time.Now().UTC()

	if err := s.writeTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Store) DeleteTask(id int64) error {
	path := s.taskPath(id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("task %d not found", id)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting task %d: %w", id, err)
	}
	// Cascade: clear task_id on any events that reference this task.
	s.clearEventTaskID(id)
	return nil
}

func (s *Store) TaskCountsByStatus() (map[string]int, error) {
	tasks, err := s.ListTasks("", "", "", "")
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, t := range tasks {
		counts[t.Status]++
	}
	return counts, nil
}

func (s *Store) MilestoneInfo() ([]models.MilestoneInfo, error) {
	tasks, err := s.ListTasks("", "", "", "")
	if err != nil {
		return nil, err
	}

	type milestoneData struct {
		total int
		done  int
	}
	mmap := make(map[string]*milestoneData)
	var order []string
	for _, t := range tasks {
		if t.Milestone == "" {
			continue
		}
		if _, ok := mmap[t.Milestone]; !ok {
			mmap[t.Milestone] = &milestoneData{}
			order = append(order, t.Milestone)
		}
		mmap[t.Milestone].total++
		if t.Status == "done" {
			mmap[t.Milestone].done++
		}
	}

	sort.Strings(order)
	var milestones []models.MilestoneInfo
	for _, name := range order {
		m := mmap[name]
		milestones = append(milestones, models.MilestoneInfo{
			Name:  name,
			Total: m.total,
			Done:  m.done,
		})
	}
	return milestones, nil
}

func (s *Store) writeTask(t *models.Task) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling task: %w", err)
	}
	return os.WriteFile(s.taskPath(t.ID), data, 0644)
}
