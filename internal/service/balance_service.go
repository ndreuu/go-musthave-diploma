package service

import (
	"context"
	"errors"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
)

// ErrNotEnoughBalance is returned when attempting to withdraw more points
// than are available in the user's current balance.
var ErrNotEnoughBalance = errors.New("not enough balance")

// Balance represents a user's loyalty points balance.
//
// It includes both the current available balance and the total amount
// withdrawn by the user.
type Balance struct {
	// Current is the available balance (accrued - withdrawn).
	Current float64

	// Withdrawn is the total amount of points withdrawn by the user.
	Withdrawn float64
}

// BalanceService handles balance queries and withdrawal operations.
//
// It calculates user balances based on accrued points from orders
// and tracks withdrawals. All withdrawals are validated using the
// Luhn algorithm for order number verification.
type BalanceService struct {
	// balance is the repository for balance-related queries.
	balance repository.BalanceRepository

	// withdrawals is the repository for withdrawal operations.
	withdrawals repository.WithdrawalRepository
}

// NewBalanceService creates a new BalanceService instance.
//
// Parameters:
//   - balance: BalanceRepository implementation for accrual/withdrawal sum queries
//   - withdrawals: WithdrawalRepository implementation for withdrawal operations
//
// Returns:
//   - *BalanceService: New balance service instance
//
// Example usage:
//
//	balanceService := service.NewBalanceService(balanceRepo, withdrawalRepo)
//	balance, err := balanceService.GetBalance(ctx, userID)
func NewBalanceService(
	balance repository.BalanceRepository,
	withdrawals repository.WithdrawalRepository,
) *BalanceService {
	return &BalanceService{
		balance:     balance,
		withdrawals: withdrawals,
	}
}

// GetBalance retrieves the current balance for the specified user.
//
// The balance is calculated as:
//   - Current = Total Accrued - Total Withdrawn
//   - Withdrawn = Sum of all withdrawals
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - userID: ID of the user to retrieve balance for
//
// Returns:
//   - *Balance: Pointer to the user's balance information
//   - error: Repository error
//
// Example usage:
//
//	balance, err := balanceService.GetBalance(ctx, userID)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Current: %.2f, Withdrawn: %.2f\n", balance.Current, balance.Withdrawn)
func (s *BalanceService) GetBalance(ctx context.Context, userID int64) (*Balance, error) {
	accrued, err := s.balance.GetAccrualSumByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawn, err := s.balance.GetWithdrawalSumByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Current:   accrued - withdrawn,
		Withdrawn: withdrawn,
	}, nil
}

// Withdraw processes a withdrawal request for the specified user.
//
// The function performs the following steps:
//  1. Validates the order number using the Luhn algorithm
//  2. Retrieves the user's current balance
//  3. Checks if the user has sufficient balance
//  4. Creates a withdrawal record if all checks pass
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - userID: ID of the user making the withdrawal
//   - order: Order number associated with the withdrawal (must pass Luhn check)
//   - sum: Amount of points to withdraw (must be positive)
//
// Returns:
//   - error: ErrInvalidOrderNumber if Luhn check fails,
//     ErrNotEnoughBalance if insufficient funds, or other error
//
// Example usage:
//
//	err := balanceService.Withdraw(ctx, userID, "12345678901", 100.5)
//	if err != nil {
//	    if errors.Is(err, service.ErrNotEnoughBalance) {
//	        // Handle insufficient balance
//	    }
//	}
func (s *BalanceService) Withdraw(ctx context.Context, userID int64, order string, sum float64) error {
	if !IsValidLuhn(order) {
		return ErrInvalidOrderNumber
	}

	balance, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current < sum {
		return ErrNotEnoughBalance
	}

	return s.withdrawals.CreateWithdrawal(ctx, userID, order, sum)
}

// GetWithdrawals retrieves all withdrawals for the specified user.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - userID: ID of the user to retrieve withdrawals for
//
// Returns:
//   - []model.Withdrawal: Slice of withdrawal records (may be empty)
//   - error: Repository error
//
// Example usage:
//
//	withdrawals, err := balanceService.GetWithdrawals(ctx, userID)
//	if err != nil {
//	    return err
//	}
//	for _, w := range withdrawals {
//	    fmt.Printf("Order: %s, Sum: %.2f, Date: %s\n", w.Order, w.Sum, w.ProcessedAt)
//	}
func (s *BalanceService) GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	return s.withdrawals.GetWithdrawalsByUserID(ctx, userID)
}
