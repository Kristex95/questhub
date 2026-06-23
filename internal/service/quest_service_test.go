package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"testing"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCtx = context.Background()

// Mocks
// Quest Mock
type MockQuestRepository struct {
	data      map[int64]*models.Quest
	counter   int64
	CreateErr error
	GetErr    error
	GetAllErr error
	UpdateErr error
	DeleteErr error
}

func NewMockQuestRepository() *MockQuestRepository {
	return &MockQuestRepository{data: make(map[int64]*models.Quest)}
}

func (m *MockQuestRepository) Create(ctx context.Context, quest *models.Quest) (*models.Quest, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	m.counter++
	quest.ID = m.counter
	m.data[quest.ID] = quest
	return quest, nil
}

func (m *MockQuestRepository) GetByID(ctx context.Context, id int64) (*models.Quest, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	quest, ok := m.data[id]
	if !ok {
		return nil, &domain.NotFoundError{Entity: "Quest", Value: strconv.FormatInt(id, 10)}
	}
	return quest, nil
}

func (m *MockQuestRepository) GetAll(ctx context.Context) ([]*models.Quest, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	quests := make([]*models.Quest, 0, len(m.data))
	for _, q := range m.data {
		quests = append(quests, q)
	}
	return quests, nil
}

func (m *MockQuestRepository) Update(ctx context.Context, quest *models.Quest) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if _, ok := m.data[quest.ID]; !ok {
		return &domain.NotFoundError{Entity: "Quest", Value: strconv.FormatInt(quest.ID, 10)}
	}
	m.data[quest.ID] = quest
	return nil
}

func (m *MockQuestRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.data[id]; !ok {
		return &domain.NotFoundError{Entity: "Quest", Value: strconv.FormatInt(id, 10)}
	}
	delete(m.data, id)
	return nil
}

// Task Mock
type MockTaskRepository struct {
	data          map[int64]*models.Task
	counter       int64
	CreateErr     error
	GetErr        error
	GetAllErr     error
	GetByQuestErr error
	UpdateErr     error
	DeleteErr     error
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{data: make(map[int64]*models.Task)}
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	m.counter++
	task.ID = m.counter
	m.data[task.ID] = task
	return task, nil
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	task, ok := m.data[id]
	if !ok {
		return nil, &domain.NotFoundError{Entity: "Task", Value: strconv.FormatInt(id, 10)}
	}
	return task, nil
}

func (m *MockTaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	tasks := make([]*models.Task, 0, len(m.data))
	for _, t := range m.data {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (m *MockTaskRepository) GetByQuestID(ctx context.Context, questID int64) ([]*models.Task, error) {
	if m.GetByQuestErr != nil {
		return nil, m.GetByQuestErr
	}
	var tasks []*models.Task
	for _, t := range m.data {
		if t.QuestID == questID {
			taskCopy := *t
			tasks = append(tasks, &taskCopy)
		}
	}
	return tasks, nil
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if _, ok := m.data[task.ID]; !ok {
		return &domain.NotFoundError{Entity: "Task", Value: strconv.FormatInt(task.ID, 10)}
	}
	m.data[task.ID] = task
	return nil
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.data[id]; !ok {
		return &domain.NotFoundError{Entity: "Task", Value: strconv.FormatInt(id, 10)}
	}
	delete(m.data, id)
	return nil
}

// User Mock
type MockUserRepository struct {
	data             map[int64]*models.User
	counter          int64
	CreateErr        error
	GetErr           error
	GetAllErr        error
	GetByUsernameErr error
	UpdateErr        error
	DeleteErr        error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{data: make(map[int64]*models.User)}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	m.counter++
	user.ID = int(m.counter)
	m.data[m.counter] = user
	return user, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	user, ok := m.data[id]
	if !ok {
		return nil, &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(id, 10)}
	}
	return user, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	if m.GetAllErr != nil {
		return nil, m.GetAllErr
	}
	users := make([]*models.User, 0, len(m.data))
	for _, u := range m.data {
		users = append(users, u)
	}
	return users, nil
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	if m.GetByUsernameErr != nil {
		return nil, m.GetByUsernameErr
	}
	for _, u := range m.data {
		if u.Username == username {
			userCopy := *u
			return &userCopy, nil
		}
	}
	return nil, &domain.NotFoundError{Entity: "User"}
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	id := int64(user.ID)
	if _, ok := m.data[id]; !ok {
		return &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(id, 10)}
	}
	m.data[id] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.data[id]; !ok {
		return &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(id, 10)}
	}
	delete(m.data, id)
	return nil
}

func (m *MockUserRepository) AddXP(ctx context.Context, userID int64, amount int) error {
	user, ok := m.data[userID]
	if !ok {
		return &domain.NotFoundError{Entity: "User", Value: strconv.FormatInt(userID, 10)}
	}
	user.XP += amount
	return nil
}

