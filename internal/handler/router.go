package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(questH *QuestHandler, leaderboardH *LeaderboardHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/quests/{questID}/complete", questH.Complete)
	r.Get("/api/leaderboard", leaderboardH.List)

	return r
}
