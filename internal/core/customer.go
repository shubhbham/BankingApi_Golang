package core

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Customer struct {
	CustomerID  uuid.UUID  `json:"customer_id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Mobile      string     `json:"mobile"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	KYCStatus   string     `json:"kyc_status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CustomerService struct {
	db *pgxpool.Pool
}

func NewCustomerService(db *pgxpool.Pool) *CustomerService {
	return &CustomerService{db: db}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, c *Customer) error {
	query := `
		INSERT INTO customers (name, email, mobile, date_of_birth, kyc_status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING customer_id, created_at, updated_at`

	err := s.db.QueryRow(ctx, query, c.Name, c.Email, c.Mobile, c.DateOfBirth, c.KYCStatus).
		Scan(&c.CustomerID, &c.CreatedAt, &c.UpdatedAt)

	return err
}

func (s *CustomerService) GetCustomer(ctx context.Context, id uuid.UUID) (*Customer, error) {
	query := `
		SELECT customer_id, name, email, mobile, date_of_birth, kyc_status, created_at, updated_at
		FROM customers
		WHERE customer_id = $1`

	c := &Customer{}
	err := s.db.QueryRow(ctx, query, id).Scan(
		&c.CustomerID, &c.Name, &c.Email, &c.Mobile, &c.DateOfBirth,
		&c.KYCStatus, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, c *Customer) error {
	query := `
		UPDATE customers
		SET name = $1, email = $2, mobile = $3, date_of_birth = $4, 
		    kyc_status = $5, updated_at = now()
		WHERE customer_id = $6
		RETURNING updated_at`

	err := s.db.QueryRow(ctx, query, c.Name, c.Email, c.Mobile, c.DateOfBirth,
		c.KYCStatus, c.CustomerID).Scan(&c.UpdatedAt)

	return err
}

func (s *CustomerService) ListCustomers(ctx context.Context, limit, offset int) ([]*Customer, error) {
	query := `
		SELECT customer_id, name, email, mobile, date_of_birth, kyc_status, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		c := &Customer{}
		err := rows.Scan(
			&c.CustomerID, &c.Name, &c.Email, &c.Mobile, &c.DateOfBirth,
			&c.KYCStatus, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, rows.Err()
}