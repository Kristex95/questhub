package repository

import (
	"context"

	"github.com/Kristex95/questhub/internal/models"
)

type QuestRepository interface {
	Create(ctx context.Context, quest *models.Quest) (*models.Quest, error)
	GetByID(ctx context.Context, id int64) (*models.Quest, error)
	GetAll(ctx context.Context) ([]*models.Quest, error)
	Update(ctx context.Context, quest *models.Quest) error
	Delete(ctx context.Context, id int64) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	GetAll(ctx context.Context) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id int64) error
	GetByQuestID(ctx context.Context, questID int64) ([]*models.Task, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	AddXP(ctx context.Context, userID int64, amount int) error
}

type RewardRepository interface {
	Create(ctx context.Context, reward *models.Reward) (*models.Reward, error)
	GetByQuestID(ctx context.Context, questID int64) ([]*models.Reward, error)
}

type ProgressRepository interface {
	Create(ctx context.Context, progress *models.Progress) (*models.Progress, error)
	GetByUserAndQuest(ctx context.Context, userID, questID int64) (*models.Progress, error)
	MarkCompleted(ctx context.Context, userID, questID int64) error
}
