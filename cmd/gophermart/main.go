// GopherMart is a loyalty rewards system for processing customer orders and managing points.
//
// This application provides a REST API for:
//   - User registration and authentication
//   - Order number upload for loyalty point accrual
//   - Balance checking and point withdrawal
//   - Order status tracking
//
// The system integrates with an external accrual system to determine points for orders.
// A background poller periodically checks the accrual system for order status updates.
//
// Configuration:
//   - Command-line flags: -a (address), -d (database URI), -r (accrual URL), -s (JWT secret)
//   - Environment variables: RUN_ADDRESS, DATABASE_URI, ACCRUAL_SYSTEM_ADDRESS, JWT_SECRET
//
// If DATABASE_URI is empty, the application runs with in-memory storage (for testing).
// Otherwise, it uses PostgreSQL with automatic migration application.
//
// Example usage:
//
//	./gophermart -a :8080 -d "postgres://user:pass@localhost/db" -r "http://accrual:8080"
package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go-musthave-diploma/internal/config"
	"go-musthave-diploma/internal/handler"
	"go-musthave-diploma/internal/logger"
	"go-musthave-diploma/internal/repository"
	"go-musthave-diploma/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Initialize structured logger
	log, err := logger.NewLogger("info")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Load configuration from flags and environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	// Set up context with graceful shutdown on SIGINT/SIGTERM
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancel(signalCtx)
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	var storage repository.Storage
	
	// Initialize storage: in-memory for testing, PostgreSQL for production
	if cfg.DatabaseURI == "" {
		storage = repository.NewMemoryRepository()
	} else {
		pool, err := pgxpool.New(ctx, cfg.DatabaseURI)
		if err != nil {
			log.Fatal("failed to create database pool", zap.Error(err))
		}

		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			log.Fatal("failed to ping database", zap.Error(err))
		}

		migrationDB, err := sql.Open("pgx", cfg.DatabaseURI)
		if err != nil {
			pool.Close()
			log.Fatal("failed to open migration database", zap.Error(err))
		}

		if err := repository.RunMigrations(migrationDB, log); err != nil {
			log.Fatal("failed to run migrations", zap.Error(err))
		}

		if err := migrationDB.Close(); err != nil {
			log.Error("failed to close migration database", zap.Error(err))
		}

		defer pool.Close()

		storage = repository.NewPostgresRepository(pool)
	}

	// Initialize services
	authService := service.NewAuthService(storage)
	orderService := service.NewOrderService(storage)
	tokenService := service.NewTokenService(cfg.JWTSecret)
	balanceService := service.NewBalanceService(storage, storage)

	// Start accrual poller if accrual system is configured
	if cfg.AccrualSystemAddress != "" {
		accrualClient := service.NewAccrualClient(cfg.AccrualSystemAddress)
		accrualPoller := service.NewAccrualPoller(
			storage,
			accrualClient,
			log,
			cfg.AccrualPollInterval,
			cfg.AccrualBatchSize,
		)

		group.Go(func() error {
			accrualPoller.Start(groupCtx)
			return nil
		})

	}

	// Set up HTTP router with all handlers
	router := handler.NewRouter(
		log,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	// Configure HTTP server
	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	// Start server in a goroutine
	group.Go(func() error {
		defer cancel()

		log.Info("starting server", zap.String("addr", cfg.RunAddress))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	// Wait for shutdown signal
	group.Go(func() error {
		<-groupCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		log.Error("application stopped with error", zap.Error(err))
	} else {
		log.Info("application stopped gracefully")
	}
}
