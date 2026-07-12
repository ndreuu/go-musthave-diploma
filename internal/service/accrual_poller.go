package service

import (
	"context"
	"errors"
	"time"

	"go-musthave-diploma/internal/repository"

	"go.uber.org/zap"
)

// AccrualPoller is a background worker that periodically polls the accrual system
// for order status updates.
//
// It runs on a configurable interval, fetching orders in NEW or PROCESSING status
// and querying the accrual system for their current status. The poller handles
// rate limiting by respecting the Retry-After header and scheduling the next
// request accordingly.
type AccrualPoller struct {
	// repo is the repository for fetching and updating orders.
	repo repository.AccrualRepository

	// client is the HTTP client for communicating with the accrual system.
	client *AccrualClient

	// log is the logger for poller activity and errors.
	log *zap.Logger

	// interval is the polling interval between run cycles.
	interval time.Duration

	// limit is the maximum number of orders to fetch per polling cycle.
	limit int

	// nextRequestAfter is the earliest time when the next accrual request can be made.
	// Used for rate limiting.
	nextRequestAfter time.Time
}

// NewAccrualPoller creates a new AccrualPoller instance.
//
// Parameters:
//   - repo: AccrualRepository implementation for order data access
//   - client: AccrualClient for communicating with the accrual system
//   - log: Structured logger for poller activity
//   - interval: Time between polling cycles
//   - limit: Maximum number of orders to process per cycle
//
// Returns:
//   - *AccrualPoller: New accrual poller instance
//
// Example usage:
//
//	poller := service.NewAccrualPoller(repo, client, logger, 10*time.Second, 100)
//	go poller.Start(ctx)
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

// Start begins the accrual polling loop.
//
// The method runs an initial poll immediately, then continues polling at the
// configured interval until the context is cancelled. Each polling cycle:
//  1. Checks if rate limiting is in effect (skips if before nextRequestAfter)
//  2. Fetches orders in NEW or PROCESSING status
//  3. Queries the accrual system for each order's status
//  4. Updates order status and accrual amounts in the repository
//
// The method blocks until the context is cancelled. It should be called as a
// goroutine in production code.
//
// Parameters:
//   - ctx: Context for cancellation (stops the poller when cancelled)
//
// Example usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	go poller.Start(ctx)
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

// runOnce executes a single polling cycle.
//
// It fetches orders needing accrual processing and queries the accrual system
// for each one. Rate limiting is handled by setting nextRequestAfter when a
// 429 response is received.
//
// Errors during individual order processing are logged but don't stop the
// entire cycle.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
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
				p.nextRequestAfter = time.Now().Add(
					tooManyRequestsErr.RetryAfter,
				)

				p.log.Warn(
					"accrual rate limited",
					zap.Duration(
						"retry_after",
						tooManyRequestsErr.RetryAfter,
					),
					zap.Time(
						"next_request_after",
						p.nextRequestAfter,
					),
				)

				return
			}

			p.log.Error(
				"failed to get order from accrual",
				zap.String("order", order.Number),
				zap.Error(err),
			)
			continue
		}

		if err := p.repo.UpdateOrderAccrual(
			ctx,
			order.Number,
			resp.Status,
			resp.Accrual,
		); err != nil {
			p.log.Error(
				"failed to update order accrual",
				zap.String("order", order.Number),
				zap.String("status", string(resp.Status)),
				zap.Error(err),
			)
			continue
		}

		fields := []zap.Field{
			zap.String("order", order.Number),
			zap.String("status", string(resp.Status)),
		}

		if resp.Accrual != nil {
			fields = append(
				fields,
				zap.Float64("accrual", *resp.Accrual),
			)
		}

		p.log.Info("order updated from accrual", fields...)
	}
	
}
