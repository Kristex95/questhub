package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepository struct {
	pool *pgxpool.Pool
}

func NewProgressRepository(pool *pgxpool.Pool) *ProgressRepository {
	return &ProgressRepository{pool: pool}
}

func (r *ProgressRepository) Create(ctx context.Context, p *models.Progress) (*models.Progress, error) {
	const query = `
		INSERT INTO progress (user_id, quest_id, status, completed_tasks)
		VALUES ($1, $2, $3, $4)
		RETURNING id, started_at`

	err := r.pool.QueryRow(ctx, query,
		p.UserID, p.QuestID, p.Status, p.CompletedTasks,
	).Scan(&p.ID, &p.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("create progress: %w", err)
	}

	return p, nil
}

func (r *ProgressRepository) GetByUserAndQuest(ctx context.Context, userID, questID int64) (*models.Progress, error) {
	const query = `
		SELECT id, user_id, quest_id, status, completed_tasks, started_at, completed_at
		FROM progress WHERE user_id = $1 AND quest_id = $2`

	var p models.Progress
	err := r.pool.QueryRow(ctx, query, userID, questID).Scan(
		&p.ID, &p.UserID, &p.QuestID, &p.Status, &p.CompletedTasks, &p.StartedAt, &p.CompletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &domain.NotFoundError{Entity: "Progress", Value: fmt.Sprintf("user=%d quest=%d", userID, questID)}
	}
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}

	return &p, nil
}

func (r *ProgressRepository) MarkCompleted(ctx context.Context, userID, questID int64) error {
	const query = `
		UPDATE progress
		SET status = 'completed', completed_at = NOW()
		WHERE user_id = $1 AND quest_id = $2`

	tag, err := r.pool.Exec(ctx, query, userID, questID)
	if err != nil {
		return fmt.Errorf("mark progress completed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "Progress", Value: fmt.Sprintf("user=%d quest=%d", userID, questID)}
	}

	return nil
}
