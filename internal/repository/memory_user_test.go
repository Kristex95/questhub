package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Kristex95/questhub/internal/domain"
)

func TestUserRepository_CreateAndGetByID(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := domain.User{
		Username: "kirill",
		Email:    "test@mail.com",
	}

	created, err := repo.Create(u)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	got, err := repo.Get(created.ID)
	require.NoError(t, err)

	require.Equal(t, created.ID, got.ID)
	require.Equal(t, "kirill", got.Username)
	require.Equal(t, "test@mail.com", got.Email)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, err := repo.Get("missing")
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "User", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}

func TestUserRepository_GetAll_Empty(t *testing.T) {
	repo := NewInMemoryUserRepository()

	users, err := repo.GetAll()
	require.NoError(t, err)

	require.NotNil(t, users)
	require.Len(t, users, 0)
}

func TestUserRepository_GetAll_WithData(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u1 := domain.User{Username: "u1"}
	u2 := domain.User{Username: "u2"}

	_, _ = repo.Create(u1)
	_, _ = repo.Create(u2)

	users, err := repo.GetAll()
	require.NoError(t, err)

	require.Len(t, users, 2)

	found := map[string]bool{}
	for _, u := range users {
		found[u.Username] = true
	}

	require.True(t, found["u1"])
	require.True(t, found["u2"])
}

func TestUserRepository_Update_Existing(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := domain.User{
		Username: "old",
		Email:    "old@mail.com",
	}

	created, _ := repo.Create(u)

	updated := domain.User{
		Username: "new",
		Email:    "new@mail.com",
	}

	res, err := repo.Update(created.ID, updated)
	require.NoError(t, err)

	require.Equal(t, "new", res.Username)
	require.Equal(t, "new@mail.com", res.Email)

	got, _ := repo.Get(created.ID)
	require.Equal(t, "new", got.Username)
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, err := repo.Update("missing", domain.User{})
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "User", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}

func TestUserRepository_Delete_And_GetByID(t *testing.T) {
	repo := NewInMemoryUserRepository()

	u := domain.User{Username: "to-delete"}

	created, _ := repo.Create(u)

	err := repo.Delete(created.ID)
	require.NoError(t, err)

	_, err = repo.Get(created.ID)
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))
	require.Equal(t, created.ID, notFoundErr.ID)
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	err := repo.Delete("missing")
	require.Error(t, err)

	var notFoundErr *domain.NotFoundError
	require.True(t, errors.As(err, &notFoundErr))

	require.Equal(t, "User", notFoundErr.Entity)
	require.Equal(t, "missing", notFoundErr.ID)
}

func TestUserRepository_GetByUsername_NotFound_ShouldBeError(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, err := repo.GetByUsername("missing")

	require.Error(t, err)
}