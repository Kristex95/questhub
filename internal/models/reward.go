package models

import "time"

type Reward struct {
	ID         int64     `db:"id" json:"id"`
	QuestID    int64     `db:"quest_id" json:"quest_id"`
	RewardType string    `db:"reward_type" json:"reward_type"`
	XPAmount   int       `db:"xp_amount" json:"xp_amount"`
	ItemName   *string   `db:"item_name" json:"item_name,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
