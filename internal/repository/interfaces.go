package repository

import "github.com/Kristex95/questhub/internal/domain"

type QuestRepository interface {
	Create(quest domain.Quest) (domain.Quest, error)
	Get(id string) (domain.Quest, error)
	GetAll() ([]domain.Quest, error)
	Update(questID string, updatedQuest domain.Quest) (domain.Quest, error)
	Delete(questID string) error
}

type TaskRepository interface {
	Create(task domain.Task) (domain.Task, error)
	Get(id string) (domain.Task, error)
	GetAll() ([]domain.Task, error)
	GetByQuestID(questID string) ([]*domain.Task, error)
	Update(taskID string, updatedTask domain.Task) (domain.Task, error)
	Delete(taskID string) error
}

type UserRepository interface {
	Create(task domain.User) (domain.User, error)
	Get(id string) (domain.User, error)
	GetAll() ([]domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Update(userID string, updatedUser domain.User) (domain.User, error)
	Delete(userID string) error
}
