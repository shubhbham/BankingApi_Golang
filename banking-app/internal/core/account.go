package core

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Account struct {
	AccountID     uuid.UUID  `json:"account_id"`
	CustomerID    uuid.UUID  `json:"customer_id"`
	BranchID      uuid.UUID  `json:"branch_id"`
	AccountType   string     `json:"account_type"`
	AccountNumber string     `json:"account_number"`
	Balance       float64    `json:"balance"`
	Status        string     `json:"status"`
	OpenedAt      time.Time  `json:"opened_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
}

type AccountService struct {
	db *pgxpool.Pool
}

func NewAccountService(db *pgxpool.Pool) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(ctx context.Context, a *Account) error {
	query := `
		INSERT INTO accounts (customer_id, branch_id, account_type, account_number, balance, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING account_id, opened_at`

	err := s.db.QueryRow(ctx, query, a.CustomerID, a.BranchID, a.AccountType,
		a.AccountNumber, a.Balance, a.Status).Scan(&a.AccountID, &a.OpenedAt)

	return err
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (*Account, error) {
	query := `
		SELECT account_id, customer_id, branch_id, account_type, account_number, 
		       balance, status, opened_at, closed_at
		FROM accounts
		WHERE account_id = $1`

	a := &Account{}
	err := s.db.QueryRow(ctx, query, id).Scan(
		&a.AccountID, &a.CustomerID, &a.BranchID, &a.AccountType,
		&a.AccountNumber, &a.Balance, &a.Status, &a.OpenedAt, &a.ClosedAt,
	)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AccountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*Account, error) {
	query := `
		SELECT account_id, customer_id, branch_id, account_type, account_number, 
		       balance, status, opened_at, closed_at
		FROM accounts
		WHERE account_number = $1`

	a := &Account{}
	err := s.db.QueryRow(ctx, query, accountNumber).Scan(
		&a.AccountID, &a.CustomerID, &a.BranchID, &a.AccountType,
		&a.AccountNumber, &a.Balance, &a.Status, &a.OpenedAt, &a.ClosedAt,
	)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *AccountService) ListAccountsByCustomer(ctx context.Context, customerID uuid.UUID) ([]*Account, error) {
	query := `
		SELECT account_id, customer_id, branch_id, account_type, account_number, 
		       balance, status, opened_at, closed_at
		FROM accounts
		WHERE customer_id = $1
		ORDER BY opened_at DESC`

	rows, err := s.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		a := &Account{}
		err := rows.Scan(
			&a.AccountID, &a.CustomerID, &a.BranchID, &a.AccountType,
			&a.AccountNumber, &a.Balance, &a.Status, &a.OpenedAt, &a.ClosedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

func (s *AccountService) UpdateBalance(ctx context.Context, accountID uuid.UUID, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE account_id = $2 AND status = 'ACTIVE'`

	result, err := s.db.Exec(ctx, query, amount, accountID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAccountClosed
	}

	return nil
}

func (s *AccountService) CloseAccount(ctx context.Context, accountID uuid.UUID) error {
	query := `
		UPDATE accounts
		SET status = 'CLOSED', closed_at = now()
		WHERE account_id = $1 AND status = 'ACTIVE'`

	result, err := s.db.Exec(ctx, query, accountID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}