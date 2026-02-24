package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/coderstone/ghist/internal/models"
)

const taskColumns = `id, title, description, plan, status, milestone, commit_hash, priority, type, ref_id, legacy_id, created_at, updated_at`

type CreateTaskInput struct {
	Title, Description, Status, Milestone, Priority, Type, LegacyID string
}

func (s *Store) CreateTask(in CreateTaskInput) (*models.Task, error) {
	if in.Status == "" {
		in.Status = "todo"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.Exec(
		`INSERT INTO tasks (title, description, status, milestone, priority, type, legacy_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		in.Title, in.Description, in.Status, in.Milestone, in.Priority, in.Type, in.LegacyID, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting task: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert id: %w", err)
	}
	// Backfill ref_id
	refID := fmt.Sprintf("GHST-%d", id)
	if _, err := s.db.Exec(`UPDATE tasks SET ref_id = ? WHERE id = ?`, refID, id); err != nil {
		return nil, fmt.Errorf("setting ref_id: %w", err)
	}
	return s.GetTask(id)
}

func (s *Store) GetTask(id int64) (*models.Task, error) {
	row := s.db.QueryRow(`SELECT `+taskColumns+` FROM tasks WHERE id = ?`, id)
	return scanTask(row)
}

func (s *Store) ListTasks(status, milestone, priority, taskType string) ([]models.Task, error) {
	query := `SELECT ` + taskColumns + ` FROM tasks`
	var conditions []string
	var args []any

	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}
	if milestone != "" {
		conditions = append(conditions, "milestone = ?")
		args = append(args, milestone)
	}
	if priority != "" {
		conditions = append(conditions, "priority = ?")
		args = append(args, priority)
	}
	if taskType != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, taskType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY id ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		t, err := scanTaskRow(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

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

func (s *Store) UpdateTask(id int64, u TaskUpdate) (*models.Task, error) {
	var sets []string
	var args []any

	if u.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *u.Title)
	}
	if u.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *u.Description)
	}
	if u.Plan != nil {
		sets = append(sets, "plan = ?")
		args = append(args, *u.Plan)
	}
	if u.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *u.Status)
	}
	if u.Milestone != nil {
		sets = append(sets, "milestone = ?")
		args = append(args, *u.Milestone)
	}
	if u.CommitHash != nil {
		sets = append(sets, "commit_hash = ?")
		args = append(args, *u.CommitHash)
	}
	if u.Priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *u.Priority)
	}
	if u.Type != nil {
		sets = append(sets, "type = ?")
		args = append(args, *u.Type)
	}
	if u.LegacyID != nil {
		sets = append(sets, "legacy_id = ?")
		args = append(args, *u.LegacyID)
	}

	if len(sets) == 0 {
		return s.GetTask(id)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	sets = append(sets, "updated_at = ?")
	args = append(args, now)
	args = append(args, id)

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(sets, ", "))
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("updating task: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return nil, fmt.Errorf("task %d not found", id)
	}
	return s.GetTask(id)
}

func (s *Store) DeleteTask(id int64) error {
	result, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("task %d not found", id)
	}
	return nil
}

func (s *Store) TaskCountsByStatus() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT status, COUNT(*) FROM tasks GROUP BY status`)
	if err != nil {
		return nil, fmt.Errorf("counting tasks by status: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

func (s *Store) MilestoneInfo() ([]models.MilestoneInfo, error) {
	rows, err := s.db.Query(`
		SELECT milestone, COUNT(*) as total,
		       SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END) as done
		FROM tasks
		WHERE milestone != ''
		GROUP BY milestone
		ORDER BY milestone`)
	if err != nil {
		return nil, fmt.Errorf("querying milestones: %w", err)
	}
	defer rows.Close()

	var milestones []models.MilestoneInfo
	for rows.Next() {
		var m models.MilestoneInfo
		if err := rows.Scan(&m.Name, &m.Total, &m.Done); err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}
	return milestones, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTask(row scanner) (*models.Task, error) {
	var t models.Task
	var createdAt, updatedAt string
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Plan, &t.Status, &t.Milestone, &t.CommitHash, &t.Priority, &t.Type, &t.RefID, &t.LegacyID, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("scanning task: %w", err)
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &t, nil
}

func scanTaskRow(rows *sql.Rows) (*models.Task, error) {
	return scanTask(rows)
}
