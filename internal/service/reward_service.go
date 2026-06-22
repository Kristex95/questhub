package service

import (
	"context"
	"fmt"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/Kristex95/questhub/internal/repository"
)

type leaderboardUpdater interface {
	UpdateLeaderboard(ctx context.Context, userID int64, xp int) error
}

var validRewardTypes = map[string]bool{
	"xp":       true,
	"currency": true,
	"item":     true,
}

type RewardService struct {
	rewards repository.RewardRepository
	users   repository.UserRepository
	stats   leaderboardUpdater
}

func NewRewardService(rewards repository.RewardRepository, users repository.UserRepository, stats leaderboardUpdater) *RewardService {
	return &RewardService{rewards: rewards, users: users, stats: stats}
}

func (s *RewardService) CreateReward(ctx context.Context, questID int64, rewardType string, xpAmount int, itemName string) (*models.Reward, error) {
	if !validRewardTypes[rewardType] {
		return nil, fmt.Errorf("create reward: %w", &domain.ValidationError{
			Field: "reward_type", Message: "must be one of: xp, currency, item",
		})
	}
	if rewardType == "xp" && xpAmount <= 0 {
		return nil, fmt.Errorf("create reward: %w", &domain.ValidationError{
			Field: "xp_amount", Message: "must be greater than 0 for xp reward",
		})
	}
	if rewardType == "item" && itemName == "" {
		return nil, fmt.Errorf("create reward: %w", &domain.ValidationError{
			Field: "item_name", Message: "must not be empty for item reward",
		})
	}

	reward := &models.Reward{
		QuestID:    questID,
		RewardType: rewardType,
		XPAmount:   xpAmount,
		ItemName:   &itemName,
	}

	created, err := s.rewards.Create(ctx, reward)
	if err != nil {
		return nil, fmt.Errorf("create reward: %w", err)
	}

	return created, nil
}

func (s *RewardService) GrantQuestRewards(ctx context.Context, userID, questID int64) (int, error) {
	rewards, err := s.rewards.GetByQuestID(ctx, questID)
	if err != nil {
		return 0, fmt.Errorf("grant quest rewards: %w", err)
	}

	totalXP := 0
	for _, rw := range rewards {
		if rw.RewardType == "xp" {
			totalXP += rw.XPAmount
		}
	}

	if totalXP > 0 {
		if err := s.users.AddXP(ctx, userID, totalXP); err != nil {
			return 0, fmt.Errorf("grant quest rewards: %w", err)
		}

		user, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return 0, fmt.Errorf("grant quest rewards: %w", err)
		}
		if err := s.stats.UpdateLeaderboard(ctx, userID, user.XP); err != nil {
			return 0, fmt.Errorf("grant quest rewards: %w", err)
		}
	}

	return totalXP, nil
}
