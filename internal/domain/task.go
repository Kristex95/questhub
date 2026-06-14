package domain

import "fmt"

type Task struct {
	ID          string
	Title       string
	Description string
	isCompleted bool
	XPReward    int
	Reward      *Reward
	QuestId 	string
}

func NewTask(id, title string, XPReward int, questId string) (*Task, error) {
	if title == "" {
		return nil, &ValidationError{Field: "title", Message: "must not be empty"}
	}
	if id == "" {
		return nil, &ValidationError{Field: "id", Message: "must not be empty"}
	}
	if questId == "" {
		return nil, &ValidationError{Field: "questId", Message: "must not be empty"}
	}
	if XPReward < 0 {
		XPReward = 0
	}
	return &Task{
			ID:       id,
			Title:    title,
			XPReward: XPReward,
			QuestId:  questId,
		},
		nil
}

func (t *Task) GetIsCompleted() bool {
	return t.isCompleted
}

func (t *Task) Print() {
	fmt.Printf("Task: %s | %s | %s\n", t.ID, t.Title, t.Description)
}
