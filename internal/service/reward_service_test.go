package service

import (
	"context"
	"testing"

	"github.com/Kristex95/questhub/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateReward_Validation(t *testing.T) {
	svc := NewRewardService(NewMockRewardRepository(), NewMockUserRepository(), &MockLeaderboardUpdater{})

	cases := []struct {
		name       string
		rewardType string
		xpAmount   int
		itemName   string
		wantErr    bool
	}{
		{"valid xp reward", "xp", 50, "", false},
		{"invalid reward type", "gold", 0, "", true},
		{"item reward without name", "item", 0, "", true},
		{"xp reward with zero amount", "xp", 0, "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reward, err := svc.CreateReward(context.Background(), 1, tc.rewardType, tc.xpAmount, tc.itemName)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, reward)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, reward)
			assert.Equal(t, int64(1), reward.QuestID)
			assert.Equal(t, tc.rewardType, reward.RewardType)
			if tc.rewardType == "xp" {
				assert.Equal(t, tc.xpAmount, reward.XPAmount)
			}
		})
	}
}

func TestGrantQuestRewards_Table(t *testing.T) {
	tests := []struct {
		name           string
		rewards        []*models.Reward
		initialXP      int
		wantGranted    int
		wantFinalXP    int
		wantErr        bool
		leaderboardErr bool
	}{
		{"no xp rewards", []*models.Reward{{QuestID: 1, RewardType: "currency", ItemName: strPtr("Gold Coin")}}, 0, 0, 0, false, false},
		{"single xp reward", []*models.Reward{{QuestID: 1, RewardType: "xp", XPAmount: 100}}, 0, 100, 100, false, false},
		{"multiple xp rewards", []*models.Reward{{QuestID: 1, RewardType: "xp", XPAmount: 20}, {QuestID: 1, RewardType: "xp", XPAmount: 30}}, 50, 50, 100, false, false},
		{"leaderboard update fails", []*models.Reward{{QuestID: 1, RewardType: "xp", XPAmount: 10}}, 0, 0, 0, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rewardRepo := NewMockRewardRepository()
			rewardRepo.data[1] = tc.rewards
			userRepo := NewMockUserRepository()
			user, err := userRepo.Create(context.Background(), &models.User{Username: "player", XP: tc.initialXP})
			require.NoError(t, err)

			var stats leaderboardUpdater = &MockLeaderboardUpdater{}
			if tc.leaderboardErr {
				stats = &FailingLeaderboardUpdater{err: assert.AnError}
			}

			svc := NewRewardService(rewardRepo, userRepo, stats)

			granted, err := svc.GrantQuestRewards(context.Background(), int64(user.ID), 1)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, granted)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantGranted, granted)
			updated, err := userRepo.GetByID(context.Background(), int64(user.ID))
			require.NoError(t, err)
			assert.Equal(t, tc.wantFinalXP, updated.XP)
		})
	}
}

func strPtr(s string) *string {
	return &s
}

type FailingLeaderboardUpdater struct {
	err error
}

func (m *FailingLeaderboardUpdater) UpdateLeaderboard(ctx context.Context, userID int64, xp int) error {
	return m.err
}
