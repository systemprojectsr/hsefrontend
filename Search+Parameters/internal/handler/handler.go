package handler

import (
	"ServiceApi/internal/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) SearchServices(w http.ResponseWriter, r *http.Request) {
	filters := map[string]string{
		"price_from": r.URL.Query().Get("price_from"),
		"price_to":   r.URL.Query().Get("price_to"),
		"location":   r.URL.Query().Get("location"),
		"rating":     r.URL.Query().Get("rating"),
		"sort":       r.URL.Query().Get("sort"),
	}

	services, err := h.service.SearchServices(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(services)
}
