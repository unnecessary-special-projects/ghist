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

func (s *Store) eventsDir() string {
	return filepath.Join(s.root, "events")
}

func (s *Store) eventPath(id int64) string {
	return filepath.Join(s.eventsDir(), fmt.Sprintf("%d.json", id))
}

func (s *Store) CreateEvent(typ, message, metadata string, taskID *int64) (*models.Event, error) {
	if typ == "" {
		typ = "log"
	}
	if metadata == "" {
		metadata = "{}"
	}
	id, err := nextID(s.eventsDir())
	if err != nil {
		return nil, fmt.Errorf("getting next id: %w", err)
	}
	e := models.Event{
		ID:        id,
		Type:      typ,
		Message:   message,
		Metadata:  metadata,
		TaskID:    taskID,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.writeEvent(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *Store) GetEvent(id int64) (*models.Event, error) {
	data, err := os.ReadFile(s.eventPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("reading event %d: %w", id, err)
	}
	var e models.Event
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("parsing event %d: %w", id, err)
	}
	return &e, nil
}

func (s *Store) ListEvents(limit int) ([]models.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	events, err := s.readAllEvents()
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}

func (s *Store) ListEventsByTask(taskID int64) ([]models.Event, error) {
	events, err := s.readAllEvents()
	if err != nil {
		return nil, err
	}
	var filtered []models.Event
	for _, e := range events {
		if e.TaskID != nil && *e.TaskID == taskID {
			filtered = append(filtered, e)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})
	return filtered, nil
}

func (s *Store) readAllEvents() ([]models.Event, error) {
	entries, err := os.ReadDir(s.eventsDir())
	if err != nil {
		return nil, fmt.Errorf("listing events: %w", err)
	}
	var events []models.Event
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.eventsDir(), e.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading event file %s: %w", e.Name(), err)
		}
		var ev models.Event
		if err := json.Unmarshal(data, &ev); err != nil {
			return nil, fmt.Errorf("parsing event file %s: %w", e.Name(), err)
		}
		events = append(events, ev)
	}
	return events, nil
}

// clearEventTaskID sets task_id to nil on all events referencing taskID.
// Used as a cascade when a task is deleted.
func (s *Store) clearEventTaskID(taskID int64) {
	events, err := s.readAllEvents()
	if err != nil {
		return
	}
	for i := range events {
		if events[i].TaskID != nil && *events[i].TaskID == taskID {
			events[i].TaskID = nil
			s.writeEvent(&events[i]) //nolint:errcheck
		}
	}
}

func (s *Store) writeEvent(e *models.Event) error {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}
	return os.WriteFile(s.eventPath(e.ID), data, 0644)
}
