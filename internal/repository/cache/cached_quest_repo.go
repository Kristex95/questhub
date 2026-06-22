package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Kristex95/questhub/internal/models"
	"github.com/Kristex95/questhub/internal/repository"
	"github.com/redis/go-redis/v9"
)

type CachedQuestRepository struct {
	next repository.QuestRepository
	rdb  *redis.Client
	ttl  time.Duration
}

func NewCachedQuestRepository(next repository.QuestRepository, rdb *redis.Client, ttl time.Duration) *CachedQuestRepository {
	return &CachedQuestRepository{next: next, rdb: rdb, ttl: ttl}
}

func questKey(id int64) string {
	return fmt.Sprintf("quest:%d", id)
}

func (r *CachedQuestRepository) GetByID(ctx context.Context, id int64) (*models.Quest, error) {
	key := questKey(id)

	cached, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var q models.Quest
		if jsonErr := json.Unmarshal([]byte(cached), &q); jsonErr == nil {
			return &q, nil
		}
	} else if !errors.Is(err, redis.Nil) {

	}

	quest, err := r.next.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data, mErr := json.Marshal(quest); mErr == nil {
		_ = r.rdb.Set(ctx, key, data, r.ttl).Err()
	}

	return quest, nil
}

func (r *CachedQuestRepository) Create(ctx context.Context, quest *models.Quest) (*models.Quest, error) {
	return r.next.Create(ctx, quest)
}

func (r *CachedQuestRepository) Update(ctx context.Context, quest *models.Quest) error {
	if err := r.next.Update(ctx, quest); err != nil {
		return err
	}
	_ = r.rdb.Del(ctx, questKey(quest.ID)).Err()
	return nil
}

func (r *CachedQuestRepository) Delete(ctx context.Context, id int64) error {
	if err := r.next.Delete(ctx, id); err != nil {
		return err
	}
	_ = r.rdb.Del(ctx, questKey(id)).Err()
	return nil
}

func (r *CachedQuestRepository) GetAll(ctx context.Context) ([]*models.Quest, error) {
	return r.next.GetAll(ctx)
}
