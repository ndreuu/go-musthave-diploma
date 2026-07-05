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
	balance    repository.BalanceRepository
	withdrawals repository.WithdrawalRepository
}

func NewBalanceService(
	balance repository.BalanceRepository,
	withdrawals repository.WithdrawalRepository,
) *BalanceService {
	return &BalanceService{
		balance:     balance,
		withdrawals: withdrawals,
	}
}

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
