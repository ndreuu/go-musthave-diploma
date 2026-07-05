package repository

import (
	"context"
	"sync"
	"time"

	"go-musthave-diploma/internal/model"
)

// MemoryRepository is an in-memory implementation of the Storage interface.
//
// It uses maps and slices protected by a read-write mutex for thread-safe
// concurrent access. This implementation is intended for testing and
// development purposes where a full database is not required.
//
// The repository stores:
//   - users: Map keyed by login for O(1) lookup
//   - orders: Map keyed by order number for O(1) lookup
//   - withdrawals: Slice maintaining insertion order
type MemoryRepository struct {
	// mu protects all fields from concurrent access.
	mu sync.RWMutex

	// nextID is the counter for generating unique user IDs.
	nextID int64

	// users maps login strings to User pointers.
	users map[string]*model.User

	// orders maps order numbers to Order pointers.
	orders map[string]*model.Order

	// withdrawals is a slice of all withdrawal records.
	withdrawals []model.Withdrawal
}

// NewMemoryRepository creates and initializes a new in-memory repository.
//
// The repository starts empty with nextID set to 1. All internal maps
// are initialized and ready for use.
//
// Returns:
//   - *MemoryRepository: A new in-memory repository instance
//
// Example usage:
//
//	repo := repository.NewMemoryRepository()
//	user, err := repo.CreateUser(ctx, "testuser", "hash")
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		users:  make(map[string]*model.User),
		orders: make(map[string]*model.Order),
	}
}

// CreateWithdrawal adds a new withdrawal record to the repository.
//
// The withdrawal is created with the current timestamp as ProcessedAt.
// This method is thread-safe and acquires a write lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - userID: ID of the user making the withdrawal
//   - order: Order number associated with the withdrawal
//   - sum: Amount of points to withdraw
//
// Returns:
//   - error: Always nil for memory implementation
func (r *MemoryRepository) CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.withdrawals = append(r.withdrawals, model.Withdrawal{
		Order:       order,
		UserID:      userID,
		Sum:         sum,
		ProcessedAt: time.Now(),
	})

	return nil
}

// GetWithdrawalsByUserID retrieves all withdrawals for a specific user.
//
// This method is thread-safe and acquires a read lock. It returns
// withdrawals in the order they were created.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - userID: ID of the user to retrieve withdrawals for
//
// Returns:
//   - []model.Withdrawal: Slice of withdrawal records (may be empty)
//   - error: Always nil for memory implementation
func (r *MemoryRepository) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Withdrawal

	for _, withdrawal := range r.withdrawals {
		if withdrawal.UserID == userID {
			result = append(result, withdrawal)
		}
	}

	return result, nil
}

// CreateUser creates a new user with the given login and password hash.
//
// The user is assigned a unique auto-incremented ID. If a user with
// the same login already exists, returns ErrAlreadyExists.
// This method is thread-safe and acquires a write lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - login: Unique username for the new user
//   - passwordHash: Bcrypt hash of the user's password
//
// Returns:
//   - *model.User: Pointer to the created user with assigned ID
//   - error: ErrAlreadyExists if login is taken, nil otherwise
func (r *MemoryRepository) CreateUser(ctx context.Context, login string, passwordHash string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[login]; exists {
		return nil, ErrAlreadyExists
	}

	user := &model.User{
		ID:           r.nextID,
		Login:        login,
		PasswordHash: passwordHash,
	}

	r.nextID++
	r.users[login] = user

	return user, nil
}

// GetUserByLogin retrieves a user by their login username.
//
// This method is thread-safe and acquires a read lock. Returns
// ErrNotFound if no user with the given login exists.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - login: Username to look up
//
// Returns:
//   - *model.User: Pointer to the user if found
//   - error: ErrNotFound if user doesn't exist
func (r *MemoryRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[login]
	if !exists {
		return nil, ErrNotFound
	}

	return user, nil
}

