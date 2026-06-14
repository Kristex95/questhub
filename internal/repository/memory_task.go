package repository

import (
	"fmt"
	"github.com/Kristex95/questhub/internal/domain"
)

type InMemoryTaskRepository struct {
	data    map[string]domain.Task
	counter int
}

var _ TaskRepository = (*InMemoryTaskRepository)(nil)

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		data: make(map[string]domain.Task),
	}
}

func (repo *InMemoryTaskRepository) Create(task domain.Task) (domain.Task, error) {
	repo.counter++
	id := fmt.Sprintf("task-%d", repo.counter)
	if _, exists := repo.data[id]; exists {
		return domain.Task{}, &domain.DuplicateError{
			Entity: "Task",
			Field:  "ID",
			Value:  id,
		}
	}
	task.ID = id
	repo.data[id] = task
	return task, nil
}

func (repo *InMemoryTaskRepository) Get(id string) (domain.Task, error) {
	task, ok := repo.data[id]
	if !ok {
		return domain.Task{},
			&domain.NotFoundError{
				Entity: "Task",
				ID:     id,
			}
	}
	return task, nil
}

func (repo *InMemoryTaskRepository) GetAll() ([]domain.Task, error) {
	tasks := make([]domain.Task, 0, len(repo.data))
	for _, task := range repo.data {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repo *InMemoryTaskRepository) GetByQuestID(questID string) ([]*domain.Task, error) {
	result := make([]*domain.Task, 0)

	for _, task := range repo.data {
		if task.QuestId == questID {
			t := task
			result = append(result, &t)
		}
	}

	return result, nil
}

func (repo *InMemoryTaskRepository) Update(taskID string, updatedtask domain.Task) (domain.Task, error) {
	_, ok := repo.data[taskID]
	if !ok {
		return domain.Task{},
			&domain.NotFoundError{
				Entity: "Task",
				ID:     taskID,
			}
	}
	updatedtask.ID = taskID
	repo.data[taskID] = updatedtask
	return updatedtask, nil
}

func (repo *InMemoryTaskRepository) Delete(taskID string) error {
	if _, ok := repo.data[taskID]; !ok {
		return &domain.NotFoundError{
			Entity: "Task",
			ID: taskID,
		}
	}
	delete(repo.data, taskID)
	return nil
}