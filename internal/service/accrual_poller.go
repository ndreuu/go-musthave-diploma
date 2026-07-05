package service

import (
	"context"
	"errors"
	"time"

	"go-musthave-diploma/internal/repository"

	"go.uber.org/zap"
)

type AccrualPoller struct {
	repo   repository.AccrualRepository
	client *AccrualClient
	log    *zap.Logger

	interval time.Duration
	limit    int

	nextRequestAfter time.Time
}

func NewAccrualPoller(
	repo repository.AccrualRepository,
	client *AccrualClient,
	log *zap.Logger,
	interval time.Duration,
	limit int,
) *AccrualPoller {
	return &AccrualPoller{
		repo:     repo,
		client:   client,
		log:      log,
		interval: interval,
		limit:    limit,
	}
}

func (p *AccrualPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			p.log.Info("accrual poller stopped")
			return
		case <-ticker.C:
			p.runOnce(ctx)
		}
	}
}

func (p *AccrualPoller) runOnce(ctx context.Context) {
	if time.Now().Before(p.nextRequestAfter) {
		return
	}

	orders, err := p.repo.GetOrdersForAccrual(ctx, p.limit)
	if err != nil {
		p.log.Error("failed to get orders for accrual", zap.Error(err))
		return
	}

	for _, order := range orders {
		resp, err := p.client.GetOrder(ctx, order.Number)
		if err != nil {
			if errors.Is(err, ErrAccrualNoContent) {
				continue
			}

			var tooManyRequestsErr *AccrualTooManyRequestsError
			if errors.As(err, &tooManyRequestsErr) {
				p.nextRequestAfter = time.Now().Add(tooManyRequestsErr.RetryAfter)

				p.log.Warn("accrual rate limited",
					zap.Duration("retry_after", tooManyRequestsErr.RetryAfter),
					zap.Time("next_request_after", p.nextRequestAfter),
				)

				return
			}

			p.log.Error("failed to get order from accrual",
				zap.String("order", order.Number),
				zap.Error(err),
			)
			continue
		}

		if err := p.repo.UpdateOrderAccrual(ctx, order.Number, resp.Status, resp.Accrual); err != nil {
			fields := []zap.Field{
				zap.String("order", order.Number),
				zap.String("status", string(resp.Status)),
			}

			if resp.Accrual != nil {
				fields = append(fields, zap.Float64("accrual", *resp.Accrual))
			}

			p.log.Info("order updated from accrual", fields...)
			p.log.Error("failed to update order accrual",
				zap.String("order", order.Number),
				zap.Error(err),
			)
			continue
		}
	}
}
