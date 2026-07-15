// Package service provides business logic for the GopherMart loyalty system.
//
// It includes services for:
//   - User authentication (AuthService)
//   - JWT token generation and validation (TokenService)
//   - Order management (OrderService)
//   - Balance and withdrawal operations (BalanceService)
//   - Accrual system integration (AccrualClient, AccrualPoller)
//
// The package also provides utility functions like Luhn algorithm validation
// for order number verification.
package service

// IsValidLuhn validates a number string using the Luhn algorithm.
//
// The Luhn algorithm (also known as the "modulus 10" or "mod 10" algorithm)
// is a simple checksum formula used to validate a variety of identification
// numbers, such as credit card numbers, IMEI numbers, and order numbers.
//
// The algorithm works as follows:
//  1. Starting from the rightmost digit, double every second digit
//  2. If doubling results in a number > 9, subtract 9 from it
//  3. Sum all the digits
//  4. If the total modulo 10 equals 0, the number is valid
//
// This function also validates that all characters are digits (0-9).
// Empty strings and strings containing non-digit characters return false.
//
// Parameters:
//   - number: The string to validate (should contain only digits)
//
// Returns:
//   - bool: True if the number passes the Luhn check, false otherwise
//
// Example usage:
//
//	service.IsValidLuhn("12345678901")  // May return true or false depending on checksum
//	service.IsValidLuhn("")             // Returns false (empty string)
//	service.IsValidLuhn("12345abc")     // Returns false (contains non-digits)
func IsValidLuhn(number string) bool {
	if number == "" {
		return false
	}

	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		if number[i] < '0' || number[i] > '9' {
			return false
		}

		n := int(number[i] - '0')

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}
