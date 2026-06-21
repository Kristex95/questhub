package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Kristex95/questhub/internal/config"
	"github.com/Kristex95/questhub/internal/database"
	"github.com/Kristex95/questhub/internal/models"
	"github.com/Kristex95/questhub/internal/repository/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)	
	}

	pool, err := database.NewPostgresPool(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()
	log.Println("postgres connected")

	userRepo := postgres.NewUserRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	questRepo := postgres.NewQuestRepository(pool)
	rewardRepo := postgres.NewRewardRepository(pool)
	progressRepo := postgres.NewProgressRepository(pool)

	user := &models.User{
		Username: "Ezio Auditore",
		Email:    "ezio@assassin.com",
		XP:       0,
		Level:    1,
	}
	if _, err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("create user: %v", err)
	}
	log.Printf("user created with ID %d", user.ID)

	quest := &models.Quest{
		Title: "Title1",
		Description: "Description",
		Difficulty: 5,
		IsActive: true,
		XPReward: 500,
	}
	if _, err := questRepo.Create(ctx, quest); err != nil {
		log.Fatalf("create quest: %v", err)
	}
	log.Printf("quest created with ID %d", quest.ID)


	task := &models.Task{
		QuestID:     quest.ID,
		Title:       "Complete First Mission",
		Description: "Finish your first quest objective",
		IsCompleted: false,
		XPReward:    100,
	}
	if _, err := taskRepo.Create(ctx, task); err != nil {
		log.Fatalf("create task: %v", err)
	}
	log.Printf("task created with ID %d", task.ID)

}
