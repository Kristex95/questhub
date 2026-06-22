package models

import "time"

type Task struct {
	ID          int64       `db:"id" json:"id"`
	QuestID     int64       `db:"quest_id" json:"quest_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	IsCompleted bool      `db:"is_completed" json:"is_completed"`
	XPReward    int       `db:"xp_reward" json:"xp_reward"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
