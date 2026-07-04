package handler

import (
	"github.com/gin-gonic/gin"
	"go-musthave-diploma/internal/middleware"
	"go-musthave-diploma/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	log            *zap.Logger
	authService    *service.AuthService
	tokenService   *service.TokenService
	orderService   *service.OrderService
	balanceService *service.BalanceService
}

func NewRouter(
	log *zap.Logger,
	authService *service.AuthService,
	tokenService *service.TokenService,
	orderService *service.OrderService,
	balanceService *service.BalanceService,
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
