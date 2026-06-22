package stats

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

func NewStatsService(rdb *redis.Client) *Service {
	return &Service{rdb: rdb}
}

func (s *Service) IncrCompletedQuests(ctx context.Context) error {
	if err := s.rdb.Incr(ctx, "stats:completed_quests").Err(); err != nil {
		return fmt.Errorf("incr completed quests: %w", err)
	}
	return nil
}

func (s *Service) GetCompletedQuests(ctx context.Context) (int64, error) {
	n, err := s.rdb.Get(ctx, "stats:completed_quests").Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get completed quests: %w", err)
	}
	return n, nil
}

func (s *Service) UpdateLeaderboard(ctx context.Context, userID int64, xp int) error {
	member := fmt.Sprintf("user:%d", userID)
	if err := s.rdb.ZAdd(ctx, "leaderboard", redis.Z{
		Score:  float64(xp),
		Member: member,
	}).Err(); err != nil {
		return fmt.Errorf("update leaderboard: %w", err)
	}
	return nil
}

type LeaderboardEntry struct {
	UserID int64 `json:"user_id"`
	XP     int   `json:"xp"`
}

func (s *Service) TopPlayers(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	res, err := s.rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("best players: %w", err)
	}

	entries := make([]LeaderboardEntry, 0, len(res))
	for _, z := range res {
		var id int64
		if _, scanErr := fmt.Sscanf(z.Member.(string), "user:%d", &id); scanErr != nil {
			continue
		}
		entries = append(entries, LeaderboardEntry{UserID: id, XP: int(z.Score)})
	}

	return entries, nil
}
