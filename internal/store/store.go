package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const currentSchemaVersion = 3

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("setting WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	version := 0

	// Check if schema_version table exists
	row := s.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='schema_version'")
	var name string
	if err := row.Scan(&name); err == nil {
		// Table exists, read version
		row = s.db.QueryRow("SELECT version FROM schema_version LIMIT 1")
		if err := row.Scan(&version); err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("reading schema version: %w", err)
		}
	}

	if version >= currentSchemaVersion {
		return nil
	}

	// Run v1+v2 in a transaction
	if version < 2 {
		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("starting migration transaction: %w", err)
		}
		defer tx.Rollback()

		if version < 1 {
			if err := migrateV1(tx); err != nil {
				return fmt.Errorf("migration v1: %w", err)
			}
		}

		if version < 2 {
			if err := migrateV2(tx); err != nil {
				return fmt.Errorf("migration v2: %w", err)
			}
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration: %w", err)
		}

		version = 2
	}

	// Run v3 separately (needs PRAGMA foreign_keys changes outside transaction)
	if version < 3 {
		if err := s.migrateV3(); err != nil {
			return fmt.Errorf("migration v3: %w", err)
		}
	}

	return nil
}

func migrateV1(tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			title       TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status      TEXT NOT NULL DEFAULT 'todo'
			            CHECK (status IN ('todo', 'in_progress', 'done', 'blocked')),
			milestone   TEXT NOT NULL DEFAULT '',
			commit_hash TEXT NOT NULL DEFAULT '',
			created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
			updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
		)`,
		`CREATE TABLE IF NOT EXISTS opportunities (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			notes      TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
			updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			type       TEXT NOT NULL DEFAULT 'log',
			message    TEXT NOT NULL,
			metadata   TEXT NOT NULL DEFAULT '{}',
			task_id    INTEGER,
			created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER NOT NULL
		)`,
		`INSERT INTO schema_version (version) VALUES (1)`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("executing %q: %w", stmt[:40], err)
		}
	}

	return nil
}

func migrateV2(tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE tasks ADD COLUMN plan TEXT NOT NULL DEFAULT ''`,
		`UPDATE schema_version SET version = 2`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("executing %q: %w", stmt[:40], err)
		}
	}

	return nil
}

// migrateV3 runs on *Store because PRAGMA foreign_keys can't change inside a transaction.
func (s *Store) migrateV3() error {
	// Disable foreign keys (must be outside transaction)
	if _, err := s.db.Exec("PRAGMA foreign_keys=OFF"); err != nil {
		return fmt.Errorf("disabling foreign keys: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("starting v3 transaction: %w", err)
	}
	defer tx.Rollback()

	stmts := []string{
		// Create new tasks table with updated schema
		`CREATE TABLE tasks_new (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			title       TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			plan        TEXT NOT NULL DEFAULT '',
			status      TEXT NOT NULL DEFAULT 'todo'
			            CHECK (status IN ('todo', 'in_planning', 'in_progress', 'done', 'blocked')),
			milestone   TEXT NOT NULL DEFAULT '',
			commit_hash TEXT NOT NULL DEFAULT '',
			priority    TEXT NOT NULL DEFAULT '' CHECK (priority IN ('', 'low', 'medium', 'high', 'urgent')),
			type        TEXT NOT NULL DEFAULT '' CHECK (type IN ('', 'bug', 'feature', 'improvement', 'chore')),
			ref_id      TEXT NOT NULL DEFAULT '',
			legacy_id   TEXT NOT NULL DEFAULT '',
			created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
			updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
		)`,
		// Copy data from old table
		`INSERT INTO tasks_new (id, title, description, plan, status, milestone, commit_hash, created_at, updated_at)
		 SELECT id, title, description, plan, status, milestone, commit_hash, created_at, updated_at FROM tasks`,
		// Backfill ref_id
		`UPDATE tasks_new SET ref_id = 'GHST-' || CAST(id AS TEXT)`,
		// Swap tables
		`DROP TABLE tasks`,
		`ALTER TABLE tasks_new RENAME TO tasks`,
		// Update schema version
		`UPDATE schema_version SET version = 3`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("v3 executing %q: %w", stmt[:40], err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing v3: %w", err)
	}

	// Re-enable foreign keys
	if _, err := s.db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return fmt.Errorf("re-enabling foreign keys: %w", err)
	}

	return nil
}
