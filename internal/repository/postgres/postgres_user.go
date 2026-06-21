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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO users (username, email, xp, level)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		user.Username, user.Email, user.XP, user.Level,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	const query = `SELECT id, username, email, xp, level, created_at FROM users WHERE id = $1`

	var u models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.XP, &u.Level, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(id, 10)}
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	const query = `SELECT id, username, email, xp, level, created_at FROM users WHERE username = $1`

	var u models.User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.XP, &u.Level, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &domain.NotFoundError{Entity: "User", Value: username}
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	const query = `SELECT id, username, email, xp, level, created_at FROM users ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.XP, &u.Level, &u.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	const query = `UPDATE users SET username = $1, email = $2, xp = $3, level = $4 WHERE id = $5`

	tag, err := r.pool.Exec(ctx, query, user.Username, user.Email, user.XP, user.Level, user.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "User", Value: strconv.Itoa(user.ID)}
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(id, 10)}
	}

	return nil
}

func (r *UserRepository) AddXP(ctx context.Context, userID int64, amount int) error {
	const query = `UPDATE users SET xp = xp + $2 WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("add xp: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(userID, 10)}
	}

	return nil
}
