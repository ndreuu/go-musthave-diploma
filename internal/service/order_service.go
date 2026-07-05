package service

import (
	"context"
	"errors"
	"sort"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
)

var (
	// ErrInvalidOrderNumber is returned when an order number fails Luhn validation.
	ErrInvalidOrderNumber = errors.New("invalid order number")

	// ErrOrderUploadedByAnother is returned when attempting to upload an order
	// that was already uploaded by a different user.
	ErrOrderUploadedByAnother = errors.New("order uploaded by another user")

	// ErrOrderAlreadyUploaded is returned when a user attempts to upload an order
	// they have already uploaded previously.
	ErrOrderAlreadyUploaded = errors.New("order already uploaded")
)

// OrderService handles order management operations.
//
// It provides functionality for uploading order numbers and retrieving
// user orders. All order numbers are validated using the Luhn algorithm
// before being accepted.
type OrderService struct {
	// orders is the repository for order data access.
	orders repository.OrderRepository
}

// NewOrderService creates a new OrderService instance.
//
// Parameters:
//   - orders: OrderRepository implementation for data access
//
// Returns:
//   - *OrderService: New order service instance
//
// Example usage:
//
//	orderService := service.NewOrderService(orderRepo)
//	err := orderService.UploadOrder(ctx, userID, "12345678901")
func NewOrderService(orders repository.OrderRepository) *OrderService {
	return &OrderService{orders: orders}
}

// UploadOrder uploads a new order number for the specified user.
//
// The function performs the following steps:
//  1. Validates the order number using the Luhn algorithm
//  2. Checks if the order already exists in the system
//  3. If exists and belongs to the same user, returns ErrOrderAlreadyUploaded
//  4. If exists and belongs to another user, returns ErrOrderUploadedByAnother
//  5. If not found, creates a new order with status NEW
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - userID: ID of the user uploading the order
//   - number: Order number to upload (must pass Luhn validation)
//
// Returns:
//   - error: ErrInvalidOrderNumber if Luhn check fails,
//     ErrOrderAlreadyUploaded if user already uploaded this order,
//     ErrOrderUploadedByAnother if another user uploaded this order,
//     or other repository error
//
// Example usage:
//
//	err := orderService.UploadOrder(ctx, userID, "12345678901")
//	if err != nil {
//	    switch {
//	    case errors.Is(err, service.ErrInvalidOrderNumber):
//	        // Handle invalid order number
//	    case errors.Is(err, service.ErrOrderUploadedByAnother):
//	        // Handle conflict
//	    }
//	}
func (s *OrderService) UploadOrder(ctx context.Context, userID int64, number string) error {
	if !IsValidLuhn(number) {
		return ErrInvalidOrderNumber
	}

	order, err := s.orders.GetOrderByNumber(ctx, number)
	if err == nil {
		if order.UserID == userID {
			return ErrOrderAlreadyUploaded
		}
		return ErrOrderUploadedByAnother
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	return s.orders.CreateOrder(ctx, userID, number)
}

// GetOrders retrieves all orders for the specified user.
//
// Orders are sorted by upload timestamp in descending order (newest first).
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - userID: ID of the user to retrieve orders for
//
// Returns:
//   - []model.Order: Slice of orders (may be empty if user has no orders)
//   - error: Repository error
//
// Example usage:
//
//	orders, err := orderService.GetOrders(ctx, userID)
//	if err != nil {
//	    return err
//	}
//	for _, order := range orders {
//	    fmt.Printf("Order %s: %s\n", order.Number, order.Status)
//	}
func (s *OrderService) GetOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	orders, err := s.orders.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.After(orders[j].UploadedAt)
	})

	return orders, nil
}
