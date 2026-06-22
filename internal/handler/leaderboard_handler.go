package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Kristex95/questhub/internal/repository/stats"
)

type topProvider interface {
	TopPlayers(ctx context.Context, limit int) ([]stats.LeaderboardEntry, error)
}

type LeaderboardHandler struct {
	stats topProvider
}

func NewLeaderboardHandler(stats topProvider) *LeaderboardHandler {
	return &LeaderboardHandler{stats: stats}
}

func (h *LeaderboardHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if q := r.URL.Query().Get("limit"); q != "" {
		if parsed, err := strconv.Atoi(q); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	entries, err := h.stats.TopPlayers(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load leaderboard")
		return
	}

	writeJSON(w, http.StatusOK, entries)
}
