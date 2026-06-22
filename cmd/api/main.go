package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	rediscache "github.com/Kristex95/questhub/internal/cache"
	"github.com/Kristex95/questhub/internal/config"
	"github.com/Kristex95/questhub/internal/database"
	"github.com/Kristex95/questhub/internal/handler"
	"github.com/Kristex95/questhub/internal/repository/cache"
	"github.com/Kristex95/questhub/internal/repository/postgres"
	"github.com/Kristex95/questhub/internal/repository/stats"
	"github.com/Kristex95/questhub/internal/service"
	"github.com/Kristex95/questhub/internal/transport/grpc"
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

	rdb, err := rediscache.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("connect redis: %v", err)
	}
	defer rdb.Close()
	log.Println("redis connected")

	questRepo := cache.NewCachedQuestRepository(
		postgres.NewQuestRepository(pool),
		rdb,
		cfg.Cache.QuestTTL,
	)
	taskRepo := postgres.NewTaskRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	rewardRepo := postgres.NewRewardRepository(pool)
	progressRepo := postgres.NewProgressRepository(pool)

	statsService := stats.NewStatsService(rdb)
	rewardService := service.NewRewardService(rewardRepo, userRepo, statsService)
	progressService := service.NewProgressService(progressRepo)
	notifier, err := grpc.NewNotificationClient(cfg.GRPC.NotificationAddr)
	if err != nil {
		log.Fatalf("connect notification service: %v", err)
	}
	defer notifier.Close()
	log.Println("notification service connected")

	questService := service.NewQuestService(
		questRepo,
		taskRepo,
		rewardService,
		progressService,
		statsService,
		notifier,
	)

	questHandler := handler.NewQuestHandler(questService)
	leaderboardHandler := handler.NewLeaderboardHandler(statsService)
	router := handler.NewRouter(questHandler, leaderboardHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}

	go func() {
		log.Printf("http server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped gracefully")
}
