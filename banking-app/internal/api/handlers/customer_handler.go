package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shubhbham/BankingApi_Golang/internal/core"
)

type CustomerHandler struct {
	service *core.CustomerService
}

func NewCustomerHandler(service *core.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer core.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if customer.KYCStatus == "" {
		customer.KYCStatus = "PENDING"
	}

	if err := h.service.CreateCustomer(r.Context(), &customer); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.GetCustomer(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	respondJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var customer core.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer.CustomerID = id

	if err := h.service.UpdateCustomer(r.Context(), &customer); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 10
	}

	customers, err := h.service.ListCustomers(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, customers)
}