// CreateOrder creates a new order for the specified user.
//
// The order is created with status NEW and the current timestamp.
// If an order with the same number already exists, it will be
// overwritten (this is intentional for idempotency).
// This method is thread-safe and acquires a write lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - userID: ID of the user uploading the order
//   - number: Unique order number (should be validated with Luhn algorithm)
//
// Returns:
//   - error: Always nil for memory implementation
func (r *MemoryRepository) CreateOrder(ctx context.Context, userID int64, number string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[number] = &model.Order{
		Number:     number,
		UserID:     userID,
		Status:     model.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	return nil
}

// GetOrderByNumber retrieves an order by its number.
//
// This method is thread-safe and acquires a read lock. Returns
// ErrNotFound if no order with the given number exists.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - number: Order number to look up
//
// Returns:
//   - *model.Order: Pointer to the order if found
//   - error: ErrNotFound if order doesn't exist
func (r *MemoryRepository) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[number]
	if !ok {
		return nil, ErrNotFound
	}

	return order, nil
}

// GetOrdersByUserID retrieves all orders for a specific user.
//
// This method is thread-safe and acquires a read lock. Returns
// orders in arbitrary order (map iteration order).
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - userID: ID of the user to retrieve orders for
//
// Returns:
//   - []model.Order: Slice of order records (may be empty)
//   - error: Always nil for memory implementation
func (r *MemoryRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Order

	for _, order := range r.orders {
		if order.UserID == userID {
			result = append(result, *order)
		}
	}

	return result, nil
}

// GetOrdersForAccrual retrieves orders that need accrual processing.
//
// Returns orders with status NEW or PROCESSING, up to the specified limit.
// This method is thread-safe and acquires a read lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - limit: Maximum number of orders to return
//
// Returns:
//   - []model.Order: Slice of orders awaiting accrual (may be empty)
//   - error: Always nil for memory implementation
func (r *MemoryRepository) GetOrdersForAccrual(ctx context.Context, limit int) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]model.Order, 0, limit)

	for _, order := range r.orders {
		if order.Status == model.OrderStatusNew ||
			order.Status == model.OrderStatusProcessing {
			orders = append(orders, *order)

			if len(orders) >= limit {
				break
			}
		}
	}

	return orders, nil
}

// UpdateOrderAccrual updates the status and accrual amount for an order.
//
// If the order doesn't exist, returns ErrNotFound. This method is
// thread-safe and acquires a write lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - number: Order number to update
//   - status: New status for the order
//   - accrual: Accrued points amount (nil if not applicable)
//
// Returns:
//   - error: ErrNotFound if order doesn't exist, nil otherwise
func (r *MemoryRepository) UpdateOrderAccrual(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[number]
	if !ok {
		return ErrNotFound
	}

	order.Status = status
	order.Accrual = accrual

	r.orders[number] = order

	return nil
}

// GetAccrualSumByUserID calculates the total accrued points for a user.
//
// Sums the accrual amounts from all orders with non-nil accrual values.
// This method is thread-safe and acquires a read lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - userID: ID of the user to calculate accrual for
//
// Returns:
//   - float64: Total accrued points (0 if no orders or no accruals)
//   - error: Always nil for memory implementation
func (r *MemoryRepository) GetAccrualSumByUserID(ctx context.Context, userID int64) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sum float64

	for _, order := range r.orders {
		if order.UserID == userID && order.Accrual != nil {
			sum += *order.Accrual
		}
	}

	return sum, nil
}

// GetWithdrawalSumByUserID calculates the total withdrawn points for a user.
//
// Sums all withdrawal amounts for the specified user. This method is
// thread-safe and acquires a read lock.
//
// Parameters:
//   - ctx: Context for cancellation (not used in memory implementation)
//   - userID: ID of the user to calculate withdrawals for
//
// Returns:
//   - float64: Total withdrawn points (0 if no withdrawals)
//   - error: Always nil for memory implementation
func (r *MemoryRepository) GetWithdrawalSumByUserID(ctx context.Context, userID int64) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sum float64

	for _, withdrawal := range r.withdrawals {
		if withdrawal.UserID == userID {
			sum += withdrawal.Sum
		}
	}

	return sum, nil
}
