package api

import (
	"github.com/gorilla/mux"
	"github.com/shubhbham/BankingApi_Golang/internal/api/handlers"
	"github.com/shubhbham/BankingApi_Golang/internal/api/middleware"
	"github.com/shubhbham/BankingApi_Golang/internal/core"
	"github.com/shubhbham/BankingApi_Golang/internal/db"
)

func NewRouter(database *db.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize services
	customerService := core.NewCustomerService(database.Pool)
	accountService := core.NewAccountService(database.Pool)
	transactionService := core.NewTransactionService(database.Pool)
	branchService := core.NewBranchService(database.Pool)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(database)
	customerHandler := handlers.NewCustomerHandler(customerService)
	accountHandler := handlers.NewAccountHandler(accountService)
	branchHandler := handlers.NewBranchHandler(branchService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Apply global middleware
	router.Use(middleware.Recovery)
	router.Use(middleware.Logging)
	router.Use(middleware.CORS)

	// Health check
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Customer routes
	api.HandleFunc("/customers", customerHandler.CreateCustomer).Methods("POST")
	api.HandleFunc("/customers", customerHandler.ListCustomers).Methods("GET")
	api.HandleFunc("/customers/{id}", customerHandler.GetCustomer).Methods("GET")
	api.HandleFunc("/customers/{id}", customerHandler.UpdateCustomer).Methods("PUT")

	// Account routes
	api.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	api.HandleFunc("/accounts/{id}", accountHandler.GetAccount).Methods("GET")
	api.HandleFunc("/accounts", accountHandler.GetAccountByNumber).Methods("GET").Queries("number", "{number}")
	api.HandleFunc("/accounts/{id}/close", accountHandler.CloseAccount).Methods("POST")
	api.HandleFunc("/customers/{customer_id}/accounts", accountHandler.ListAccountsByCustomer).Methods("GET")

	// Transaction routes
	api.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")
	api.HandleFunc("/transactions/{id}", transactionHandler.GetTransaction).Methods("GET")
	api.HandleFunc("/transactions/transfer", transactionHandler.Transfer).Methods("POST")
	api.HandleFunc("/accounts/{account_id}/transactions", transactionHandler.ListTransactionsByAccount).Methods("GET")

	// Branch routes
	api.HandleFunc("/branches", branchHandler.CreateBranch).Methods("POST")
	api.HandleFunc("/branches/{id}", branchHandler.GetBranch).Methods("GET")
	api.HandleFunc("/branches", branchHandler.ListBranches).Methods("GET")

	return router
}
