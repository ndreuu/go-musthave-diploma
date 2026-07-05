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
	// Load configuration from flags and environment variables
	cfg := config.Load()

	// Initialize structured logger
	log, err := logger.NewLogger("info")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Set up context with graceful shutdown on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var storage repository.Storage

	// Initialize storage: in-memory for testing, PostgreSQL for production
	if cfg.DatabaseURI == "" {
		storage = repository.NewMemoryRepository()
	} else {
		// Open PostgreSQL connection
		db, err := sql.Open("pgx", cfg.DatabaseURI)
		if err != nil {
			log.Fatal("failed to open db", zap.Error(err))
		}

		// Verify database connectivity
		if err := db.PingContext(context.Background()); err != nil {
			log.Fatal("failed to ping db", zap.Error(err))
		}

		// Run database migrations
		if err := repository.RunMigrations(db, cfg.DatabaseURI, log); err != nil {
			log.Fatal("failed to run migrations", zap.Error(err))
		}

		defer db.Close()

		storage = repository.NewPostgresRepository(db)
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

		go accrualPoller.Start(ctx)
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
	go func() {
		log.Info("starting server", zap.String("addr", cfg.RunAddress))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown with 5-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
	} else {
		log.Info("server stopped gracefully")
	}
}
