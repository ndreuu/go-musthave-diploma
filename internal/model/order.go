package model

import "time"

// OrderStatus represents the current processing state of an order.
//
// Orders transition through the following states:
//   - NEW: Order has been uploaded but not yet sent to accrual system
//   - PROCESSING: Order has been sent to accrual system, awaiting response
//   - PROCESSED: Order has been processed, accrual amount is determined
//   - INVALID: Order number failed validation or was rejected
type OrderStatus string

const (
	// OrderStatusNew indicates a newly uploaded order awaiting processing.
	OrderStatusNew OrderStatus = "NEW"

	// OrderStatusProcessing indicates the order is being processed by the accrual system.
	OrderStatusProcessing OrderStatus = "PROCESSING"

	// OrderStatusInvalid indicates the order was rejected or has an invalid number.
	OrderStatusInvalid OrderStatus = "INVALID"

	// OrderStatusProcessed indicates the order was successfully processed with accrual determined.
	OrderStatusProcessed OrderStatus = "PROCESSED"
)

// Order represents a customer order uploaded for loyalty point accrual.
//
// Each order belongs to a specific user and goes through a lifecycle
// from NEW to either PROCESSED (with points accrued) or INVALID.
type Order struct {
	// Number is the unique order identifier (validated with Luhn algorithm).
	Number string

	// UserID is the ID of the user who uploaded this order.
	UserID int64

	// Status is the current processing status of the order.
	Status OrderStatus

	// Accrual is the amount of points accrued for this order.
	// Nil if the order hasn't been processed yet or was marked invalid.
	Accrual *float64

	// UploadedAt is the timestamp when the order was uploaded by the user.
	UploadedAt time.Time
}
