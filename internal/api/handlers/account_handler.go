package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shubhbham/BankingApi_Golang/internal/core"
)

type AccountHandler struct {
	service *core.AccountService
}

func NewAccountHandler(service *core.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account core.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if account.Status == "" {
		account.Status = "ACTIVE"
	}

	if err := h.service.CreateAccount(r.Context(), &account); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	account, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	respondJSON(w, http.StatusOK, account)
}

func (h *AccountHandler) GetAccountByNumber(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.URL.Query().Get("number")
	if accountNumber == "" {
		respondError(w, http.StatusBadRequest, "Account number is required")
		return
	}

	account, err := h.service.GetAccountByNumber(r.Context(), accountNumber)
	if err != nil {
		respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	respondJSON(w, http.StatusOK, account)
}

func (h *AccountHandler) ListAccountsByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := uuid.Parse(vars["customer_id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	accounts, err := h.service.ListAccountsByCustomer(r.Context(), customerID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, accounts)
}

func (h *AccountHandler) CloseAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	if err := h.service.CloseAccount(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Account closed successfully",
	})
}