type MockRewardRepository struct {
	data          map[int64][]*models.Reward
	counter       int64
	CreateErr     error
	GetByQuestErr error
}

func NewMockRewardRepository() *MockRewardRepository {
	return &MockRewardRepository{data: make(map[int64][]*models.Reward)}
}

func (m *MockRewardRepository) Create(ctx context.Context, reward *models.Reward) (*models.Reward, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	m.counter++
	reward.ID = m.counter
	m.data[reward.QuestID] = append(m.data[reward.QuestID], reward)
	return reward, nil
}

func (m *MockRewardRepository) GetByQuestID(ctx context.Context, questID int64) ([]*models.Reward, error) {
	if m.GetByQuestErr != nil {
		return nil, m.GetByQuestErr
	}
	return m.data[questID], nil
}

type MockProgressRepository struct {
	data             map[string]*models.Progress
	CreateErr        error
	GetErr           error
	MarkCompletedErr error
}

func NewMockProgressRepository() *MockProgressRepository {
	return &MockProgressRepository{data: make(map[string]*models.Progress)}
}

func (m *MockProgressRepository) Create(ctx context.Context, progress *models.Progress) (*models.Progress, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	key := fmt.Sprintf("%d:%d", progress.UserID, progress.QuestID)
	m.data[key] = progress
	return progress, nil
}

func (m *MockProgressRepository) GetByUserAndQuest(ctx context.Context, userID, questID int64) (*models.Progress, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	key := fmt.Sprintf("%d:%d", userID, questID)
	progress, ok := m.data[key]
	if !ok {
		return nil, &domain.NotFoundError{Entity: "Progress"}
	}
	return progress, nil
}

func (m *MockProgressRepository) MarkCompleted(ctx context.Context, userID, questID int64) error {
	if m.MarkCompletedErr != nil {
		return m.MarkCompletedErr
	}
	key := fmt.Sprintf("%d:%d", userID, questID)
	progress, ok := m.data[key]
	if !ok {
		return &domain.NotFoundError{Entity: "Progress"}
	}
	progress.Status = "completed"
	return nil
}

type MockLeaderboardUpdater struct{}

func (m *MockLeaderboardUpdater) UpdateLeaderboard(ctx context.Context, userID int64, xp int) error {
	return nil
}

type MockStatsIncrementer struct{}

func (m *MockStatsIncrementer) IncrCompletedQuests(ctx context.Context) error {
	return nil
}

type MockNotifier struct{}


func (m *MockNotifier) Notify(ctx context.Context, userID int64, message string) error {
	return nil
}

func TestCreateQuest_Success(t *testing.T) {
	svc, _, _, _ := makeService()

	q, err := svc.CreateQuest(testCtx, "Epic Journey", "Long description of the quest", 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if q == nil {
		t.Fatal("expected quest not to be nil")
	}
}

func TestCreateQuest_ValidationError(t *testing.T) {
	svc, _, _, _ := makeService()

	_, err := svc.CreateQuest(testCtx, "ab", "Valid description", 5)
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

	_, err := svc.CreateQuest(testCtx, "Epic Journey", "Description", 5)
	require.Error(t, err)
	assert.True(t, errors.Is(err, mockDbErr), "expected error to wrap repo error, got %v", err)
}

func TestCreateQuest_Validation(t *testing.T) {
	svc, _, _, _ := makeService()

	cases := []struct {
		name        string
		title       string
		description string
		difficulty  int
		wantErr     bool
	}{
		{"too short title", "ab", "valid desc", 5, true},
		{"empty description", "Valid title", "", 5, true},
		{"invalid difficulty", "Valid title", "Valid desc", 11, true},
		{"valid quest", "Valid title", "Valid desc", 5, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := svc.CreateQuest(testCtx, tc.title, tc.description, tc.difficulty)
			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, q)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, q)
			assert.Equal(t, tc.title, q.Title)
			assert.Equal(t, tc.description, q.Description)
			assert.Equal(t, tc.difficulty, q.Difficulty)
		})
	}
}

func TestAddTaskToQuest_Success(t *testing.T) {
	svc, questRepo, _, _ := makeService()

	q, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})

	task, err := svc.AddTaskToQuest(testCtx, q.ID, "Task 1", "Kill 5 wolves")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task == nil {
		t.Fatal("expected task not to be nil")
	}
	if task.QuestID != q.ID {
		t.Errorf("expected task to be linked to quest %d, got %d", q.ID, task.QuestID)
	}
}

func TestAddTaskToQuest_QuestNotFound(t *testing.T) {
	svc, _, _, _ := makeService()

	_, err := svc.AddTaskToQuest(testCtx, 999, "Task 1", "Desc")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestGetQuestTasks_Success(t *testing.T) {
	svc, questRepo, taskRepo, _ := makeService()

	q1, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})
	q2, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q2", Description: "Desc", Difficulty: 1})

	_, _ = taskRepo.Create(testCtx, &models.Task{QuestID: q1.ID, Title: "Task 1"})
	_, _ = taskRepo.Create(testCtx, &models.Task{QuestID: q2.ID, Title: "Task 2"})

	tasks, err := svc.GetQuestTasks(testCtx, q1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].QuestID != q1.ID {
		t.Errorf("expected task to belong to quest %d, got %d", q1.ID, tasks[0].QuestID)
	}
}

