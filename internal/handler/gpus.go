package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemPotapov52/gpurenta-backend/internal/db"
	"github.com/ArtemPotapov52/gpurenta-backend/internal/types"
)

type GPUHandler struct {
	Store *db.Store
}

func (h *GPUHandler) List(w http.ResponseWriter, r *http.Request) {
	minVRAM, _ := strconv.Atoi(r.URL.Query().Get("min_vram"))
	imageFilter := r.URL.Query().Get("image")

	agents, err := h.Store.ListOnlineGPUs(r.Context(), minVRAM, imageFilter)
	if err != nil {
		http.Error(w, `{"error":"failed to list GPUs"}`, http.StatusInternalServerError)
		return
	}
	if agents == nil {
		agents = []types.Agent{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}
