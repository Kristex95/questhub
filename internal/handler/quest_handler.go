package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/logging"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/go-chi/chi/v5"
)

type questCompleter interface {
	CompleteQuest(ctx context.Context, userID, questID int64) error
	GetQuest(ctx context.Context, id int64) (*models.Quest, error)
	CreateQuest(ctx context.Context, title, description string, difficulty int) (*models.Quest, error)
	ListQuests(ctx context.Context) ([]*models.Quest, error)
}

type QuestHandler struct {
	quests questCompleter
	logger *slog.Logger
}

func NewQuestHandler(quests questCompleter, logger *slog.Logger) *QuestHandler {
	return &QuestHandler{quests: quests, logger: logger}
}

func (h *QuestHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.LoggerFrom(ctx, h.logger)

	quests, err := h.quests.ListQuests(ctx)
	if err != nil {
		log.ErrorContext(ctx, "failed to list quests", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to list quests")
		return
	}

	writeJSON(w, http.StatusOK, quests)
}

func (h *QuestHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.LoggerFrom(ctx, h.logger)

	questIDStr := chi.URLParam(r, "questID")
	questID, err := strconv.ParseInt(questIDStr, 10, 64)
	if err != nil {
		log.WarnContext(ctx, "invalid quest ID parameter", slog.String("quest_id_raw", questIDStr), slog.String("error", err.Error()))
		writeError(w, http.StatusBadRequest, "invalid questID")
		return
	}

	log.InfoContext(ctx, "fetching quest", slog.Int64("quest_id", questID))

	quest, err := h.quests.GetQuest(ctx, questID)
	if err != nil {
		var notFound *domain.NotFoundError
		if errors.As(err, &notFound) {
			log.WarnContext(ctx, "quest not found", slog.Int64("quest_id", questID))
			writeError(w, http.StatusNotFound, notFound.Error())
			return
		}

		log.ErrorContext(ctx, "failed to fetch quest from service", slog.Int64("quest_id", questID), slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to fetch quest")
		return
	}

	writeJSON(w, http.StatusOK, quest)
}

func (h *QuestHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.LoggerFrom(ctx, h.logger)

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Difficulty  int    `json:"difficulty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WarnContext(ctx, "failed to decode create quest request body", slog.String("error", err.Error()))
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.InfoContext(ctx, "creating quest", slog.String("title", req.Title), slog.Int("difficulty", req.Difficulty))

	quest, err := h.quests.CreateQuest(ctx, req.Title, req.Description, req.Difficulty)
	if err != nil {
		var valErr *domain.ValidationError
		if errors.As(err, &valErr) {
			log.WarnContext(ctx, "quest validation failed", slog.String("title", req.Title), slog.String("error", valErr.Error()))
			writeError(w, http.StatusBadRequest, valErr.Error())
			return
		}

		log.ErrorContext(ctx, "failed to create quest in service", slog.String("title", req.Title), slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to create quest")
		return
	}

	log.InfoContext(ctx, "quest created successfully", slog.Int64("quest_id", quest.ID))
	writeJSON(w, http.StatusCreated, quest)
}

func (h *QuestHandler) Complete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.LoggerFrom(ctx, h.logger)

	questIDStr := chi.URLParam(r, "questID")
	questID, err := strconv.ParseInt(questIDStr, 10, 64)
	if err != nil {
		log.WarnContext(ctx, "invalid quest ID parameter on completion", slog.String("quest_id_raw", questIDStr))
		writeError(w, http.StatusBadRequest, "invalid questID")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.WarnContext(ctx, "invalid or missing user ID parameter on completion", slog.String("user_id_raw", userIDStr), slog.Int64("quest_id", questID))
		writeError(w, http.StatusBadRequest, "invalid or missing user_id")
		return
	}

	log.InfoContext(ctx, "attempting quest completion", slog.Int64("quest_id", questID), slog.Int64("user_id", userID))

	if err := h.quests.CompleteQuest(ctx, userID, questID); err != nil {
		var notFound *domain.NotFoundError
		if errors.As(err, &notFound) {
			log.WarnContext(ctx, "completion failed: resource not found", slog.Int64("quest_id", questID), slog.Int64("user_id", userID), slog.String("error", notFound.Error()))
			writeError(w, http.StatusNotFound, notFound.Error())
			return
		}
		var valErr *domain.ValidationError
		if errors.As(err, &valErr) {
			log.WarnContext(ctx, "completion failed: requirements not met", slog.Int64("quest_id", questID), slog.Int64("user_id", userID), slog.String("error", valErr.Error()))
			writeError(w, http.StatusBadRequest, valErr.Error())
			return
		}

		log.ErrorContext(ctx, "critical system failure during quest completion", slog.Int64("quest_id", questID), slog.Int64("user_id", userID), slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to complete quest")
		return
	}

	log.InfoContext(ctx, "quest completed successfully", slog.Int64("quest_id", questID), slog.Int64("user_id", userID))
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "completed",
		"quest_id": questID,
		"user_id":  userID,
	})
}