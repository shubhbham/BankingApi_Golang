package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shubhbham/BankingApi_Golang/internal/core"
)

type TransactionHandler struct {
	service *core.TransactionService
}

func NewTransactionHandler(service *core.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var txn core.AccountTransaction
	if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.CreateTransaction(r.Context(), &txn); err != nil {
		if err == core.ErrInsufficientFunds {
			respondError(w, http.StatusBadRequest, "Insufficient funds")
			return
		}
		if err == core.ErrAccountClosed {
			respondError(w, http.StatusBadRequest, "Account is closed")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	txn, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	respondJSON(w, http.StatusOK, txn)
}

func (h *TransactionHandler) ListTransactionsByAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := uuid.Parse(vars["account_id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 20
	}

	transactions, err := h.service.ListTransactionsByAccount(r.Context(), accountID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, transactions)
}

type TransferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "Amount must be positive")
		return
	}

	err := h.service.Transfer(r.Context(), req.FromAccountID, req.ToAccountID, req.Amount, req.Description)
	if err != nil {
		if err == core.ErrInsufficientFunds {
			respondError(w, http.StatusBadRequest, "Insufficient funds")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Transfer completed successfully",
	})
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}