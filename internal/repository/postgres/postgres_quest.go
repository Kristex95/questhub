package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestRepository struct {
	pool *pgxpool.Pool
}

func NewQuestRepository(pool *pgxpool.Pool) *QuestRepository {
	return &QuestRepository{pool: pool}
}

func (r *QuestRepository) Create(ctx context.Context, quest *models.Quest) (*models.Quest, error) {
	const query = `
		INSERT INTO quests (title, description, difficulty, xp_reward, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		quest.Title, quest.Description, quest.Difficulty, quest.XPReward, quest.IsActive,
	).Scan(&quest.ID, &quest.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create quest: %w", err)
	}

	return quest, nil
}

func (r *QuestRepository) GetByID(ctx context.Context, id int64) (*models.Quest, error) {
	const query = `
		SELECT id, title, description, difficulty, xp_reward, is_active, created_at
		FROM quests WHERE id = $1`

	var q models.Quest
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.Title, &q.Description, &q.Difficulty, &q.XPReward, &q.IsActive, &q.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &domain.NotFoundError{Entity: "Quest", Value: strconv.FormatInt(id, 10)}
	}
	if err != nil {
		return nil, fmt.Errorf("get quest by id: %w", err)
	}

	return &q, nil
}

func (r *QuestRepository) GetAll(ctx context.Context) ([]*models.Quest, error) {
	const query = `
		SELECT id, title, description, difficulty, xp_reward, is_active, created_at
		FROM quests ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all quests: %w", err)
	}
	defer rows.Close()

	quests := make([]*models.Quest, 0)
	for rows.Next() {
		var q models.Quest
		if err := rows.Scan(
			&q.ID, &q.Title, &q.Description, &q.Difficulty, &q.XPReward, &q.IsActive, &q.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan quest: %w", err)
		}
		quests = append(quests, &q)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate quests: %w", err)
	}

	return quests, nil
}

func (r *QuestRepository) Update(ctx context.Context, quest *models.Quest) error {
	const query = `
		UPDATE quests
		SET title = $1, description = $2, difficulty = $3, xp_reward = $4, is_active = $5
		WHERE id = $6`

	tag, err := r.pool.Exec(ctx, query,
		quest.Title, quest.Description, quest.Difficulty, quest.XPReward, quest.IsActive, quest.ID,
	)
	if err != nil {
		return fmt.Errorf("update quest: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "Quest", Value: strconv.Itoa(quest.ID)}
	}

	return nil
}

func (r *QuestRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM quests WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete quest: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "Quest", Value: strconv.FormatInt(id, 10)}
	}

	return nil
}

type repositoryQuest interface {
	Create(ctx context.Context, quest *models.Quest) (*models.Quest, error)
	GetByID(ctx context.Context, id int64) (*models.Quest, error)
	GetAll(ctx context.Context) ([]*models.Quest, error)
	Update(ctx context.Context, quest *models.Quest) error
	Delete(ctx context.Context, id int64) error
}

var _ repositoryQuest = (*QuestRepository)(nil)
