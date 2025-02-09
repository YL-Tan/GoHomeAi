package server

import (
	"encoding/json"
	"net/http"

	"github.com/YL-Tan/GoHomeAi/internal/db"
)

type Handler struct {
	Queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{Queries: q}
}

func (h *Handler) getDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := h.Queries.GetDevices(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(devices)
}

func InitRouter(q *db.Queries) *http.ServeMux{
	mux := http.NewServeMux()
	handler := NewHandler(q)
	mux.HandleFunc("/devices", handler.getDevicesHandler)
	return mux
}
