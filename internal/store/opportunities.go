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

func (s *Store) opportunitiesDir() string {
	return filepath.Join(s.root, "opportunities")
}

func (s *Store) opportunityPath(id int64) string {
	return filepath.Join(s.opportunitiesDir(), fmt.Sprintf("%d.json", id))
}

func (s *Store) CreateOpportunity(name, notes string) (*models.Opportunity, error) {
	id, err := nextID(s.opportunitiesDir())
	if err != nil {
		return nil, fmt.Errorf("getting next id: %w", err)
	}
	now := time.Now().UTC()
	o := models.Opportunity{
		ID:        id,
		Name:      name,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.writeOpportunity(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Store) GetOpportunity(id int64) (*models.Opportunity, error) {
	data, err := os.ReadFile(s.opportunityPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("opportunity not found")
		}
		return nil, fmt.Errorf("reading opportunity %d: %w", id, err)
	}
	var o models.Opportunity
	if err := json.Unmarshal(data, &o); err != nil {
		return nil, fmt.Errorf("parsing opportunity %d: %w", id, err)
	}
	return &o, nil
}

func (s *Store) ListOpportunities() ([]models.Opportunity, error) {
	entries, err := os.ReadDir(s.opportunitiesDir())
	if err != nil {
		return nil, fmt.Errorf("listing opportunities: %w", err)
	}
	var opps []models.Opportunity
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.opportunitiesDir(), e.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading opportunity file %s: %w", e.Name(), err)
		}
		var o models.Opportunity
		if err := json.Unmarshal(data, &o); err != nil {
			return nil, fmt.Errorf("parsing opportunity file %s: %w", e.Name(), err)
		}
		opps = append(opps, o)
	}
	sort.Slice(opps, func(i, j int) bool {
		return opps[i].ID < opps[j].ID
	})
	return opps, nil
}

func (s *Store) writeOpportunity(o *models.Opportunity) error {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling opportunity: %w", err)
	}
	return os.WriteFile(s.opportunityPath(o.ID), data, 0644)
}
