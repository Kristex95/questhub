package handler

import (
	"context"
	"net/http"
	"time"
)

type dbPinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	pool dbPinger
}

func NewHealthHandler(pool dbPinger) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) Livez(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("database unavailable"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
