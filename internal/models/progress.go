package models

import "time"

type Progress struct {
	ID             int64      `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"user_id"`
	QuestID        int64      `db:"quest_id" json:"quest_id"`
	Status         string     `db:"status" json:"status"`
	CompletedTasks int        `db:"completed_tasks" json:"completed_tasks"`
	StartedAt      time.Time  `db:"started_at" json:"started_at"`
	CompletedAt    *time.Time `db:"completed_at" json:"completed_at,omitempty"`
}
