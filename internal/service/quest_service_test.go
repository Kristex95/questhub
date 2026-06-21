package service

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/Kristex95/questhub/internal/domain"
)

// Mocks
// Quest Mock
type MockQuestRepository struct {
	data      map[int]domain.Quest
	counter   int
	CreateErr error
	GetErr    error
	GetAllErr error
	UpdateErr error
	DeleteErr error
}

func NewMockQuestRepository() *MockQuestRepository {
	return &MockQuestRepository{data: make(map[int]domain.Quest)}
}

func (m *MockQuestRepository) Create(quest domain.Quest) (domain.Quest, error) {
	if m.CreateErr != nil {
		return domain.Quest{}, m.CreateErr
	}
	m.counter++
	quest.ID = m.counter
	m.data[quest.ID] = quest
	return quest, nil
}

func (m *MockQuestRepository) Get(id int) (domain.Quest, error) {
	if m.GetErr != nil {
		return domain.Quest{}, m.GetErr
	}
	quest, ok := m.data[id]
	if !ok {
		return domain.Quest{}, &domain.NotFoundError{Entity: "Quest", Value: strconv.Itoa(id)}
	}
	return quest, nil
}

func (m *MockQuestRepository) GetAll() ([]domain.Quest, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	quests := make([]domain.Quest, 0, len(m.data))
	for _, q := range m.data {
		quests = append(quests, q)
	}
	return quests, nil
}

func (m *MockQuestRepository) Update(questID int, updatedQuest domain.Quest) (domain.Quest, error) {
	if m.UpdateErr != nil {
		return domain.Quest{}, m.UpdateErr
	}
	if _, ok := m.data[questID]; !ok {
		return domain.Quest{}, &domain.NotFoundError{Entity: "Quest", Value: strconv.Itoa(questID)}
	}
	updatedQuest.ID = questID
	m.data[questID] = updatedQuest
	return updatedQuest, nil
}

func (m *MockQuestRepository) Delete(questID int) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.data[questID]; !ok {
		return &domain.NotFoundError{Entity: "Quest", Value: strconv.Itoa(questID)}
	}
	delete(m.data, questID)
	return nil
}

// Task Mock
type MockTaskRepository struct {
	data          map[int]domain.Task
	counter       int
	CreateErr     error
	GetErr        error
	GetAllErr     error
	GetByQuestErr error
	UpdateErr     error
	DeleteErr     error
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{data: make(map[int]domain.Task)}
}

func (m *MockTaskRepository) Create(task domain.Task) (domain.Task, error) {
	if m.CreateErr != nil {
		return domain.Task{}, m.CreateErr
	}
	m.counter++
	task.ID = m.counter
	m.data[task.ID] = task
	return task, nil
}

func (m *MockTaskRepository) Get(id int) (domain.Task, error) {
	if m.GetErr != nil {
		return domain.Task{}, m.GetErr
	}
	task, ok := m.data[id]
	if !ok {
		return domain.Task{}, &domain.NotFoundError{Entity: "Task", Value: strconv.Itoa(id)}
	}
	return task, nil
}

