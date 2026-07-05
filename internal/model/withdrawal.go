package model

import "time"

// Withdrawal represents a withdrawal of loyalty points from a user's balance.
//
// Withdrawals reduce the user's current balance and are recorded for auditing
// purposes. Each withdrawal is associated with an order number and has a
// specific amount being withdrawn.
type Withdrawal struct {
	// Order is the order number associated with this withdrawal.
	Order string

	// UserID is the ID of the user who made this withdrawal.
	UserID int64

	// Sum is the amount of points withdrawn.
	Sum float64

	// ProcessedAt is the timestamp when the withdrawal was processed.
	ProcessedAt time.Time
}
