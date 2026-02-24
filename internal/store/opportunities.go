package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/coderstone/ghist/internal/models"
)

func (s *Store) CreateOpportunity(name, notes string) (*models.Opportunity, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.Exec(
		`INSERT INTO opportunities (name, notes, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		name, notes, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting opportunity: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert id: %w", err)
	}
	return s.GetOpportunity(id)
}

func (s *Store) GetOpportunity(id int64) (*models.Opportunity, error) {
	row := s.db.QueryRow(`SELECT id, name, notes, created_at, updated_at FROM opportunities WHERE id = ?`, id)
	return scanOpportunity(row)
}

func (s *Store) ListOpportunities() ([]models.Opportunity, error) {
	rows, err := s.db.Query(`SELECT id, name, notes, created_at, updated_at FROM opportunities ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("listing opportunities: %w", err)
	}
	defer rows.Close()

	var opps []models.Opportunity
	for rows.Next() {
		o, err := scanOpportunity(rows)
		if err != nil {
			return nil, err
		}
		opps = append(opps, *o)
	}
	return opps, rows.Err()
}

func scanOpportunity(row scanner) (*models.Opportunity, error) {
	var o models.Opportunity
	var createdAt, updatedAt string
	err := row.Scan(&o.ID, &o.Name, &o.Notes, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("opportunity not found")
		}
		return nil, fmt.Errorf("scanning opportunity: %w", err)
	}
	o.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	o.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &o, nil
}
