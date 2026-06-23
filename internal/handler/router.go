package handler

import (
	"log/slog"

	"github.com/Kristex95/questhub/internal/logging"
	"github.com/Kristex95/questhub/internal/observability"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(
	serviceName string,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	questH *QuestHandler,
	leaderboardH *LeaderboardHandler,
) *chi.Mux {
	r := chi.NewRouter()


	r.Use(logging.RequestIDMiddleware)
	r.Use(observability.TracingMiddleware(serviceName))
	r.Use(logging.AccessLog(logger))
	r.Use(middleware.Recoverer)

	health := NewHealthHandler(pool)
	r.Get("/livez", health.Livez)
	r.Get("/readyz", health.Readyz)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/quests", questH.List)
		r.Get("/quests/{questID:[0-9]+}", questH.Get)
		r.Post("/quests", questH.Create)
		r.Get("/leaderboard", leaderboardH.List)
	})

	return r
}
