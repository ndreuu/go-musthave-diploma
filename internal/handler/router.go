// Package handler provides HTTP request handlers for the GopherMart REST API.
//
// It implements the following endpoints:
//   - POST /api/user/register - User registration
//   - POST /api/user/login - User authentication
//   - POST /api/user/orders - Upload order number for processing
//   - GET /api/user/orders - Get user's uploaded orders
//   - GET /api/user/balance - Get user's balance information
//   - POST /api/user/balance/withdraw - Withdraw points from balance
//   - GET /api/user/withdrawals - Get withdrawal history
//
// All authenticated endpoints require a valid JWT token in the Authorization header.
package handler

import (
	"context"
	"go-musthave-diploma/internal/middleware"
	"go-musthave-diploma/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthService defines the interface for user authentication operations.
//
// Implementations handle user registration and login, returning
// user data on success or appropriate errors on failure.
type AuthService interface {
	Register(ctx context.Context, login, password string) (*model.User, error)
	Login(ctx context.Context, login, password string) (*model.User, error)
}

type TokenService interface {
	Generate(userID int64) (string, error)
	Parse(token string) (int64, error)
}

type OrderService interface {
	UploadOrder(ctx context.Context, userID int64, number string) error
	GetOrders(ctx context.Context, userID int64) ([]model.Order, error)
}

type BalanceService interface {
	GetBalance(ctx context.Context, userID int64) (*model.Balance, error)
	Withdraw(ctx context.Context, userID int64, order string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdrawal, error)
}

// Handler is the main HTTP handler that processes incoming requests
// and delegates them to appropriate service layer methods.
//
// It holds references to all required services and a logger for
// request logging and error reporting.
type Handler struct {
	// log is the structured logger for request logging.
	log *zap.Logger

	// authService handles user registration and login operations.
	authService AuthService

	// tokenService handles JWT token generation and validation.
	tokenService TokenService

	// orderService handles order upload and retrieval operations.
	orderService OrderService

	// balanceService handles balance queries and withdrawal operations.
	balanceService BalanceService
}

// NewRouter creates and configures a new Gin router with all API endpoints.
//
// The router is set up with the following middleware and routes:
//   - gin.Recovery() middleware for panic recovery
//   - Public routes: /api/user/register, /api/user/login
//   - Protected routes (require JWT auth): /api/user/orders, /api/user/balance, /api/user/withdrawals
//
// Parameters:
//   - log: Structured logger for request logging
//   - authService: Service for user authentication
//   - tokenService: Service for JWT token operations
//   - orderService: Service for order management
//   - balanceService: Service for balance and withdrawal operations
//
// Returns:
//   - Configured *gin.Engine ready to serve HTTP requests
func NewRouter(
	log *zap.Logger,
	authService AuthService,
	tokenService TokenService,
	orderService OrderService,
	balanceService BalanceService,
) *gin.Engine {
	h := &Handler{
		log:            log,
		authService:    authService,
		tokenService:   tokenService,
		orderService:   orderService,
		balanceService: balanceService,
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/api/user/register", h.Register)
	r.POST("/api/user/login", h.Login)

	auth := r.Group("/api/user")
	auth.Use(middleware.Auth(tokenService))

	auth.POST("/orders", h.UploadOrder)
	auth.GET("/orders", h.GetOrders)

	auth.GET("/balance", h.GetBalance)
	auth.POST("/balance/withdraw", h.Withdraw)

	auth.GET("/withdrawals", h.GetWithdrawals)
	return r
}
