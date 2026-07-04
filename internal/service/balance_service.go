package service

import (
	"context"
	"errors"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
)

var ErrNotEnoughBalance = errors.New("not enough balance")

type Balance struct {
	Current   float64
	Withdrawn float64
}

type BalanceService struct {
	orders      repository.OrderRepository
	withdrawals repository.WithdrawalRepository
}

func NewBalanceService(
	orders repository.OrderRepository,
	withdrawals repository.WithdrawalRepository,
) *BalanceService {
	return &BalanceService{
		orders:      orders,
		withdrawals: withdrawals,
	}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int64) (*Balance, error) {
	orders, err := s.orders.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawals, err := s.withdrawals.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var accrued float64
	for _, order := range orders {
		if order.Accrual != nil {
			accrued += *order.Accrual
		}
	}

	var withdrawn float64
	for _, withdrawal := range withdrawals {
		withdrawn += withdrawal.Sum
	}

	return &Balance{
		Current:   accrued - withdrawn,
		Withdrawn: withdrawn,
	}, nil
}

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

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	return s.withdrawals.GetWithdrawalsByUserID(ctx, userID)
}
