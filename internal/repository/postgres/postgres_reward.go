package postgres

import (
	"context"
	"fmt"

	"github.com/Kristex95/questhub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RewardRepository struct {
	pool *pgxpool.Pool
}

func NewRewardRepository(pool *pgxpool.Pool) *RewardRepository {
	return &RewardRepository{pool: pool}
}

func (r *RewardRepository) Create(ctx context.Context, reward *models.Reward) (*models.Reward, error) {
	const query = `
		INSERT INTO rewards (quest_id, reward_type, xp_amount, item_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		reward.QuestID, reward.RewardType, reward.XPAmount, reward.ItemName,
	).Scan(&reward.ID, &reward.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create reward: %w", err)
	}

	return reward, nil
}

func (r *RewardRepository) GetByQuestID(ctx context.Context, questID int64) ([]*models.Reward, error) {
	const query = `
		SELECT id, quest_id, reward_type, xp_amount, COALESCE(item_name, '') , created_at
		FROM rewards WHERE quest_id = $1 ORDER BY id`

	rows, err := r.pool.Query(ctx, query, questID)
	if err != nil {
		return nil, fmt.Errorf("get rewards by quest: %w", err)
	}
	defer rows.Close()

	rewards := make([]*models.Reward, 0)
	for rows.Next() {
		var rw models.Reward
		if err := rows.Scan(
			&rw.ID, &rw.QuestID, &rw.RewardType, &rw.XPAmount, &rw.ItemName, &rw.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan reward: %w", err)
		}
		rewards = append(rewards, &rw)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rewards: %w", err)
	}

	return rewards, nil
}
