package service

import (
	"errors"
	"testing"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuestStore struct {
	items map[int]domain.Quest
}

func (m *mockQuestStore) Create(quest domain.Quest) (domain.Quest, error) { return quest, nil }
func (m *mockQuestStore) Get(id int) (domain.Quest, error) {
	quest, ok := m.items[id]
	if !ok {
		return domain.Quest{}, errors.New("quest not found")
	}
	return quest, nil
}
func (m *mockQuestStore) Update(id int, quest domain.Quest) (domain.Quest, error) {
	m.items[id] = quest
	return quest, nil
}

type mockTaskStore struct {
	items map[int]*domain.Task
}

func (m *mockTaskStore) Create(task domain.Task) (domain.Task, error) {
	m.items[task.ID] = &task
	return task, nil
}
func (m *mockTaskStore) Get(id int) (domain.Task, error) {
	task, ok := m.items[id]
	if !ok {
		return domain.Task{}, errors.New("task not found")
	}
	return *task, nil
}
func (m *mockTaskStore) GetByQuestID(questID int) ([]*domain.Task, error) {
	var tasks []*domain.Task
	for _, task := range m.items {
		if task.QuestId == questID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
func (m *mockTaskStore) Update(id int, task domain.Task) (domain.Task, error) {
	m.items[id] = &task
	return task, nil
}

type mockUserStore struct {
	items map[int]domain.User
}

func (m *mockUserStore) Create(user domain.User) (domain.User, error) { return user, nil }
func (m *mockUserStore) Get(id int) (domain.User, error) {
	user, ok := m.items[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}
func (m *mockUserStore) Update(id int, user domain.User) (domain.User, error) {
	m.items[id] = user
	return user, nil
}

func newGameService(t *testing.T) (*GameService, *mockQuestStore, *mockTaskStore, *mockUserStore) {
	quests := &mockQuestStore{items: make(map[int]domain.Quest)}
	tasks := &mockTaskStore{items: make(map[int]*domain.Task)}
	users := &mockUserStore{items: make(map[int]domain.User)}
	return NewGameService(quests, tasks, users), quests, tasks, users
}

func TestGameService_StartQuest(t *testing.T) {
	svc, quests, _, users := newGameService(t)
	users.items[1] = domain.User{ID: 1, Username: "player"}
	quests.items[1] = domain.Quest{ID: 1, Title: "Test"}

	cases := []struct {
		name      string
		userID    int
		questID   int
		wantError bool
	}{
		{"success", 1, 1, false},
		{"invalid user", 2, 1, true},
		{"invalid quest", 1, 2, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.StartQuest(tc.userID, tc.questID)
			if tc.wantError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

	err := svc.StartQuest(1, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrAlreadyStarted)
}

func TestGameService_CompleteTask_Success(t *testing.T) {
	svc, quests, tasks, users := newGameService(t)
	users.items[1] = domain.User{ID: 1, Username: "player"}
	quests.items[1] = domain.Quest{ID: 1, Title: "Test", Tasks: []*domain.Task{{ID: 10, QuestId: 1}}}
	tasks.items[10] = &domain.Task{ID: 10, Title: "Task 10", QuestId: 1}

	err := svc.StartQuest(1, 1)
	require.NoError(t, err)

	err = svc.CompleteTask(1, 10)
	require.NoError(t, err)

	updatedTask, err := tasks.Get(10)
	require.NoError(t, err)
	assert.True(t, updatedTask.GetIsCompleted())
}

func TestGameService_GetProgress(t *testing.T) {
	svc, quests, tasks, users := newGameService(t)
	users.items[1] = domain.User{ID: 1, Username: "player"}
	quests.items[1] = domain.Quest{ID: 1, Title: "Test", Tasks: []*domain.Task{{ID: 10, QuestId: 1}}}
	tasks.items[10] = &domain.Task{ID: 10, Title: "Task 10", QuestId: 1}

	err := svc.StartQuest(1, 1)
	require.NoError(t, err)

	progress, err := svc.GetProgress(1)
	require.NoError(t, err)
	assert.Equal(t, "Test", progress.QuestTitle)
	assert.Equal(t, 0, progress.CompletedTasks)
	assert.Equal(t, 1, progress.TotalTasks)
}

func TestGameService_FinishQuest_Success(t *testing.T) {
	svc, quests, tasks, users := newGameService(t)
	users.items[1] = domain.User{ID: 1, Username: "player", XP: 0}
	taskPtr := &domain.Task{ID: 10, Title: "Task 10", QuestId: 1}
	quests.items[1] = domain.Quest{ID: 1, Title: "Test", Difficulty: 1, Tasks: []*domain.Task{taskPtr}}
	tasks.items[10] = taskPtr

	quest := quests.items[1]
	err := quest.CompleteTask(10)
	require.NoError(t, err)
	quests.items[1] = quest

	err = svc.StartQuest(1, 1)
	require.NoError(t, err)

	reward, err := svc.FinishQuest(1)
	require.NoError(t, err)
	assert.Equal(t, "Completed: Test", reward.Title)
	assert.Equal(t, 100, users.items[1].XP)
}

func TestGameService_FinishQuest_ValidationError(t *testing.T) {
	svc, quests, tasks, users := newGameService(t)
	users.items[1] = domain.User{ID: 1, Username: "player"}
	quests.items[1] = domain.Quest{ID: 1, Title: "Test", Difficulty: 1, Tasks: []*domain.Task{{ID: 10, QuestId: 1}}}
	tasks.items[10] = &domain.Task{ID: 10, Title: "Task 10", QuestId: 1}

	err := svc.StartQuest(1, 1)
	require.NoError(t, err)

	reward, err := svc.FinishQuest(1)
	require.Error(t, err)
	assert.Nil(t, reward)
	var validationErr *domain.ValidationError
	assert.ErrorAs(t, err, &validationErr)
}
