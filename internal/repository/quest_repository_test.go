package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Kristex95/questhub/internal/domain"
)

func TestQuestRepository_CreateAndGetByID(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	q := domain.Quest{
		Title:       "Quest 1",
		Description: "desc",
		Difficulty:  5,
	}

	created, err := repo.Create(q)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	got, err := repo.Get(created.ID)
	require.NoError(t, err)

	require.Equal(t, created.ID, got.ID)
	require.Equal(t, created.Title, got.Title)
	require.Equal(t, created.Description, got.Description)
	require.Equal(t, created.Difficulty, got.Difficulty)
}

func TestQuestRepository_GetByID_NotFound(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	_, err := repo.Get("missing-id")
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "Quest", notFoundErr.Entity)
	require.Equal(t, "missing-id", notFoundErr.ID)
}

func TestQuestRepository_GetAll_Empty(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	quests, err := repo.GetAll()
	require.NoError(t, err)

	require.NotNil(t, quests)
	require.Len(t, quests, 0)
}

func TestQuestRepository_GetAll_WithData(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	q1 := domain.Quest{Title: "Q1", Description: "d1", Difficulty: 1}
	q2 := domain.Quest{Title: "Q2", Description: "d2", Difficulty: 2}

	_, _ = repo.Create(q1)
	_, _ = repo.Create(q2)

	quests, err := repo.GetAll()
	require.NoError(t, err)

	require.Len(t, quests, 2)

	titles := map[string]bool{}
	for _, q := range quests {
		titles[q.Title] = true
	}

	require.True(t, titles["Q1"])
	require.True(t, titles["Q2"])
}

func TestQuestRepository_Update_Existing(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	q := domain.Quest{
		Title:       "Old",
		Description: "old desc",
		Difficulty:  1,
	}

	created, _ := repo.Create(q)

	updated := domain.Quest{
		Title:       "New",
		Description: "new desc",
		Difficulty:  10,
	}

	res, err := repo.Update(created.ID, updated)
	require.NoError(t, err)

	require.Equal(t, created.ID, res.ID)
	require.Equal(t, "New", res.Title)
	require.Equal(t, "new desc", res.Description)
	require.Equal(t, 10, res.Difficulty)

	// verify persistence
	got, _ := repo.Get(created.ID)
	require.Equal(t, "New", got.Title)
}

func TestQuestRepository_Update_NotFound(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	_, err := repo.Update("missing", domain.Quest{})
	require.Error(t, err)

	var nf *domain.NotFoundError
	require.True(t, errors.As(err, &nf))

	require.Equal(t, "Quest", nf.Entity)
	require.Equal(t, "missing", nf.ID)
}

func TestQuestRepository_Delete_And_GetByID(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	q := domain.Quest{
		Title:       "To delete",
		Description: "desc",
		Difficulty:  3,
	}

	created, _ := repo.Create(q)

	err := repo.Delete(created.ID)
	require.NoError(t, err)

	_, err = repo.Get(created.ID)
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))
	require.Equal(t, created.ID, notFoundErr.ID)
}

func TestQuestRepository_Delete_NotFound(t *testing.T) {
	repo := NewInMemoryQuestRepository()

	err := repo.Delete("missing")
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "Quest", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}