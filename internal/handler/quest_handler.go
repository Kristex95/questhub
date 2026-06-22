package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/go-chi/chi/v5"
)

type questCompleter interface {
	CompleteQuest(ctx context.Context, userID, questID int64) error
}

type QuestHandler struct {
	quests questCompleter
}

func NewQuestHandler(quests questCompleter) *QuestHandler {
	return &QuestHandler{quests: quests}
}

func (h *QuestHandler) Complete(w http.ResponseWriter, r *http.Request) {
	questID, err := strconv.ParseInt(chi.URLParam(r, "questID"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid questID")
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid or missing user_id")
		return
	}

	if err := h.quests.CompleteQuest(r.Context(), userID, questID); err != nil {
		var notFound *domain.NotFoundError
		if errors.As(err, &notFound) {
			writeError(w, http.StatusNotFound, notFound.Error())
			return
		}
		var valErr *domain.ValidationError
		if errors.As(err, &valErr) {
			writeError(w, http.StatusBadRequest, valErr.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to complete quest")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "completed",
		"quest_id": questID,
		"user_id":  userID,
	})
}
