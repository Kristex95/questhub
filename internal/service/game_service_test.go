package service

import (
	"errors"
	"testing"

	"github.com/Kristex95/questhub/internal/domain"
)

func setupGameTest() (*GameService, *MockQuestRepository, *MockTaskRepository, *MockUserRepository) {
	questRepo := NewMockQuestRepository()
	taskRepo := NewMockTaskRepository()
	userRepo := NewMockUserRepository()
	svc := NewGameService(questRepo, taskRepo, userRepo)
	return svc, questRepo, taskRepo, userRepo
}

func TestGameService_HappyPath(t *testing.T) {
	svc, questRepo, taskRepo, userRepo := setupGameTest()

	user, _ := userRepo.Create(domain.User{Username: "Hero", Email: "hero@test.com"})
	quest, _ := questRepo.Create(domain.Quest{Title: "Save the Village", Description: "Defeat local monsters", Difficulty: 4})
	
	t1, _ := taskRepo.Create(domain.Task{QuestId: quest.ID, Title: "Gather herbs"})
	t2, _ := taskRepo.Create(domain.Task{QuestId: quest.ID, Title: "Defeat Goblin Boss"})

	err := svc.StartQuest(user.ID, quest.ID)
	if err != nil {
		t.Fatalf("expected no error on StartQuest, got %v", err)
	}

	err = svc.CompleteTask(user.ID, t1.ID)
	if err != nil {
		t.Fatalf("expected no error on CompleteTask 1, got %v", err)
	}

	progress, err := svc.GetProgress(user.ID)
	if err != nil {
		t.Fatalf("expected no error on GetProgress, got %v", err)
	}
	if progress.Percentage != 50.0 {
		t.Errorf("expected progress percentage to be 50.0, got %.1f", progress.Percentage)
	}

	err = svc.CompleteTask(user.ID, t2.ID)
	if err != nil {
		t.Fatalf("expected no error on CompleteTask 2, got %v", err)
	}

	progress, err = svc.GetProgress(user.ID)
	if err != nil {
		t.Fatalf("expected no error on GetProgress, got %v", err)
	}
	if progress.Percentage != 100.0 {
		t.Errorf("expected progress percentage to be 100.0, got %.1f", progress.Percentage)
	}
	if progress.CompletedTasks != 2 || progress.TotalTasks != 2 {
		t.Errorf("expected 2/2 tasks, got %d/%d", progress.CompletedTasks, progress.TotalTasks)
	}

	reward, err := svc.FinishQuest(user.ID)
	if err != nil {
		t.Fatalf("expected no error on FinishQuest, got %v", err)
	}

	expectedRewardTitle := "Completed: Save the Village"
	if reward.Title != expectedRewardTitle {
		t.Errorf("expected reward title '%s', got '%s'", expectedRewardTitle, reward.Title)
	}
	expectedXP := 4 * 100 // difficulty * 100
	if reward.XPAmount != expectedXP {
		t.Errorf("expected reward XP %d, got %d", expectedXP, reward.XPAmount)
	}

	updatedUser, _ := userRepo.Get(user.ID)
	if updatedUser.XP != expectedXP {
		t.Errorf("expected user to have %d XP, got %d", expectedXP, updatedUser.XP)
	}

	_, err = svc.GetProgress(user.ID)
	if err == nil {
		t.Error("expected error when getting progress after finishing the quest, but got nil")
	}
}

func TestGameService_Start_NoActiveQuest(t *testing.T) {
	svc, questRepo, _, userRepo := setupGameTest()

	user, _ := userRepo.Create(domain.User{Username: "Player"})
	quest, _ := questRepo.Create(domain.Quest{Title: "Solo Leveling", Difficulty: 1})

	err := svc.StartQuest(user.ID, quest.ID)
	if err != nil {
		t.Errorf("expected single quest start to succeed, got %v", err)
	}
}

// 3. Start коли вже є активний квест - помилка
func TestGameService_Start_AlreadyHasActiveQuest(t *testing.T) {
	svc, questRepo, _, userRepo := setupGameTest()

	u, _ := userRepo.Create(domain.User{Username: "BusyPlayer"})
	q1, _ := questRepo.Create(domain.Quest{Title: "First Quest", Difficulty: 1})
	q2, _ := questRepo.Create(domain.Quest{Title: "Second Quest", Difficulty: 2})

	_ = svc.StartQuest(u.ID, q1.ID)

	err := svc.StartQuest(u.ID, q2.ID)
	if err == nil {
		t.Fatal("expected error when starting a second quest, got nil")
	}

	if !errors.Is(err, domain.ErrAlreadyStarted) {
		t.Errorf("expected error to be ErrAlreadyStarted, got %v", err)
	}
}

func TestGameService_CompleteTask_WrongQuest(t *testing.T) {
	svc, questRepo, taskRepo, userRepo := setupGameTest()

	user, _ := userRepo.Create(domain.User{Username: "Cheater"})
	q1, _ := questRepo.Create(domain.Quest{Title: "Active Quest", Difficulty: 1})
	q2, _ := questRepo.Create(domain.Quest{Title: "Other Quest", Difficulty: 1})

	_ = svc.StartQuest(user.ID, q1.ID) 

	foreignTask, _ := taskRepo.Create(domain.Task{QuestId: q2.ID, Title: "Foreign Objective"})

	err := svc.CompleteTask(user.ID, foreignTask.ID)
	if err == nil {
		t.Fatal("expected error when completing a task from an inactive quest, got nil")
	}

	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ErrValidation for out-of-scope task, got %v", err)
	}
}

func TestGameService_FinishQuest_IncompleteTasks(t *testing.T) {
	svc, questRepo, taskRepo, userRepo := setupGameTest()

	user, _ := userRepo.Create(domain.User{Username: "Slacker"})
	quest, _ := questRepo.Create(domain.Quest{Title: "Hard Quest", Difficulty: 5})
	
	_, _ = taskRepo.Create(domain.Task{QuestId: quest.ID, Title: "Do some actual work"})

	_ = svc.StartQuest(user.ID, quest.ID)

	_, err := svc.FinishQuest(user.ID)
	if err == nil {
		t.Fatal("expected error when finishing quest with incomplete tasks, got nil")
	}

	var valErr *domain.ValidationError
	if !errors.As(err, &valErr) {
		t.Errorf("expected domain.ValidationError, got type %T (%v)", err, err)
	}
}

func TestGameService_GetProgress_NoActiveQuest(t *testing.T) {
	svc, _, _, userRepo := setupGameTest()

	u, _ := userRepo.Create(domain.User{Username: "LoverOfIdleness"})

	_, err := svc.GetProgress(u.ID)
	if err == nil {
		t.Fatal("expected error when getting progress without an active quest, got nil")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound when no quest is active, got %v", err)
	}
}