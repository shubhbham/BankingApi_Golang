package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountTransaction struct {
	TxnID       uuid.UUID `json:"txn_id"`
	AccountID   uuid.UUID `json:"account_id"`
	TxnType     string    `json:"txn_type"`
	Amount      float64   `json:"amount"`
	Description *string   `json:"description,omitempty"`
	Channel     *string   `json:"channel,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransactionService struct {
	db *pgxpool.Pool
}

func NewTransactionService(db *pgxpool.Pool) *TransactionService {
	return &TransactionService{db: db}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, txn *AccountTransaction) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check account status and balance
	var balance float64
	var status string
	err = tx.QueryRow(ctx, `
		SELECT balance, status FROM accounts WHERE account_id = $1 FOR UPDATE`,
		txn.AccountID,
	).Scan(&balance, &status)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	if status != "ACTIVE" {
		return ErrAccountClosed
	}

	// For DEBIT, check sufficient funds
	if txn.TxnType == "DEBIT" && balance < txn.Amount {
		return ErrInsufficientFunds
	}

	// Insert transaction
	query := `
		INSERT INTO account_transactions (account_id, txn_type, amount, description, channel)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING txn_id, created_at`

	err = tx.QueryRow(ctx, query, txn.AccountID, txn.TxnType, txn.Amount,
		txn.Description, txn.Channel).Scan(&txn.TxnID, &txn.CreatedAt)

	if err != nil {
		return err
	}

	// Update account balance
	balanceChange := txn.Amount
	if txn.TxnType == "DEBIT" {
		balanceChange = -txn.Amount
	}

	_, err = tx.Exec(ctx, `
		UPDATE accounts SET balance = balance + $1 WHERE account_id = $2`,
		balanceChange, txn.AccountID,
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *TransactionService) GetTransaction(ctx context.Context, id uuid.UUID) (*AccountTransaction, error) {
	query := `
		SELECT txn_id, account_id, txn_type, amount, description, channel, created_at
		FROM account_transactions
		WHERE txn_id = $1`

	txn := &AccountTransaction{}
	err := s.db.QueryRow(ctx, query, id).Scan(
		&txn.TxnID, &txn.AccountID, &txn.TxnType, &txn.Amount,
		&txn.Description, &txn.Channel, &txn.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return txn, nil
}

func (s *TransactionService) ListTransactionsByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*AccountTransaction, error) {
	query := `
		SELECT txn_id, account_id, txn_type, amount, description, channel, created_at
		FROM account_transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.Query(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*AccountTransaction
	for rows.Next() {
		txn := &AccountTransaction{}
		err := rows.Scan(
			&txn.TxnID, &txn.AccountID, &txn.TxnType, &txn.Amount,
			&txn.Description, &txn.Channel, &txn.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}

	return transactions, rows.Err()
}

func (s *TransactionService) Transfer(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount float64, description string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock and check source account
	var fromBalance float64
	var fromStatus string
	err = tx.QueryRow(ctx, `
		SELECT balance, status FROM accounts WHERE account_id = $1 FOR UPDATE`,
		fromAccountID,
	).Scan(&fromBalance, &fromStatus)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("source account: %w", ErrNotFound)
		}
		return err
	}

	if fromStatus != "ACTIVE" {
		return fmt.Errorf("source account: %w", ErrAccountClosed)
	}

	if fromBalance < amount {
		return ErrInsufficientFunds
	}

	// Lock and check destination account
	var toStatus string
	err = tx.QueryRow(ctx, `
		SELECT status FROM accounts WHERE account_id = $1 FOR UPDATE`,
		toAccountID,
	).Scan(&toStatus)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("destination account: %w", ErrNotFound)
		}
		return err
	}

	if toStatus != "ACTIVE" {
		return fmt.Errorf("destination account: %w", ErrAccountClosed)
	}

	// Create debit transaction
	_, err = tx.Exec(ctx, `
		INSERT INTO account_transactions (account_id, txn_type, amount, description, channel)
		VALUES ($1, 'DEBIT', $2, $3, 'TRANSFER')`,
		fromAccountID, amount, description,
	)
	if err != nil {
		return err
	}

	// Create credit transaction
	_, err = tx.Exec(ctx, `
		INSERT INTO account_transactions (account_id, txn_type, amount, description, channel)
		VALUES ($1, 'CREDIT', $2, $3, 'TRANSFER')`,
		toAccountID, amount, description,
	)
	if err != nil {
		return err
	}

	// Update balances
	_, err = tx.Exec(ctx, `
		UPDATE accounts SET balance = balance - $1 WHERE account_id = $2`,
		amount, fromAccountID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE accounts SET balance = balance + $1 WHERE account_id = $2`,
		amount, toAccountID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}