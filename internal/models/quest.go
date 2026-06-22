package models

import "time"

type Quest struct {
	ID          int64     `json:"id"          db:"id"`
	Title       string    `json:"title"       db:"title"`
	Description string    `json:"description" db:"description"`
	Difficulty  int       `json:"difficulty"  db:"difficulty"`
	IsActive    bool      `json:"is_active"   db:"is_active"` 
	XPReward    int       `json:"xp_reward"   db:"xp_reward"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"  db:"updated_at"`
}
