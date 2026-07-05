package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-musthave-diploma/internal/handler"
	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
	"go-musthave-diploma/internal/service"
	"go.uber.org/zap"
)

func TestNewRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	assert.NotNil(t, router)

	routes := router.Routes()
	assert.GreaterOrEqual(t, len(routes), 7)
}

func TestRegister_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	body := bytes.NewBufferString(`{"login":"","password":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBalance_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetOrders_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUploadOrder_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("12345678901"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithdraw_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	body := bytes.NewBufferString(`{"order":"12345678901","sum":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetWithdrawals_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	body := bytes.NewBufferString(`{"login":"testuser","password":"testpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	body := bytes.NewBufferString(`{"login":"testuser","password":"testpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(`{"login":"testuser","password":"testpass"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	body := bytes.NewBufferString(`{"login":"nonexistent","password":"wrongpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUploadOrder_InvalidLuhn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user2","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("12345"))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUploadOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user3","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("2377225624"))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestUploadOrder_AlreadyUploaded(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user4","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("2377225624"))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("2377225624"))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrders_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user5","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("2377225624"))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrders_NoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user6","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetBalance_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user7","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWithdraw_InvalidOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user8","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	body := bytes.NewBufferString(`{"order":"12345","sum":10}`)
	req = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestWithdraw_NotEnoughBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user9","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	body := bytes.NewBufferString(`{"order":"2377225624","sum":1000}`)
	req = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPaymentRequired, w.Code)
}

func TestWithdraw_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user10","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	body := bytes.NewBufferString(`{"order":"","sum":0}`)
	req = httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetWithdrawals_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger, _ := zap.NewDevelopment()
	repo := repository.NewMemoryRepository()

	authService := service.NewAuthService(repo)
	orderService := service.NewOrderService(repo)
	tokenService := service.NewTokenService("test-secret")
	balanceService := service.NewBalanceService(repo, repo)

	router := handler.NewRouter(
		logger,
		authService,
		tokenService,
		orderService,
		balanceService,
	)

	registerBody := bytes.NewBufferString(`{"login":"user11","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", registerBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := w.Header().Get("Authorization")[7:]

	req = httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}


func TestBalanceResponse(t *testing.T) {
	resp := balanceResponse{
		Current:   100.5,
		Withdrawn: 50.0,
	}
	assert.Equal(t, 100.5, resp.Current)
	assert.Equal(t, 50.0, resp.Withdrawn)
}

func TestWithdrawRequest(t *testing.T) {
	req := withdrawRequest{
		Order: "12345678901",
		Sum:   100.5,
	}
	assert.Equal(t, "12345678901", req.Order)
	assert.Equal(t, 100.5, req.Sum)
}

func TestOrderResponse(t *testing.T) {
	status := model.OrderStatusNew
	resp := orderResponse{
		Number: "12345678901",
		Status: status,
	}
	assert.Equal(t, "12345678901", resp.Number)
	assert.Equal(t, model.OrderStatusNew, resp.Status)
}

func TestWithdrawalResponse(t *testing.T) {
	resp := withdrawalResponse{
		Order: "12345678901",
		Sum:   100.5,
	}
	assert.Equal(t, "12345678901", resp.Order)
	assert.Equal(t, 100.5, resp.Sum)
}

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type orderResponse struct {
	Number     string            `json:"number"`
	Status     model.OrderStatus `json:"status"`
	Accrual    *float64          `json:"accrual,omitempty"`
	UploadedAt interface{}       `json:"uploaded_at"`
}

type withdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
