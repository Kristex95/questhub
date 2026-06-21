package domain

import (
	"fmt"
	"strconv"
)

type Quest struct {
	ID          int
	Title       string
	Description string
	Difficulty  int
	isActive    bool
	XPReward    int
	Tasks       []*Task
}

func NewQuest(id int, title, description string, difficulty int) (*Quest, error) {
	if len(title) < 3 {
		return nil, &ValidationError{Field: "title", Message: "must be at least 3 characters"}
	}
	if description == "" {
		return nil, &ValidationError{Field: "description", Message: "must not be empty"}
	}
	if difficulty < 1 || difficulty > 10 {
		return nil, &ValidationError{Field: "difficulty", Message: "must be between 1 and 10"}
	}
	return &Quest{
			ID:          id,
			Title:       title,
			Description: description,
			Difficulty:  difficulty,
			XPReward:    difficulty * 100,
			isActive:    true,
			Tasks:       make([]*Task, 0),
		},
		nil
}

func (q *Quest) AddTask(task *Task) {
	q.Tasks = append(q.Tasks, task)
}

func (q *Quest) CompleteTask(taskId int) error {
	for i := range q.Tasks {
		if q.Tasks[i].ID == taskId {
			q.Tasks[i].isCompleted = true
			return nil
		}
	}
	return &NotFoundError{Entity: "task", Value: strconv.Itoa(taskId)}
}

func (q *Quest) Summary() string {
	completed := 0
	for _, i := range q.Tasks {
		if i.isCompleted {
			completed++
		}
	}
	status := "active"
	if !q.isActive {
		status = "inactive"
	}
	return fmt.Sprintf("[q%d] %s | Difficulty: %d | Progress : %d/%d | XP: %d | %s", q.ID, q.Title, q.Difficulty, completed, len(q.Tasks), q.TotalXP(), status)

}

func (q *Quest) TotalXP() int {
	total := q.XPReward
	for _, i := range q.Tasks {
		total += i.XPReward
		if i.Reward != nil {
			total += i.Reward.XPAmount
		}
	}
	return total
}

func (q *Quest) Activate() {
	q.isActive = true
}

func (q *Quest) Deactivate() {
	q.isActive = false
}

func (q *Quest) GetCompletedTasks() []*Task {
	completed := make([]*Task, 0, len(q.Tasks))
	for _, task := range q.Tasks {
		if task.isCompleted {
			completed = append(completed, task)
		}
	}
	return completed
}

func (q *Quest) GetRemainingTasks() []*Task {
	remaining := make([]*Task, 0, len(q.Tasks))
	for _, task := range q.Tasks {
		if !task.isCompleted {
			remaining = append(remaining, task)
		}
	}
	return remaining
}

func (q *Quest) GetProgressPercentage() float64 {
	completedCount := len(q.GetCompletedTasks())
	totalCount := len(q.Tasks)
	if completedCount == 0 {
		return 0.0
	}
	return float64(completedCount) / float64(totalCount) * 100
}

func (q *Quest) IsCompleted() bool {
	for _, task := range q.Tasks {
		if !task.GetIsCompleted() {
			return false
		}
	}
	return true
}
