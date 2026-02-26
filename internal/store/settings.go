package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type settings struct {
	MilestoneOrder []string `json:"milestone_order"`
}

func (s *Store) settingsPath() string {
	return filepath.Join(s.root, "settings.json")
}

func (s *Store) readSettings() (settings, error) {
	var st settings
	data, err := os.ReadFile(s.settingsPath())
	if err != nil {
		return st, nil // missing file → zero value
	}
	_ = json.Unmarshal(data, &st)
	return st, nil
}

func (s *Store) writeSettings(st settings) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.settingsPath(), data, 0644)
}

// GetMilestoneOrder returns the saved milestone ordering.
func (s *Store) GetMilestoneOrder() ([]string, error) {
	st, err := s.readSettings()
	if err != nil {
		return nil, err
	}
	if st.MilestoneOrder == nil {
		return []string{}, nil
	}
	return st.MilestoneOrder, nil
}

// SetMilestoneOrder saves the milestone ordering.
func (s *Store) SetMilestoneOrder(order []string) error {
	st, err := s.readSettings()
	if err != nil {
		return err
	}
	st.MilestoneOrder = order
	return s.writeSettings(st)
}
