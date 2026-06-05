package main

import "fmt"

type Quest struct {
    ID          string
    Title       string
    Description string
    Difficulty  int
    isActive    bool
    XPReward    int
    Tasks       []*Task
}

func NewQuest(id, title, description string, difficulty int) *Quest {
    if difficulty < 1 {
       difficulty = 1
    }
    if difficulty > 10 {
       difficulty = 10
    }
    return &Quest{
       ID:          id,
       Title:       title,
       Description: description,
	   Difficulty:  difficulty,
       XPReward:    difficulty * 100,
       isActive:    true,
       Tasks:       make([]*Task, 0),
    }
}

func (q *Quest) AddTask(task *Task) {
    q.Tasks = append(q.Tasks, task)
}

func (q *Quest) CompleteTask(taskId string) bool {
    for i := range q.Tasks {
       if q.Tasks[i].ID == taskId {
          q.Tasks[i].isCompleted = true
          return true
       }
    }
    return false
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
    return fmt.Sprintf("[%s] %s | Difficulty: %d | Progress : %d/%d | XP: %d | %s", q.ID, q.Title, q.Difficulty, completed, len(q.Tasks), q.TotalXP(), status)

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
