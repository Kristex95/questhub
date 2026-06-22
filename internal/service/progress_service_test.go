package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Kristex95/questhub/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartProgress_Table(t *testing.T) {
	svc := NewProgressService(NewMockProgressRepository())

	cases := []struct {
		name    string
		userID  int64
		questID int64
		wantErr bool
	}{
		{"valid progress", 1, 2, false},
		{"another valid progress", 2, 3, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			progress, err := svc.StartProgress(context.Background(), tc.userID, tc.questID)
			require.NoError(t, err)
			require.NotNil(t, progress)
			assert.Equal(t, tc.userID, progress.UserID)
			assert.Equal(t, tc.questID, progress.QuestID)
			assert.Equal(t, "in_progress", progress.Status)
			assert.Equal(t, 0, progress.CompletedTasks)
		})
	}
}

func TestStartProgress_RepoError(t *testing.T) {
	mockRepo := NewMockProgressRepository()
	mockRepo.CreateErr = errors.New("database unavailable")
	svc := NewProgressService(mockRepo)

	progress, err := svc.StartProgress(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.Nil(t, progress)
}

func TestGetProgress_Table(t *testing.T) {
	mockRepo := NewMockProgressRepository()
	stored, err := mockRepo.Create(context.Background(), &models.Progress{UserID: 1, QuestID: 2, Status: "in_progress"})
	require.NoError(t, err)

	svc := NewProgressService(mockRepo)

	cases := []struct {
		name    string
		userID  int64
		questID int64
		wantErr bool
	}{
		{"progress exists", stored.UserID, stored.QuestID, false},
		{"progress not found", 10, 20, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			progress, err := svc.GetProgress(context.Background(), tc.userID, tc.questID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, progress)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, progress)
			assert.Equal(t, tc.userID, progress.UserID)
			assert.Equal(t, tc.questID, progress.QuestID)
		})
	}
}

func TestMarkCompleted_Error(t *testing.T) {
	mockRepo := NewMockProgressRepository()
	mockRepo.MarkCompletedErr = errors.New("cannot mark completed")
	svc := NewProgressService(mockRepo)

	err := svc.MarkCompleted(context.Background(), 1, 1)
	assert.Error(t, err)
}
