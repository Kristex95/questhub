package service

import (
	"context"
	"fmt"

	"github.com/Kristex95/questhub/internal/models"
	"github.com/Kristex95/questhub/internal/repository"
)

var validStatuses = map[string]bool{
	"in_progress": true,
	"completed":   true,
	"abandoned":   true,
}

type ProgressService struct {
	progress repository.ProgressRepository
}

func NewProgressService(progress repository.ProgressRepository) *ProgressService {
	return &ProgressService{progress: progress}
}

func (s *ProgressService) StartProgress(ctx context.Context, userID, questID int64) (*models.Progress, error) {
	p := &models.Progress{
		UserID:         userID,
		QuestID:        questID,
		Status:         "in_progress",
		CompletedTasks: 0,
	}

	created, err := s.progress.Create(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("start progress: %w", err)
	}

	return created, nil
}

func (s *ProgressService) MarkCompleted(ctx context.Context, userID, questID int64) error {
	if err := s.progress.MarkCompleted(ctx, userID, questID); err != nil {
		return fmt.Errorf("mark completed: %w", err)
	}
	return nil
}

func (s *ProgressService) GetProgress(ctx context.Context, userID, questID int64) (*models.Progress, error) {
	p, err := s.progress.GetByUserAndQuest(ctx, userID, questID)
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}
	return p, nil
}
