package service

import (
	"fmt"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/repository"
)

type GameService struct {
	quests       repository.QuestRepository
	tasks        repository.TaskRepository
	users        repository.UserRepository
	activeQuests map[int]int // userID -> questID
}

func NewGameService(quests repository.QuestRepository, tasks repository.TaskRepository, users repository.UserRepository) *GameService {
	return &GameService{
		quests:       quests,
		tasks:        tasks,
		users:        users,
		activeQuests: make(map[int]int),
	}
}

func (service *GameService) StartQuest(userID int, questID int) error {
	if _, err := service.users.Get(userID); err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if _, err := service.quests.Get(questID); err != nil {
		return fmt.Errorf("get quest: %w", err)
	}

	if _, exists := service.activeQuests[userID]; exists {
		return fmt.Errorf("%w: user already has an active quest", domain.ErrAlreadyStarted)
	}

	service.activeQuests[userID] = questID
	return nil
}

func (service *GameService) CompleteTask(userID, taskID int) error {
	activeQuestID, exists := service.activeQuests[userID]

	if !exists {
		return fmt.Errorf("%w: user has no active quest", domain.ErrInvalidState)
	}

	task, err := service.tasks.Get(taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if task.QuestId != activeQuestID {
		return fmt.Errorf("%w: task does not belong to the user's active quest", domain.ErrValidation)
	}

	quest, err := service.quests.Get(activeQuestID)
	if err != nil {
		return fmt.Errorf("get quest: %w", err)
	}

	tasks, err := service.tasks.GetByQuestID(activeQuestID)
	if err != nil {
		return fmt.Errorf("get quest tasks: %w", err)
	}
	quest.Tasks = tasks

	if err := quest.CompleteTask(taskID); err != nil {
		return fmt.Errorf("domain complete task: %w", err)
	}

	var updatedTask domain.Task
	for _, t := range quest.Tasks {
		if t.ID == taskID {
			updatedTask = *t
			break
		}
	}

	if _, err := service.tasks.Update(taskID, updatedTask); err != nil {
		return fmt.Errorf("update task repository: %w", err)
	}

	return nil
}

func (service *GameService) GetProgress(userID int) (*domain.Progress, error) {
	activeQuestID, hasActive := service.activeQuests[userID]

	if !hasActive {
		return nil, fmt.Errorf("%w: user has no active quest", domain.ErrNotFound)
	}

	quest, err := service.quests.Get(activeQuestID)
	if err != nil {
		return nil, fmt.Errorf("get quest: %w", err)
	}

	tasks, err := service.tasks.GetByQuestID(activeQuestID)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	quest.Tasks = tasks

	completedCount := len(quest.GetCompletedTasks())
	totalCount := len(quest.Tasks)
	percentage := quest.GetProgressPercentage()

	return &domain.Progress{
		QuestTitle:     quest.Title,
		CompletedTasks: completedCount,
		TotalTasks:     totalCount,
		Percentage:     percentage,
	}, nil
}

func (service *GameService) FinishQuest(userID int) (*domain.Reward, error) {
	activeQuestID, hasActive := service.activeQuests[userID]
	if !hasActive {
		return nil, fmt.Errorf("%w: user has no active quest to finish", domain.ErrInvalidState)
	}

	quest, err := service.quests.Get(activeQuestID)
	if err != nil {
		return nil, fmt.Errorf("get quest: %w", err)
	}

	user, err := service.users.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	tasks, err := service.tasks.GetByQuestID(activeQuestID)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	quest.Tasks = tasks

	if !quest.IsCompleted() {
		return nil, &domain.ValidationError{
			Field:   "quest",
			Message: "not all tasks are completed",
		}
	}

	xpReward := quest.Difficulty * 100
	user.AddXP(xpReward)

	if _, err := service.users.Update(userID, user); err != nil {
		return nil, fmt.Errorf("update user XP: %w", err)
	}

	rewardTitle := fmt.Sprintf("Completed: %s", quest.Title)
	reward, err := domain.NewReward(activeQuestID, rewardTitle, xpReward, "common")
	if err != nil {
		return nil, fmt.Errorf("create reward: %w", err)
	}

	delete(service.activeQuests, userID)

	return reward, nil
}
