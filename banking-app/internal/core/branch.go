package core

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Branch struct {
	BranchID   uuid.UUID `json:"branch_id"`
	BranchCode string    `json:"branch_code"`
	Name       string    `json:"name"`
	Address    *string   `json:"address,omitempty"`
	City       *string   `json:"city,omitempty"`
	State      *string   `json:"state,omitempty"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
}

type BranchService struct {
	db *pgxpool.Pool
}

func NewBranchService(db *pgxpool.Pool) *BranchService {
	return &BranchService{db: db}
}

func (s *BranchService) CreateBranch(ctx context.Context, b *Branch) error {
	query := `
		INSERT INTO branches (branch_code, name, address, city, state, country)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING branch_id, created_at`

	err := s.db.QueryRow(ctx, query, b.BranchCode, b.Name, b.Address,
		b.City, b.State, b.Country).Scan(&b.BranchID, &b.CreatedAt)

	return err
}

func (s *BranchService) GetBranch(ctx context.Context, id uuid.UUID) (*Branch, error) {
	query := `
		SELECT branch_id, branch_code, name, address, city, state, country, created_at
		FROM branches
		WHERE branch_id = $1`

	b := &Branch{}
	err := s.db.QueryRow(ctx, query, id).Scan(
		&b.BranchID, &b.BranchCode, &b.Name, &b.Address,
		&b.City, &b.State, &b.Country, &b.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *BranchService) ListBranches(ctx context.Context) ([]*Branch, error) {
	query := `
		SELECT branch_id, branch_code, name, address, city, state, country, created_at
		FROM branches
		ORDER BY name`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []*Branch
	for rows.Next() {
		b := &Branch{}
		err := rows.Scan(
			&b.BranchID, &b.BranchCode, &b.Name, &b.Address,
			&b.City, &b.State, &b.Country, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}

	return branches, rows.Err()
}