package main

import (
	"errors"
	"fmt"

	"github.com/Kristex95/questhub/internal/domain"
	"github.com/Kristex95/questhub/internal/repository"
	"github.com/Kristex95/questhub/internal/service"
)

func main() {
	questRepository := repository.NewInMemoryQuestRepository()
	taskRepository := repository.NewInMemoryTaskRepository()
	userRepository := repository.NewInMemoryUserRepository()
	questService := service.NewQuestService(questRepository, taskRepository, userRepository)

	// === CREATE USERS ===
	fmt.Println("\n[1] Creating users")

	user1, err := domain.NewUser("u1", "Ezio Auditore", "ezio@assassin.com")
	if err != nil {
		fmt.Println("User1 creation error:", err)
	}

	user2, err := domain.NewUser("u2", "Altair Ibn-La'Ahad", "altair@assassin.com")
	if err != nil {
		fmt.Println("User2 creation error:", err)
	}

	fmt.Println("Users created:")
	fmt.Println(" -", user1.Username)
	fmt.Println(" -", user2.Username)

	// === CREATE QUESTS ===
	fmt.Println("\n[2] Creating quests")

	q1, err := questService.CreateQuest("Training", "Basic training", 2)
	if err != nil {
		fmt.Printf("Quest1 creating: %v\n", err)
	}

	q2, err := questService.CreateQuest("Stealth Ops", "Advanced infiltration", 5)
	if err != nil {
		fmt.Printf("Quest2 creating: %v\n", err)
	}

	// Tasks for Q1
	task1, err := domain.NewTask("t1", "Run 1km", 50, q1.ID)
	if err != nil {
		fmt.Printf("Task1 creating: %v\n", err)
	}
	*task1, err = taskRepository.Create(*task1)
	if err != nil {
		fmt.Printf("Task1 creating in repo: %v\n", err)
	}
	q1.AddTask(task1)

	task2, err := domain.NewTask("t2", "Climb tower", 70, q1.ID)
	if err != nil {
		fmt.Printf("Task2 creating: %v\n", err)
	}
	*task2, err = taskRepository.Create(*task2)
	if err != nil {
		fmt.Printf("Task2 creating in repo: %v\n", err)
	}
	q1.AddTask(task2)

	// Tasks for Q2
	task3, err := domain.NewTask("t3", "Scout area", 80, q2.ID)
	if err != nil {
		fmt.Printf("Task3 creating: %v\n", err)
	}
	*task3, err = taskRepository.Create(*task3)
	if err != nil {
		fmt.Printf("Task3 creating in repo: %v\n", err)
	}
	q2.AddTask(task3)

	task4, err := domain.NewTask("t4", "Avoid guards", 120, q2.ID)
	if err != nil {
		fmt.Printf("Task4 creating: %v\n", err)
	}
	*task4, err = taskRepository.Create(*task4)
	if err != nil {
		fmt.Printf("Task4 creating in repo: %v\n", err)
	}
	q2.AddTask(task4)

	task5, err := domain.NewTask("t5", "Steal document", 150, q2.ID)
	if err != nil {
		fmt.Printf("Task5 creating: %v\n", err)
	}
	*task5, err = taskRepository.Create(*task5)
	if err != nil {
		fmt.Printf("Task5 creating in repo: %v\n", err)
	}
	q2.AddTask(task5)

	fmt.Println("Quests created with tasks")

	// === SIMULATE INCOMPLETE COMPLETION ERROR ===
	fmt.Println("\n[4] Attempting to complete quest with incomplete tasks")

	if err := user1.CompleteQuest(q1); err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// === MARK TASKS AS COMPLETED ===
	fmt.Println("\n[5] Marking all tasks in quest q1 as completed")

	for _, task := range q1.Tasks {
		q1.CompleteTask(task.ID)

		updatedTask, err := taskRepository.Update(task.ID, *task)
		if err != nil {
			fmt.Printf("Failed to update task %s: %v\n", task.ID, err)
			continue
		}

		fmt.Printf("Task '%s' marked as completed\n", updatedTask.Title)
	}

	// === COMPLETE QUEST SUCCESSFULLY ===
	fmt.Println("\n[6] Completing quest after all tasks are done")

	beforeXP := user1.XP

	if err := user1.CompleteQuest(q1); err != nil {
		fmt.Printf("Unexpected completion error: %v\n", err)
	} else {
		fmt.Printf("Quest '%s' completed successfully\n", q1.Title)
		fmt.Printf("XP before: %d | XP after: %d | Gained: %d\n",
			beforeXP,
			user1.XP,
			user1.XP-beforeXP,
		)
	}

	// === INVALID QUEST CREATION ===
	fmt.Println("\n[7] Creating invalid quest and inspecting error chain")

	_, err = domain.NewQuest("bad", "", "", -1)
	if err != nil {
		fmt.Printf("Returned error: %v\n", err)

		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf(
				"Validation error: field=%s, message=%s\n",
				validationErr.Field,
				validationErr.Message,
			)
		}
	}

	// === NOT FOUND EXAMPLE ===
	fmt.Println("\n[8] Loading non-existing quest")

	_, err = questRepository.Get("quest-does-not-exist")
	if err != nil {
		var notFoundError *domain.NotFoundError
		if errors.As(err, &notFoundError) {
			fmt.Printf(
				"NotFound error: entity=%s, id=%s\n",
				notFoundError.Entity,
				notFoundError.ID,
			)
		}
	}

	// === DELETE QUEST AND VERIFY TASKS ===
	fmt.Println("\n[9] Deleting quest q2 and verifying related tasks")

	if err := questService.DeleteQuest(q2.ID); err != nil {
		fmt.Printf("Delete failed: %v\n", err)
	} else {
		fmt.Printf("Quest '%s' deleted\n", q2.Title)
	}

	for _, task := range q2.Tasks {
		_, err := taskRepository.Get(task.ID)
		if err != nil {
			fmt.Printf("Task %s is no longer available: %v\n", task.ID, err)
		} else {
			fmt.Printf("WARNING: task %s still exists\n", task.ID)
		}
	}

	fmt.Println("\nDemo completed")

}
