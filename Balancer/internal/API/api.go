package API

import (
	"cloud/Balancer/internal/models"
	"cloud/Balancer/internal/service"
	"context"
	"encoding/json"
	"net/http"
)

type RateLimitHandler struct {
	rl *service.RateLimiter
}

func NewRateLimitHandler(rl *service.RateLimiter) *RateLimitHandler {
	return &RateLimitHandler{rl: rl}
}

func (h *RateLimitHandler) SetLimit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID string  `json:"client_id"`
		Capacity int     `json:"capacity"`
		Rate     float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.rl.SetLimit(models.ClientIdentifier(req.ClientID), req.Capacity, req.Rate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RateLimitHandler) GetLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if storage, ok := h.rl.Storage.(interface {
		GetAllClientLimits(ctx context.Context) ([]models.ClientLimit, error)
	}); ok {
		limits, err := storage.GetAllClientLimits(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(limits)
		return
	}

	http.Error(w, "Storage doesn't support getting all limits", http.StatusNotImplemented)
}

//curl -X POST -H "Content-Type: application/json" \
//-d '{"client_id":"test_client", "capacity":5, "rate":1.0}' \
//http://localhost:8080/api/set

//curl -s http://localhost:8080/api/get