func (m *MockTaskRepository) GetAll() ([]domain.Task, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	tasks := make([]domain.Task, 0, len(m.data))
	for _, t := range m.data {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (m *MockTaskRepository) GetByQuestID(questID int) ([]*domain.Task, error) {
	if m.GetByQuestErr != nil {
		return nil, m.GetByQuestErr
	}
	var tasks []*domain.Task
	for _, t := range m.data {
		if t.QuestId == questID {
			taskCopy := t
			tasks = append(tasks, &taskCopy)
		}
	}
	return tasks, nil
}

func (m *MockTaskRepository) Update(taskID int, updatedTask domain.Task) (domain.Task, error) {
	if m.UpdateErr != nil {
		return domain.Task{}, m.UpdateErr
	}
	if _, ok := m.data[taskID]; !ok {
		return domain.Task{}, &domain.NotFoundError{Entity: "Task", Value: strconv.Itoa(taskID)}
	}
	updatedTask.ID = taskID
	m.data[taskID] = updatedTask
	return updatedTask, nil
}

func (m *MockTaskRepository) Delete(taskID int) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.data[taskID]; !ok {
		return &domain.NotFoundError{Entity: "Task", Value: strconv.Itoa(taskID)}
	}
	delete(m.data, taskID)
	return nil
}

// User Mock
type MockUserRepository struct {
	data             map[int]domain.User
	counter          int
	CreateErr        error
	GetErr           error
	GetAllErr        error
	GetByUsernameErr error
	UpdateErr        error
	DeleteErr        error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{data: make(map[int]domain.User)}
}

func (m *MockUserRepository) Create(user domain.User) (domain.User, error) {
	if m.CreateErr != nil {
		return domain.User{}, m.CreateErr
	}
	m.counter++
	user.ID = m.counter
	m.data[user.ID] = user
	return user, nil
}

func (m *MockUserRepository) Get(id int) (domain.User, error) {
	if m.GetErr != nil {
		return domain.User{}, m.GetErr
	}
	user, ok := m.data[id]
	if !ok {
		return domain.User{}, &domain.NotFoundError{Entity: "User", Value: strconv.Itoa(id)}
	}
	return user, nil
}

func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	users := make([]domain.User, 0, len(m.data))
	for _, u := range m.data {
		users = append(users, u)
	}
	return users, nil
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	if m.GetByUsernameErr != nil {
		return nil, m.GetByUsernameErr
	}
	for _, u := range m.data {
		if u.Username == username {
			userCopy := u
			return &userCopy, nil
		}
	}
	return nil, &domain.NotFoundError{Entity: "User"}
}

func (m *MockUserRepository) Update(userID int, updatedUser domain.User) (domain.User, error) {
	if m.UpdateErr != nil {
		return domain.User{}, m.UpdateErr
	}
	if _, ok := m.data[userID]; !ok {
		return domain.User{}, &domain.NotFoundError{Entity: "User", Value: strconv.Itoa(userID)}
	}
	updatedUser.ID = userID
	m.data[userID] = updatedUser
	return updatedUser, nil
}

func (m *MockUserRepository) Delete(userID int) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.data[userID]; !ok {
		return &domain.NotFoundError{Entity: "User", Value: strconv.Itoa(userID)}
	}
	delete(m.data, userID)
	return nil
}

func TestCreateQuest_Success(t *testing.T) {
	svc, _, _, _ := makeService()

	q, err := svc.CreateQuest("Epic Journey", "Long description of the quest", 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if q == nil {
		t.Fatal("expected quest not to be nil")
	}
}

func TestCreateQuest_ValidationError(t *testing.T) {
	svc, _, _, _ := makeService()

	_, err := svc.CreateQuest("ab", "Valid description", 5)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected error to wrap ErrValidation, got %v", err)
	}
}

func TestCreateQuest_RepoError(t *testing.T) {
	svc, questRepo, _, _ := makeService()

	mockDbErr := errors.New("db connection timeout")
	questRepo.CreateErr = mockDbErr

	_, err := svc.CreateQuest("Epic Journey", "Description", 5)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, mockDbErr) {
		t.Errorf("expected error to wrap repo error, got %v", err)
	}
}

func TestAddTaskToQuest_Success(t *testing.T) {
	svc, questRepo, _, _ := makeService()

	q, _ := questRepo.Create(domain.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})

	task, err := svc.AddTaskToQuest(q.ID, "Task 1", "Kill 5 wolves")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task == nil {
		t.Fatal("expected task not to be nil")
	}
	if task.QuestId != q.ID {
		t.Errorf("expected task to be linked to quest %d, got %d", q.ID, task.QuestId)
	}
}

func TestAddTaskToQuest_QuestNotFound(t *testing.T) {
	svc, _, _, _ := makeService()

	_, err := svc.AddTaskToQuest(999, "Task 1", "Desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestGetQuestTasks_Success(t *testing.T) {
	svc, questRepo, taskRepo, _ := makeService()

	q1, _ := questRepo.Create(domain.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})
	q2, _ := questRepo.Create(domain.Quest{Title: "Q2", Description: "Desc", Difficulty: 1})

	_, _ = taskRepo.Create(domain.Task{QuestId: q1.ID, Title: "Task 1"})
	_, _ = taskRepo.Create(domain.Task{QuestId: q2.ID, Title: "Task 2"})

	tasks, err := svc.GetQuestTasks(q1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].QuestId != q1.ID {
		t.Errorf("expected task to belong to quest %d, got %d", q1.ID, tasks[0].QuestId)
	}
}

