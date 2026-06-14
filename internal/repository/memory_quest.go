package repository

import (
	"fmt"
	"github.com/Kristex95/questhub/internal/domain"
)

type InMemoryQuestRepository struct {
	data    map[string]domain.Quest
	counter int
}

var _ QuestRepository = (*InMemoryQuestRepository)(nil)

func NewInMemoryQuestRepository() *InMemoryQuestRepository {
	return &InMemoryQuestRepository{
		data: make(map[string]domain.Quest),
	}
}

func (repo *InMemoryQuestRepository) Create(quest domain.Quest) (domain.Quest, error) {
	repo.counter++
	id := fmt.Sprintf("quest-%d", repo.counter)
	if _, exists := repo.data[id]; exists {
		return domain.Quest{}, &domain.DuplicateError{
			Entity: "Quest",
			Field:  "ID",
			Value:  id,
		}
	}
	quest.ID = id
	repo.data[id] = quest
	return quest, nil
}

func (repo *InMemoryQuestRepository) Get(id string) (domain.Quest, error) {
	quest, ok := repo.data[id]
	if !ok {
		return domain.Quest{},
			&domain.NotFoundError{
				Entity: "Quest",
				ID:     id,
			}
	}
	return quest, nil
}

func (repo *InMemoryQuestRepository) GetAll() ([]domain.Quest, error) {
	quests := make([]domain.Quest, 0, len(repo.data))
	for _, quest := range repo.data {
		quests = append(quests, quest)
	}
	return quests, nil
}

func (repo *InMemoryQuestRepository) Update(questID string, updatedQuest domain.Quest) (domain.Quest, error) {
	_, ok := repo.data[questID]
	if !ok {
		return domain.Quest{},
			&domain.NotFoundError{
				Entity: "Quest",
				ID:     questID,
			}
	}
	updatedQuest.ID = questID
	repo.data[questID] = updatedQuest
	return updatedQuest, nil
}

func (repo *InMemoryQuestRepository) Delete(questID string) error {
	if _, ok := repo.data[questID]; !ok {
		return &domain.NotFoundError{
			Entity: "Quest",
			ID: questID,
		}
	}
	delete(repo.data, questID)
	return nil
}