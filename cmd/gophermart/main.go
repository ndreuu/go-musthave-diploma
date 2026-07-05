package main

import (
	"context"
	"database/sql"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go-musthave-diploma/internal/config"
	"go-musthave-diploma/internal/handler"
	"go-musthave-diploma/internal/logger"
	"go-musthave-diploma/internal/repository"
	"go-musthave-diploma/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log, err := logger.NewLogger("info")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var storage repository.Storage

	if cfg.DatabaseURI == "" {
		storage = repository.NewMemoryRepository()
	} else {
		db, err := sql.Open("pgx", cfg.DatabaseURI)
		if err != nil {
			log.Fatal("failed to open db", zap.Error(err))
		}

		if err := db.PingContext(context.Background()); err != nil {
			log.Fatal("failed to ping db", zap.Error(err))
		}

		if err := repository.RunMigrations(db, cfg.DatabaseURI, log); err != nil {
			log.Fatal("failed to run migrations", zap.Error(err))
		}

		defer db.Close()

		storage = repository.NewPostgresRepository(db)
	}

	authService := service.NewAuthService(storage)
	orderService := service.NewOrderService(storage)
	tokenService := service.NewTokenService(cfg.JWTSecret)
	balanceService := service.NewBalanceService(storage, storage)

	if cfg.AccrualSystemAddress != "" {
		accrualClient := service.NewAccrualClient(cfg.AccrualSystemAddress)
		accrualPoller := service.NewAccrualPoller(
			storage,
			accrualClient,
			log,
			cfg.AccrualPollInterval,
			cfg.AccrualBatchSize,
		)

		go accrualPoller.Start(ctx)
	}

	router := handler.NewRouter(
		log,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		log.Info("starting server", zap.String("addr", cfg.RunAddress))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
	} else {
		log.Info("server stopped gracefully")
	}
}
