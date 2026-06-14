package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Kristex95/questhub/internal/domain"
)

func TestTaskRepository_CreateAndGetByID(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task := domain.Task{
		Title:    "Task 1",
		QuestId:  "q1",
		XPReward: 10,
	}

	created, err := repo.Create(task)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	got, err := repo.Get(created.ID)
	require.NoError(t, err)

	require.Equal(t, created.ID, got.ID)
	require.Equal(t, created.Title, got.Title)
	require.Equal(t, created.QuestId, got.QuestId)
	require.Equal(t, created.XPReward, got.XPReward)
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	_, err := repo.Get("missing")
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "Task", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}

func TestTaskRepository_GetAll_Empty(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	tasks, err := repo.GetAll()
	require.NoError(t, err)

	require.NotNil(t, tasks)
	require.Len(t, tasks, 0)
}

func TestTaskRepository_GetAll_WithData(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	t1 := domain.Task{Title: "T1", QuestId: "q1"}
	t2 := domain.Task{Title: "T2", QuestId: "q2"}

	_, _ = repo.Create(t1)
	_, _ = repo.Create(t2)

	tasks, err := repo.GetAll()
	require.NoError(t, err)

	require.Len(t, tasks, 2)

	found := map[string]bool{}
	for _, t := range tasks {
		found[t.Title] = true
	}

	require.True(t, found["T1"])
	require.True(t, found["T2"])
}

func TestTaskRepository_Update_Existing(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task := domain.Task{
		Title:    "Old",
		QuestId:  "q1",
		XPReward: 5,
	}

	created, _ := repo.Create(task)

	updated := domain.Task{
		Title:    "New",
		QuestId:  "q1",
		XPReward: 99,
	}

	res, err := repo.Update(created.ID, updated)
	require.NoError(t, err)

	require.Equal(t, "New", res.Title)
	require.Equal(t, 99, res.XPReward)

	got, _ := repo.Get(created.ID)
	require.Equal(t, "New", got.Title)
}

func TestTaskRepository_Update_NotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	_, err := repo.Update("missing", domain.Task{})
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "Task", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}

func TestTaskRepository_Delete_And_GetByID(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task := domain.Task{
		Title:   "To delete",
		QuestId: "q1",
	}

	created, _ := repo.Create(task)

	err := repo.Delete(created.ID)
	require.NoError(t, err)

	_, err = repo.Get(created.ID)
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))
	require.Equal(t, created.ID, notFoundErr.ID)
}

func TestTaskRepository_Delete_NotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	err := repo.Delete("missing")
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "Task", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}

func TestTaskRepository_GetByQuestID_Filtering(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	t1 := domain.Task{Title: "T1", QuestId: "q1"}
	t2 := domain.Task{Title: "T2", QuestId: "q1"}
	t3 := domain.Task{Title: "T3", QuestId: "q2"}

	_, _ = repo.Create(t1)
	_, _ = repo.Create(t2)
	_, _ = repo.Create(t3)

	result, err := repo.GetByQuestID("q1")
	require.NoError(t, err)

	require.Len(t, result, 2)

	for _, task := range result {
		require.Equal(t, "q1", task.QuestId)
	}
}

func TestTaskRepository_GetByQuestID_Empty(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	result, err := repo.GetByQuestID("missing")
	require.NoError(t, err)

	require.NotNil(t, result)
	require.Len(t, result, 0)
}