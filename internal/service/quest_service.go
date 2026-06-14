package service

import (
	"fmt"
	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/repository"
)

type QuestService struct {
	quests repository.QuestRepository
	tasks  repository.TaskRepository
	users  repository.UserRepository
}

func NewQuestService(quests repository.QuestRepository, tasks repository.TaskRepository, users repository.UserRepository) *QuestService {
	return &QuestService{
		quests: quests,
		tasks:  tasks,
		users:  users,
	}
}

func (service *QuestService) CreateQuest(title, description string, difficulty int) (*domain.Quest, error) {
	incomingQuest, err := domain.NewQuest(
		"",
		title,
		description,
		difficulty,
	)
	if err != nil {
		return nil, fmt.Errorf("create quest domain model: %w", err)
	}

	resultQuest, err := service.quests.Create(*incomingQuest)
	if err != nil {
		return nil, fmt.Errorf("persist quest: %w", err)
	}

	return &resultQuest, nil
}

func (service *QuestService) GetQuest(id string) (*domain.Quest, error) {
	quest, err := service.quests.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get quest: %w", err)
	}
	return &quest, nil
}

func (service *QuestService) ListQuests() ([]*domain.Quest, error) {
	quests, err := service.quests.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all quests: %w", err)
	}

	result := make([]*domain.Quest, 0, len(quests))
	for i := range quests {
		result = append(result, &quests[i])
	}

	return result, nil
}

func (service *QuestService) DeleteQuest(id string) error {
	quest, err := service.quests.Get(id)
	if err != nil {
		return fmt.Errorf("get quest: %w", err)
	}

	tasks, err := service.tasks.GetByQuestID(quest.ID)
	if err != nil {
		return fmt.Errorf("get quest tasks: %w", err)
	}

	for _, t := range tasks {
		if err := service.tasks.Delete(t.ID); err != nil {
			return fmt.Errorf("delete task: %w", err)
		}
	}

	if err := service.quests.Delete(id); err != nil {
		return fmt.Errorf("delete quest: %w", err)
	}

	return nil
}

func (service *QuestService) AddTaskToQuest(questID, title, description string) (*domain.Task, error) {
	if _, err := service.quests.Get(questID); err != nil {
		return nil, fmt.Errorf("get quest: %w", err)
	}

	task := domain.Task{
		ID:          "",
		QuestId:     questID,
		Title:       title,
		Description: description,
	}

	created, err := service.tasks.Create(task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return &created, nil
}

func (service *QuestService) GetQuestTasks(questID string) ([]*domain.Task, error) {
	if _, err := service.quests.Get(questID); err != nil {
		return nil, fmt.Errorf("get quest: %w", err)
	}

	tasks, err := service.tasks.GetByQuestID(questID)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	return tasks, nil
}

func (service *QuestService) CompleteQuest(questID, userID string) error {
	quest, err := service.quests.Get(questID)
	if err != nil {
		return fmt.Errorf("get quest: %w", err)
	}

	user, err := service.users.Get(userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	tasks, err := service.tasks.GetByQuestID(questID)
	if err != nil {
		return fmt.Errorf("get tasks: %w", err)
	}

	for _, task := range tasks {
		if !task.GetIsCompleted() {
			return fmt.Errorf("validation error: quest not completed")
		}
	}

	user.AddXP(quest.Difficulty * 100)

	if _, err := service.users.Update(userID, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}
