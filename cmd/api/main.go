package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	rediscache "github.com/Kristex95/questhub/internal/cache"
	"github.com/Kristex95/questhub/internal/config"
	"github.com/Kristex95/questhub/internal/database"
	"github.com/Kristex95/questhub/internal/handler"
	"github.com/Kristex95/questhub/internal/logging"
	"github.com/Kristex95/questhub/internal/observability"
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

	logger, err := logging.New(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.File)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	logger.Info("starting QuestHub",
		slog.String("env", cfg.Env),
		slog.String("port", strconv.Itoa(cfg.HTTP.Port)),
	)

	shutdownTracing, err := observability.InitTracing(ctx, cfg.Observability.OTLPEndpoint, cfg.Observability.ServiceName)
	if err != nil {
		logger.Error("init tracing failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := shutdownTracing(context.Background()); err != nil {
			logger.Error("tracing shutdown error", slog.String("error", err.Error()))
		}
	}()

	pool, err := database.NewPostgresPool(ctx, cfg.Postgres.DSN(), cfg.Postgres.MaxConns)
	if err != nil {
		logger.Error("connect postgres failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("postgres connected")

	rdb, err := rediscache.NewRedisClient(ctx, rediscache.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		logger.Error("connect redis failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer rdb.Close()
	logger.Info("redis connected")

	questTTL, err := cfg.CacheQuestTTLDuration()
	if err != nil {
		logger.Error("invalid cache quest ttl", slog.String("error", err.Error()))
		os.Exit(1)
	}

	questRepo := cache.NewCachedQuestRepository(
		postgres.NewQuestRepository(pool),
		rdb,
		questTTL,
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
		logger.Error("connect notification service failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer notifier.Close()
	logger.Info("notification service connected")

	questService := service.NewQuestService(
		questRepo,
		taskRepo,
		rewardService,
		progressService,
		statsService,
		notifier,
		logger,
	)

	questHandler := handler.NewQuestHandler(questService, logger)
	leaderboardHandler := handler.NewLeaderboardHandler(statsService)
	router := handler.NewRouter(cfg.Observability.ServiceName, logger, pool, questHandler, leaderboardHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Host, strconv.Itoa(cfg.HTTP.Port)),
		Handler: router,
	}

	go func() {
		logger.Info("http server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownTimeout, err := cfg.ShutdownTimeoutDuration()
	if err != nil {
		logger.Error("invalid shutdown timeout", slog.String("error", err.Error()))
		shutdownTimeout = 15 * time.Second
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", slog.String("error", err.Error()))
	}

	logger.Info("server stopped gracefully")
}
