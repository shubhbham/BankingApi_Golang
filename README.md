# Banking API

A RESTful API for banking operations built with Go and PostgreSQL (Supabase).

## Features

- Customer management (CRUD operations)
- Account management (Savings, Current, Term Deposit)
- Transaction processing (Credit, Debit, Transfer)
- Branch management
- RESTful API design
- Database connection pooling
- Graceful shutdown
- Middleware (Logging, Recovery, CORS)

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── router.go            # Route definitions
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── account_handler.go
│   │   │   ├── customer_handler.go
│   │   │   ├── transaction_handler.go
│   │   │   └── health_handler.go
│   │   └── middleware/          # Middleware
│   │       ├── logging.go
│   │       ├── recovery.go
│   │       └── cors.go
│   ├── core/                    # Business logic
│   │   ├── account.go
│   │   ├── customer.go
│   │   ├── transaction.go
│   │   └── errors.go
│   ├── db/                      # Database layer
│   │   └── db.go
│   ├── config/                  # Configuration
│   │   └── config.go
│   └── server/                  # Server setup
│       └── server.go
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Supabase)
- Make (optional, for using Makefile commands)

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd banking-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file:
```bash
cp .env.example .env
```

4. Update `.env` with your Supabase credentials:
```properties
DATABASE_URL=postgresql://user:password@host:port/database
SERVER_PORT=8080
ENVIRONMENT=development
```

5. Run the database migrations on Supabase using the provided schema.

## Running the Application

### Using Make:
```bash
make run
```

### Using Go directly:
```bash
go run cmd/server/main.go
```

### Build and run:
```bash
make build
./bin/server
```

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### Customers
- `POST /api/v1/customers` - Create a new customer
- `GET /api/v1/customers` - List all customers (with pagination)
- `GET /api/v1/customers/{id}` - Get customer by ID
- `PUT /api/v1/customers/{id}` - Update customer

### Accounts
- `POST /api/v1/accounts` - Create a new account
- `GET /api/v1/accounts/{id}` - Get account by ID
- `GET /api/v1/accounts?number={number}` - Get account by account number
- `GET /api/v1/customers/{customer_id}/accounts` - List customer accounts
- `POST /api/v1/accounts/{id}/close` - Close an account

### Transactions
- `POST /api/v1/transactions` - Create a transaction (credit/debit)
- `GET /api/v1/transactions/{id}` - Get transaction by ID
- `GET /api/v1/accounts/{account_id}/transactions` - List account transactions
- `POST /api/v1/transactions/transfer` - Transfer between accounts

## Branch
- `POST /api/v1/branches` - POST branch details
- `GET /api/v1/branches/{id}` - Get branch detail with id
- `GET /api/v1/branches` - Get the branches

## Example API Calls

### Create Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "mobile": "+919876543210",
    "date_of_birth": "1990-01-15",
    "kyc_status": "PENDING"
  }'
```

### Create Account
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "uuid-here",
    "branch_id": "uuid-here",
    "account_type": "SAVINGS",
    "account_number": "1234567890",
    "balance": 1000.00
  }'
```

### Create Transaction
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "uuid-here",
    "txn_type": "CREDIT",
    "amount": 500.00,
    "description": "Deposit",
    "channel": "ONLINE"
  }'
```

### Transfer Money
```bash
curl -X POST http://localhost:8080/api/v1/transactions/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": "uuid-here",
    "to_account_id": "uuid-here",
    "amount": 100.00,
    "description": "Transfer to savings"
  }'
```

## Development

### Running Tests
```bash
make test
```

### Code Formatting
```bash
make fmt
```

### Clean Build Artifacts
```bash
make clean
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DATABASE_URL | PostgreSQL connection string | Required |
| SERVER_PORT | Server port | 8080 |
| ENVIRONMENT | Environment (development/production) | development |

## Features to Implement

- [ ] Authentication & Authorization (JWT)
- [ ] Card management
- [ ] Loan management
- [ ] Rate limiting
- [ ] API documentation (Swagger)
- [ ] Unit tests
- [ ] Integration tests
- [ ] Docker support
- [ ] CI/CD pipeline

## License

MIT
