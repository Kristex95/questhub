package main

import "fmt"

type Task struct {
    ID          string
    Title       string
    Description string
    isCompleted bool
    XPReward    int
    Reward      *Reward
}

func NewTask(id, title string, XPReward int) *Task {
    if XPReward < 0 {
       XPReward = 0
    }
    return &Task{
       ID:       id,
       Title:    title,
       XPReward: XPReward,
    }
}

func (t *Task) GetIsCompleted() bool {
	return t.isCompleted
}

func (t *Task) Print() {
	fmt.Printf("Task: %s | %s | %s\n", t.ID, t.Title, t.Description)
}