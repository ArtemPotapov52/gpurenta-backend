package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemPotapov52/gpurenta/internal/db"
	"github.com/ArtemPotapov52/gpurenta/internal/middleware"
	"github.com/ArtemPotapov52/gpurenta/internal/types"
)

type AgentHandler struct {
	Store *db.Store
}

type registerRequest struct {
	GPUModel        string   `json:"gpu_model"`
	VRAMGB          int      `json:"vram_gb"`
	OS              string   `json:"os"`
	SupportedImages []string `json:"supported_images"`
	PricePerHour    int      `json:"price_per_hour"`
}

type registerResponse struct {
	AgentID string `json:"agent_id"`
	Secret  string `json:"secret"`
}

type heartbeatRequest struct {
	FRPURL string `json:"frp_url"`
}

func (h *AgentHandler) Register(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		middleware.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.GPUModel == "" || req.VRAMGB == 0 {
		middleware.JSONError(w, "gpu_model and vram_gb are required", http.StatusBadRequest)
		return
	}
	if req.PricePerHour == 0 {
		req.PricePerHour = 20
	}

	agent, err := h.Store.CreateAgent(r.Context(), userID, req.GPUModel, req.OS, req.VRAMGB, req.SupportedImages, req.PricePerHour)
	if err != nil {
		middleware.JSONError(w, "failed to register agent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerResponse{
		AgentID: agent.ID,
		Secret:  agent.Secret,
	})
}

func (h *AgentHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	agentID := r.Header.Get("X-Agent-ID")
	agentSecret := r.Header.Get("X-Agent-Secret")
	if agentID == "" || agentSecret == "" {
		middleware.JSONError(w, "missing agent credentials", http.StatusUnauthorized)
		return
	}

	_, err := h.Store.GetAgentBySecret(r.Context(), agentID, agentSecret)
	if err != nil {
		middleware.JSONError(w, "invalid agent credentials", http.StatusUnauthorized)
		return
	}

	var req heartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Store.Heartbeat(r.Context(), agentID, req.FRPURL); err != nil {
		middleware.JSONError(w, "failed to update heartbeat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *AgentHandler) GetSupportedImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types.SupportedImages)
}
