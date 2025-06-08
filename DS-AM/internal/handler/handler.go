package handler

import (
	"encoding/json"
	"net/http"
	"ServiceApi/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) SearchServices(w http.ResponseWriter, r *http.Request) {
	// Пример обработки запроса с фильтрами
	filters := map[string]string{
		"price_from": r.URL.Query().Get("price_from"),
		"price_to":   r.URL.Query().Get("price_to"),
		"location":   r.URL.Query().Get("location"),
		"rating":     r.URL.Query().Get("rating"),
	}

	services, err := h.service.SearchServices(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}