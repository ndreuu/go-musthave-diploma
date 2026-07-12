package model

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
