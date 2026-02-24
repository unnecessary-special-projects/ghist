package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/coderstone/ghist/internal/models"
)

func (s *Store) CreateEvent(typ, message, metadata string, taskID *int64) (*models.Event, error) {
	if typ == "" {
		typ = "log"
	}
	if metadata == "" {
		metadata = "{}"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.Exec(
		`INSERT INTO events (type, message, metadata, task_id, created_at) VALUES (?, ?, ?, ?, ?)`,
		typ, message, metadata, taskID, now,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting event: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert id: %w", err)
	}
	return s.GetEvent(id)
}

func (s *Store) GetEvent(id int64) (*models.Event, error) {
	row := s.db.QueryRow(`SELECT id, type, message, metadata, task_id, created_at FROM events WHERE id = ?`, id)
	return scanEvent(row)
}

func (s *Store) ListEvents(limit int) ([]models.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(`SELECT id, type, message, metadata, task_id, created_at FROM events ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("listing events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *e)
	}
	return events, rows.Err()
}

func (s *Store) ListEventsByTask(taskID int64) ([]models.Event, error) {
	rows, err := s.db.Query(`SELECT id, type, message, metadata, task_id, created_at FROM events WHERE task_id = ? ORDER BY created_at DESC`, taskID)
	if err != nil {
		return nil, fmt.Errorf("listing events for task %d: %w", taskID, err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *e)
	}
	return events, rows.Err()
}

func scanEvent(row scanner) (*models.Event, error) {
	var e models.Event
	var createdAt string
	err := row.Scan(&e.ID, &e.Type, &e.Message, &e.Metadata, &e.TaskID, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("scanning event: %w", err)
	}
	e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &e, nil
}
