package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"go-musthave-diploma/internal/middleware"
	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UploadOrder(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	number := strings.TrimSpace(string(body))
	if number == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.orderService.UploadOrder(c.Request.Context(), userID, number)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrderNumber):
			c.Status(http.StatusUnprocessableEntity)
		case errors.Is(err, service.ErrOrderAlreadyUploaded):
			c.Status(http.StatusOK)
		case errors.Is(err, service.ErrOrderUploadedByAnother):
			c.Status(http.StatusConflict)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusAccepted)
}

type orderResponse struct {
	Number     string            `json:"number"`
	Status     model.OrderStatus `json:"status"`
	Accrual    *float64          `json:"accrual,omitempty"`
	UploadedAt time.Time         `json:"uploaded_at"`
}

func (h *Handler) GetOrders(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetOrders(c.Request.Context(), userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resp := make([]orderResponse, 0, len(orders))
	for _, order := range orders {
		resp = append(resp, orderResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}
