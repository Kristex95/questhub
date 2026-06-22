package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/Kristex95/questhub/internal/repository"
)

type statsIncrementer interface {
	IncrCompletedQuests(ctx context.Context) error
}

type QuestService struct {
	quests   repository.QuestRepository
	tasks    repository.TaskRepository
	rewards  *RewardService
	progress *ProgressService
	stats    statsIncrementer
	notifier Notifier
}

func NewQuestService(
	quests repository.QuestRepository,
	tasks repository.TaskRepository,
	rewards *RewardService,
	progress *ProgressService,
	stats statsIncrementer,
	notifier Notifier,
) *QuestService {
	return &QuestService{
		quests:   quests,
		tasks:    tasks,
		rewards:  rewards,
		progress: progress,
		stats:    stats,
		notifier: notifier,
	}
}

func (s *QuestService) CreateQuest(ctx context.Context, title, description string, difficulty int) (*models.Quest, error) {
	if len(title) < 3 {
		return nil, fmt.Errorf("create quest: %w", &domain.ValidationError{
			Field: "title", Message: "must be at least 3 characters long",
		})
	}
	if description == "" {
		return nil, fmt.Errorf("create quest: %w", &domain.ValidationError{
			Field: "description", Message: "must not be empty",
		})
	}
	if difficulty < 1 || difficulty > 10 {
		return nil, fmt.Errorf("create quest: %w", &domain.ValidationError{
			Field: "difficulty", Message: "must be between 1 and 10",
		})
	}

	quest := &models.Quest{
		Title:       title,
		Description: description,
		Difficulty:  difficulty,
		XPReward:    difficulty * 100,
		IsActive:    true,
	}

	created, err := s.quests.Create(ctx, quest)
	if err != nil {
		return nil, fmt.Errorf("create quest: %w", err)
	}

	return created, nil
}

func (s *QuestService) GetQuest(ctx context.Context, id int64) (*models.Quest, error) {
	quest, err := s.quests.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get quest: %w", err)
	}
	return quest, nil
}

func (s *QuestService) ListQuests(ctx context.Context) ([]*models.Quest, error) {
	quests, err := s.quests.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list quests: %w", err)
	}
	return quests, nil
}

func (s *QuestService) DeleteQuest(ctx context.Context, id int64) error {
	if _, err := s.quests.GetByID(ctx, id); err != nil {
		return fmt.Errorf("delete quest: %w", err)
	}

	tasks, err := s.tasks.GetByQuestID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete quest: %w", err)
	}

	for _, task := range tasks {
		if err := s.tasks.Delete(ctx, task.ID); err != nil {
			return fmt.Errorf("delete quest: %w", err)
		}
	}

	if err := s.quests.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete quest: %w", err)
	}

	return nil
}

func (s *QuestService) AddTaskToQuest(ctx context.Context, questID int64, title, description string) (*models.Task, error) {
	if _, err := s.quests.GetByID(ctx, questID); err != nil {
		return nil, fmt.Errorf("add task to quest: %w", err)
	}
	if title == "" {
		return nil, fmt.Errorf("add task to quest: %w", &domain.ValidationError{
			Field: "title", Message: "must not be empty",
		})
	}

	task := &models.Task{
		QuestID:     questID,
		Title:       title,
		Description: description,
	}

	created, err := s.tasks.Create(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("add task to quest: %w", err)
	}

	return created, nil
}

func (s *QuestService) GetQuestTasks(ctx context.Context, questID int64) ([]*models.Task, error) {
	if _, err := s.quests.GetByID(ctx, questID); err != nil {
		return nil, fmt.Errorf("get quest tasks: %w", err)
	}
	tasks, err := s.tasks.GetByQuestID(ctx, questID)
	if err != nil {
		return nil, fmt.Errorf("get quest tasks: %w", err)
	}
	return tasks, nil
}

func (s *QuestService) CompleteQuest(ctx context.Context, userID, questID int64) error {
	quest, err := s.quests.GetByID(ctx, questID)
	if err != nil {
		return fmt.Errorf("complete quest: %w", err)
	}

	tasks, err := s.tasks.GetByQuestID(ctx, questID)
	if err != nil {
		return fmt.Errorf("complete quest: %w", err)
	}

	for _, task := range tasks {
		if !task.IsCompleted {
			return fmt.Errorf("complete quest: %w", &domain.ValidationError{
				Field:   "tasks",
				Message: "not all quest tasks are completed",
			})
		}
	}

	if _, err := s.rewards.users.GetByID(ctx, userID); err != nil {
		return fmt.Errorf("complete quest: %w", err)
	}

	var wg sync.WaitGroup
	errs := make([]error, 4)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := s.rewards.GrantQuestRewards(ctx, userID, questID); err != nil {
			errs[0] = fmt.Errorf("grant rewards: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.progress.MarkCompleted(ctx, userID, questID); err != nil {
			errs[1] = fmt.Errorf("mark progress: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		msg := fmt.Sprintf("Quest #%d completed", questID)
		if err := s.notifier.Notify(ctx, userID, msg); err != nil {
			errs[2] = fmt.Errorf("notify: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.stats.IncrCompletedQuests(ctx); err != nil {
			errs[3] = fmt.Errorf("incr stats: %w", err)
		}
	}()

	wg.Wait()

	if quest.XPReward > 0 {
		if err := s.rewards.users.AddXP(ctx, userID, quest.XPReward); err != nil {
			return fmt.Errorf("complete quest: %w", err)
		}
		user, err := s.rewards.users.GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("complete quest: %w", err)
		}
		if err := s.rewards.stats.UpdateLeaderboard(ctx, userID, user.XP); err != nil {
			return fmt.Errorf("complete quest: %w", err)
		}
	}

	return errors.Join(errs...)
}
