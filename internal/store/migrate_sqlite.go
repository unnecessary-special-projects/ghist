package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/unnecessary-special-projects/ghist/internal/models"
	_ "modernc.org/sqlite"
)

// MigrateSQLiteToJSON migrates data from a legacy ghist.sqlite to individual
// JSON files. It is idempotent: if ghist.sqlite does not exist it returns nil.
// After a successful migration ghist.sqlite is renamed to ghist.sqlite.bak.
func MigrateSQLiteToJSON(ghistDir string) error {
	sqlitePath := filepath.Join(ghistDir, "ghist.sqlite")
	if _, err := os.Stat(sqlitePath); os.IsNotExist(err) {
		return nil
	}

	// Ensure subdirectories exist before writing.
	for _, dir := range []string{"tasks", "events", "opportunities"} {
		if err := os.MkdirAll(filepath.Join(ghistDir, dir), 0755); err != nil {
			return fmt.Errorf("creating %s directory: %w", dir, err)
		}
	}

	db, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return fmt.Errorf("opening sqlite for migration: %w", err)
	}
	defer db.Close()

	if err := migrateTasksFromSQLite(db, ghistDir); err != nil {
		return fmt.Errorf("migrating tasks: %w", err)
	}
	if err := migrateEventsFromSQLite(db, ghistDir); err != nil {
		return fmt.Errorf("migrating events: %w", err)
	}
	if err := migrateOpportunitiesFromSQLite(db, ghistDir); err != nil {
		return fmt.Errorf("migrating opportunities: %w", err)
	}

	db.Close()

	bakPath := sqlitePath + ".bak"
	if err := os.Rename(sqlitePath, bakPath); err != nil {
		return fmt.Errorf("renaming sqlite to .bak: %w", err)
	}

	fmt.Println("  Migrated ghist.sqlite → JSON files (backup kept at ghist.sqlite.bak)")
	return nil
}

func migrateTasksFromSQLite(db *sql.DB, ghistDir string) error {
	rows, err := db.Query(`SELECT id, title, description, plan, status, milestone, commit_hash, priority, type, ref_id, legacy_id, created_at, updated_at FROM tasks`)
	if err != nil {
		if isSQLiteNoSuchTable(err) {
			return nil
		}
		return fmt.Errorf("querying tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		var createdAt, updatedAt string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Plan, &t.Status, &t.Milestone, &t.CommitHash, &t.Priority, &t.Type, &t.RefID, &t.LegacyID, &createdAt, &updatedAt); err != nil {
			return fmt.Errorf("scanning task row: %w", err)
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		path := filepath.Join(ghistDir, "tasks", fmt.Sprintf("%d.json", t.ID))
		if _, err := os.Stat(path); err == nil {
			continue // already written
		}
		data, _ := json.MarshalIndent(t, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("writing task %d: %w", t.ID, err)
		}
	}
	return rows.Err()
}

func migrateEventsFromSQLite(db *sql.DB, ghistDir string) error {
	rows, err := db.Query(`SELECT id, type, message, metadata, task_id, created_at FROM events`)
	if err != nil {
		if isSQLiteNoSuchTable(err) {
			return nil
		}
		return fmt.Errorf("querying events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Event
		var createdAt string
		if err := rows.Scan(&e.ID, &e.Type, &e.Message, &e.Metadata, &e.TaskID, &createdAt); err != nil {
			return fmt.Errorf("scanning event row: %w", err)
		}
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		path := filepath.Join(ghistDir, "events", fmt.Sprintf("%d.json", e.ID))
		if _, err := os.Stat(path); err == nil {
			continue // already written
		}
		data, _ := json.MarshalIndent(e, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("writing event %d: %w", e.ID, err)
		}
	}
	return rows.Err()
}

func migrateOpportunitiesFromSQLite(db *sql.DB, ghistDir string) error {
	rows, err := db.Query(`SELECT id, name, notes, created_at, updated_at FROM opportunities`)
	if err != nil {
		if isSQLiteNoSuchTable(err) {
			return nil
		}
		return fmt.Errorf("querying opportunities: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Opportunity
		var createdAt, updatedAt string
		if err := rows.Scan(&o.ID, &o.Name, &o.Notes, &createdAt, &updatedAt); err != nil {
			return fmt.Errorf("scanning opportunity row: %w", err)
		}
		o.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		o.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		path := filepath.Join(ghistDir, "opportunities", fmt.Sprintf("%d.json", o.ID))
		if _, err := os.Stat(path); err == nil {
			continue // already written
		}
		data, _ := json.MarshalIndent(o, "", "  ")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("writing opportunity %d: %w", o.ID, err)
		}
	}
	return rows.Err()
}

func isSQLiteNoSuchTable(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no such table")
}
