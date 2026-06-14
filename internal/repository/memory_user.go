package repository

import (
	"fmt"
	"github.com/Kristex95/questhub/internal/domain"
)

type InMemoryUserRepository struct {
	data    map[string]domain.User
	counter int
}

var _ UserRepository = (*InMemoryUserRepository)(nil)

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		data: make(map[string]domain.User),
	}
}

func (repo *InMemoryUserRepository) Create(user domain.User) (domain.User, error) {
	repo.counter++
	id := fmt.Sprintf("user-%d", repo.counter)
	if _, exists := repo.data[id]; exists {
		return domain.User{}, &domain.DuplicateError{
			Entity: "User",
			Field:  "ID",
			Value:  id,
		}
	}
	user.ID = id
	repo.data[id] = user
	return user, nil
}

func (repo *InMemoryUserRepository) Get(id string) (domain.User, error) {
	user, ok := repo.data[id]
	if !ok {
		return domain.User{},
			&domain.NotFoundError{
				Entity: "User",
				ID:     id,
			}
	}
	return user, nil
}

func (repo *InMemoryUserRepository) GetAll() ([]domain.User, error) {
	users := make([]domain.User, 0, len(repo.data))
	for _, user := range repo.data {
		users = append(users, user)
	}
	return users, nil
}

func (repo *InMemoryUserRepository) GetByUsername(username string) (*domain.User, error) {
	if username == "" {
		return &domain.User{}, nil
	} 
	for _, user := range repo.data {
		if user.Username == username {
			return &user, nil
		}
	}
	return &domain.User{}, &domain.NotFoundError{Entity: "User"}
}

func (repo *InMemoryUserRepository) Update(userID string, updatedUser domain.User) (domain.User, error) {
	_, ok := repo.data[userID]
	if !ok {
		return domain.User{},
			&domain.NotFoundError{
				Entity: "User",
				ID:     userID,
			}
	}
	updatedUser.ID = userID
	repo.data[userID] = updatedUser
	return updatedUser, nil
}

func (repo *InMemoryUserRepository) Delete(userID string) error {
	if _, ok := repo.data[userID]; !ok {
		return &domain.NotFoundError{
			Entity: "User",
			ID:     userID,
		}
	}
	delete(repo.data, userID)
	return nil
}
