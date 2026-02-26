package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Store holds the root .ghist/ directory path.
type Store struct {
	root string
}

// Open initialises a file-based store rooted at ghistDir (the .ghist/ directory).
// It runs SQLite-to-JSON migration if a legacy ghist.sqlite is present, then
// ensures the tasks/, events/, and opportunities/ subdirectories exist.
func Open(ghistDir string) (*Store, error) {
	if err := MigrateSQLiteToJSON(ghistDir); err != nil {
		return nil, fmt.Errorf("migrating sqlite: %w", err)
	}

	for _, dir := range []string{"tasks", "events", "opportunities"} {
		if err := os.MkdirAll(filepath.Join(ghistDir, dir), 0755); err != nil {
			return nil, fmt.Errorf("creating %s directory: %w", dir, err)
		}
	}

	// Ensure settings.json exists
	settingsPath := filepath.Join(ghistDir, "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		if err := os.WriteFile(settingsPath, []byte("{}"), 0644); err != nil {
			return nil, fmt.Errorf("creating settings.json: %w", err)
		}
	}

	return &Store{root: ghistDir}, nil
}

// Close is a no-op for the file-based store; retained for interface compatibility.
func (s *Store) Close() error {
	return nil
}

// nextID returns the next available integer ID for a given subdirectory by
// scanning existing JSON filenames and returning max+1.
func nextID(dir string) (int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("reading directory %s: %w", dir, err)
	}
	var max int64
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		base := strings.TrimSuffix(e.Name(), ".json")
		id, err := strconv.ParseInt(base, 10, 64)
		if err != nil {
			continue
		}
		if id > max {
			max = id
		}
	}
	return max + 1, nil
}
