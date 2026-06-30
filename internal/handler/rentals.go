package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemPotapov52/gpurenta-backend/internal/db"
	"github.com/ArtemPotapov52/gpurenta-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type RentalHandler struct {
	Store *db.Store
}

type startRentalRequest struct {
	AgentID string `json:"agent_id"`
	Image   string `json:"image"`
	Hours   int    `json:"hours"`
}

func (h *RentalHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		middleware.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req startRentalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.AgentID == "" || req.Image == "" || req.Hours <= 0 {
		middleware.JSONError(w, "agent_id, image, and hours are required", http.StatusBadRequest)
		return
	}

	agent, err := h.Store.GetAgentByID(r.Context(), req.AgentID)
	if err != nil {
		middleware.JSONError(w, "agent not found", http.StatusNotFound)
		return
	}
	if agent.Status != "online" {
		middleware.JSONError(w, "agent is not online", http.StatusConflict)
		return
	}

	activeRental, _ := h.Store.GetActiveRentalByAgentID(r.Context(), req.AgentID)
	if activeRental != nil {
		middleware.JSONError(w, "GPU is already rented", http.StatusConflict)
		return
	}

	rental, err := h.Store.CreateRental(r.Context(), req.AgentID, userID, req.Image, agent.FRPURL, req.Hours)
	if err != nil {
		middleware.JSONError(w, "failed to create rental", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rental)
}

func (h *RentalHandler) Stop(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		middleware.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rentalID := chi.URLParam(r, "id")
	if rentalID == "" {
		middleware.JSONError(w, "rental id is required", http.StatusBadRequest)
		return
	}

	rental, err := h.Store.StopRental(r.Context(), rentalID)
	if err != nil {
		middleware.JSONError(w, "failed to stop rental", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rental)
}

func (h *RentalHandler) Get(w http.ResponseWriter, r *http.Request) {
	rentalID := chi.URLParam(r, "id")
	if rentalID == "" {
		middleware.JSONError(w, "rental id is required", http.StatusBadRequest)
		return
	}

	rental, err := h.Store.GetRentalByID(r.Context(), rentalID)
	if err != nil {
		middleware.JSONError(w, "rental not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rental)
}

func (h *RentalHandler) ListByAgent(w http.ResponseWriter, r *http.Request) {
	agentID := chi.URLParam(r, "id")
	if agentID == "" {
		middleware.JSONError(w, "agent id is required", http.StatusBadRequest)
		return
	}

	rentals, err := h.Store.ListRentalsByAgentID(r.Context(), agentID)
	if err != nil {
		middleware.JSONError(w, "failed to list rentals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rentals)
}

func (h *RentalHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		middleware.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	status := r.URL.Query().Get("status")
	rentals, err := h.Store.ListRentalsByUserID(r.Context(), userID, status)
	if err != nil {
		middleware.JSONError(w, "failed to list rentals", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rentals)
}

func (h *RentalHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	agentID := r.URL.Query().Get("agent_id")
	if token == "" || agentID == "" {
		middleware.JSONError(w, "token and agent_id required", http.StatusBadRequest)
		return
	}

	rental, err := h.Store.ValidateRentalToken(r.Context(), agentID, token)
	if err != nil {
		middleware.JSONError(w, "invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rental)
}
