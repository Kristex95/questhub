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

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	const query = `
		INSERT INTO tasks (quest_id, title, description, is_completed, xp_reward)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.pool.QueryRow(ctx, query,
		task.QuestID, task.Title, task.Description, task.IsCompleted, task.XPReward,
	).Scan(&task.ID)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	const query = `
		SELECT id, quest_id, title, description, is_completed, xp_reward
		FROM tasks WHERE id = $1`

	var t models.Task
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.QuestID, &t.Title, &t.Description, &t.IsCompleted, &t.XPReward,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &domain.NotFoundError{Entity: "Task", Value: strconv.FormatInt(id, 10)}
	}
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	return &t, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	const query = `
		SELECT id, quest_id, title, description, is_completed, xp_reward
		FROM tasks ORDER BY id`

	return r.queryTasks(ctx, query)
}

func (r *TaskRepository) GetByQuestID(ctx context.Context, questID int64) ([]*models.Task, error) {
	const query = `
		SELECT id, quest_id, title, description, is_completed, xp_reward
		FROM tasks WHERE quest_id = $1 ORDER BY id`

	return r.queryTasks(ctx, query, questID)
}

func (r *TaskRepository) queryTasks(ctx context.Context, query string, args ...any) ([]*models.Task, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0)
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID, &t.QuestID, &t.Title, &t.Description, &t.IsCompleted, &t.XPReward,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	const query = `
		UPDATE tasks
		SET title = $1, description = $2, is_completed = $3, xp_reward = $4
		WHERE id = $5`

	tag, err := r.pool.Exec(ctx, query,
		task.Title, task.Description, task.IsCompleted, task.XPReward, task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "Task", Value: strconv.FormatInt(task.ID, 10)}
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "Task", Value: strconv.FormatInt(id, 10)}
	}

	return nil
}