func TestCompleteQuest_Success(t *testing.T) {
	svc, questRepo, taskRepo, userRepo := makeService()

	q, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "Desc", Difficulty: 3, XPReward: 300})
	u, _ := userRepo.Create(testCtx, &models.User{Username: "Player1"})

	createdTask, _ := taskRepo.Create(testCtx, &models.Task{QuestID: q.ID, Title: "Task 1", IsCompleted: true})
	if createdTask == nil {
		t.Fatal("expected task to be created")
	}

	_, err := svc.progress.StartProgress(testCtx, int64(u.ID), q.ID)
	if err != nil {
		t.Fatalf("expected progress to start, got %v", err)
	}

	err = svc.CompleteQuest(testCtx, int64(u.ID), q.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedUser, _ := userRepo.GetByID(testCtx, int64(u.ID))
	expectedXP := 300
	if updatedUser.XP != expectedXP {
		t.Errorf("expected user XP to be %d, got %d", expectedXP, updatedUser.XP)
	}
}

func TestCompleteQuest_IncompleteTasks(t *testing.T) {
	svc, questRepo, taskRepo, userRepo := makeService()

	q, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})
	u, _ := userRepo.Create(testCtx, &models.User{Username: "Player1"})

	_, _ = taskRepo.Create(testCtx, &models.Task{QuestID: q.ID, Title: "Incomplete Task"})

	err := svc.CompleteQuest(testCtx, int64(u.ID), q.ID)
	if err == nil {
		t.Fatalf("expected error due to incomplete tasks, got nil")
	}

	if !strings.Contains(err.Error(), "not all quest tasks are completed") {
		t.Errorf("expected validation error message, got: %v", err)
	}
}

func TestCompleteQuest_QuestNotFound(t *testing.T) {
	svc, _, _, userRepo := makeService()

	u, _ := userRepo.Create(testCtx, &models.User{Username: "Player1"})

	err := svc.CompleteQuest(testCtx, int64(u.ID), 999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestCompleteQuest_UserNotFound(t *testing.T) {
	svc, questRepo, _, _ := makeService()

	q, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})

	err := svc.CompleteQuest(testCtx, 999, q.ID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestDeleteQuest_DeletesTasks(t *testing.T) {
	svc, questRepo, taskRepo, _ := makeService()

	q, _ := questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "Desc", Difficulty: 1})
	_, _ = taskRepo.Create(testCtx, &models.Task{QuestID: q.ID, Title: "Task 1"})

	err := svc.DeleteQuest(testCtx, q.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tasks, err := taskRepo.GetByQuestID(testCtx, q.ID)
	if err != nil {
		t.Fatalf("unexpected error from repository: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected tasks to be fully deleted, but found %d remaining tasks", len(tasks))
	}
}

func TestDeleteQuest_QuestNotFound(t *testing.T) {
	svc, _, _, _ := makeService()

	err := svc.DeleteQuest(testCtx, 999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got %v", err)
	}
}

func TestGetQuest_Success(t *testing.T) {
	svc, questRepo, _, _ := makeService()
	q, _ := questRepo.Create(testCtx, &models.Quest{Title: "Find the Holy Grail", Description: "Desc", Difficulty: 7})

	got, err := svc.GetQuest(testCtx, q.ID)
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

	_, err := svc.GetQuest(testCtx, 999)
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
	if !strings.Contains(err.Error(), "get quest") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestListQuests_Success(t *testing.T) {
	svc, questRepo, _, _ := makeService()
	_, _ = questRepo.Create(testCtx, &models.Quest{Title: "Q1", Description: "D1", Difficulty: 1})
	_, _ = questRepo.Create(testCtx, &models.Quest{Title: "Q2", Description: "D2", Difficulty: 2})

	list, err := svc.ListQuests(testCtx)
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

	_, err := svc.ListQuests(testCtx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "list quests") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func makeService() (*QuestService, *MockQuestRepository, *MockTaskRepository, *MockUserRepository) {
	questRepo := NewMockQuestRepository()
	taskRepo := NewMockTaskRepository()
	userRepo := NewMockUserRepository()
	progressRepo := NewMockProgressRepository()
	rewardRepo := NewMockRewardRepository()

	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	questService := NewQuestService(
		questRepo,
		taskRepo,
		NewRewardService(rewardRepo, userRepo, &MockLeaderboardUpdater{}),
		NewProgressService(progressRepo),
		&MockStatsIncrementer{},
		&MockNotifier{},
		discardLogger,
	)
	return questService, questRepo, taskRepo, userRepo
}
