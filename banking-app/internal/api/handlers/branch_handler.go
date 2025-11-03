package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shubhbham/BankingApi_Golang/internal/core"
)

type BranchHandler struct {
	service *core.BranchService
}

func NewBranchHandler(service *core.BranchService) *BranchHandler {
	return &BranchHandler{service: service}
}

func (h *BranchHandler) GetBranch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid branch ID")
		return
	}
	branch, err := h.service.GetBranch(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Branch not found")
		return
	}
	respondJSON(w, http.StatusOK, branch)
}
func (h *BranchHandler) ListBranches(w http.ResponseWriter, r *http.Request) {
	branches, err := h.service.ListBranches(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, branches)
}
func (h *BranchHandler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var branch core.Branch
	if err := json.NewDecoder(r.Body).Decode(&branch); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.CreateBranch(r.Context(), &branch); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, branch)
}