func TestCompleteQuest_Success(t *testing.T) {
	svc, questRepo, _, userRepo := makeService()

	difficulty := 3
	q, _ := questRepo.Create(domain.Quest{Title: "Q1", Description: "Desc", Difficulty: difficulty})
	u, _ := userRepo.Create(domain.User{Username: "Player1"})

	err := svc.CompleteQuest(q.ID, u.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedUser, _ := userRepo.Get(u.ID)
	expectedXP := difficulty * 100
	if updatedUser.XP != expectedXP {
		t.Errorf("expected user XP to be %d, got %d", expectedXP, updatedUser.XP)
	}
}

func TestCompleteQuest_IncompleteTasks(t *testing.T) {
	svc, questRepo, taskRepo, userRepo := makeService()

	q, _ := questRepo.Create(domain.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})
	u, _ := userRepo.Create(domain.User{Username: "Player1"})

	_, _ = taskRepo.Create(domain.Task{QuestId: q.ID, Title: "Incomplete Task"})

	err := svc.CompleteQuest(q.ID, u.ID)
	if err == nil {
		t.Fatalf("expected error due to incomplete tasks, got nil")
	}

	if !strings.Contains(err.Error(), "validation error") {
		t.Errorf("expected validation error message, got: %v", err)
	}
}

func TestCompleteQuest_QuestNotFound(t *testing.T) {
	svc, _, _, userRepo := makeService()

	u, _ := userRepo.Create(domain.User{Username: "Player1"})

	err := svc.CompleteQuest(999, u.ID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestCompleteQuest_UserNotFound(t *testing.T) {
	svc, questRepo, _, _ := makeService()

	q, _ := questRepo.Create(domain.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})

	err := svc.CompleteQuest(q.ID, 999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestDeleteQuest_DeletesTasks(t *testing.T) {
	svc, questRepo, taskRepo, _ := makeService()

	q, _ := questRepo.Create(domain.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})
	_, _ = taskRepo.Create(domain.Task{QuestId: q.ID, Title: "Task 1"})

	err := svc.DeleteQuest(q.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tasks, err := taskRepo.GetByQuestID(q.ID)
	if err != nil {
		t.Fatalf("unexpected error from repository: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected tasks to be fully deleted, but found %d remaining tasks", len(tasks))
	}
}

func TestDeleteQuest_QuestNotFound(t *testing.T) {
	svc, _, _, _ := makeService()

	err := svc.DeleteQuest(999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestGetQuest_Success(t *testing.T) {
	svc, questRepo, _, _ := makeService()
	q, _ := questRepo.Create(domain.Quest{Title: "Find the Holy Grail", Description: "Desc", Difficulty: 7})

	got, err := svc.GetQuest(q.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Title != "Find the Holy Grail" {
		t.Errorf("expected title 'Find the Holy Grail', got '%s'", got.Title)
	}
}

func TestGetQuest_RepoError(t *testing.T) {
	svc, questRepo, _, _ := makeService()
	questRepo.GetErr = errors.New("read timeout")

	_, err := svc.GetQuest(999)
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
	if !strings.Contains(err.Error(), "get quest") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestListQuests_Success(t *testing.T) {
	svc, questRepo, _, _ := makeService()
	_, _ = questRepo.Create(domain.Quest{Title: "Q1", Description: "D1", Difficulty: 1})
	_, _ = questRepo.Create(domain.Quest{Title: "Q2", Description: "D2", Difficulty: 2})

	list, err := svc.ListQuests()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 quests, got %d", len(list))
	}
}

func TestListQuests_RepoError(t *testing.T) {
	svc, questRepo, _, _ := makeService()
	questRepo.GetAllErr = errors.New("cluster down")

	_, err := svc.ListQuests()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get all quests") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func makeService() (*QuestService, *MockQuestRepository, *MockTaskRepository, *MockUserRepository) {

	questRepo := NewMockQuestRepository()
	taskRepo := NewMockTaskRepository()
	userRepo := NewMockUserRepository()
	questService := NewQuestService(questRepo, taskRepo, userRepo)
	return questService, questRepo, taskRepo, userRepo

}
