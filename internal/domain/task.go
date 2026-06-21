package domain

import "fmt"

type Task struct {
	ID          int
	Title       string
	Description string
	isCompleted bool
	XPReward    int
	Reward      *Reward
	QuestId     int
}

func NewTask(id int, title string, XPReward int, questId int) (*Task, error) {
	if title == "" {
		return nil, &ValidationError{Field: "title", Message: "must not be empty"}
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
	fmt.Printf("Task: %d | %s | %s\n", t.ID, t.Title, t.Description)
}
