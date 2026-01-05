// Package checks provides base functionality for all checker implementations.
package checks

import (
	"time"

	"dideban/internal/storage"
)

// BaseChecker provides common functionality for all checker implementations.
// It contains shared methods for configuration parsing, error handling, and result creation.
type BaseChecker struct{}

// NewBaseChecker creates a new base checker instance.
func NewBaseChecker() *BaseChecker {
	return &BaseChecker{}
}

// CreateErrorResult creates a standardized error result for failed checks.
// This ensures consistent error reporting across all checker types.
//
// Parameters:
//   - checkID: Check ID
//   - status: Error status (error, timeout, down)
//   - err: Error that occurred
//   - responseTime: Response time in milliseconds
//   - statusCode: Optional HTTP status code
//
// Returns:
//   - *storage.CheckHistory: Standardized error result
func (b *BaseChecker) CreateErrorResult(checkID int64, status string, err error, responseTime int64, statusCode *int) *storage.CheckHistory {
	errorMsg := err.Error()
	responseTimeMs := int(responseTime)

	return &storage.CheckHistory{
		CheckID:        checkID,
		Status:         status,
		ResponseTimeMs: &responseTimeMs,
		StatusCode:     statusCode,
		ErrorMessage:   &errorMsg,
		CheckedAt:      time.Now(),
	}
}

// CreateSuccessResult creates a standardized success result for passed checks.
// This ensures consistent success reporting across all checker types.
//
// Parameters:
//   - checkID: Check ID
//   - responseTime: Response time in milliseconds
//   - statusCode: Optional HTTP status code
//   - message: Optional success message
//
// Returns:
//   - *storage.CheckHistory: Standardized success result
func (b *BaseChecker) CreateSuccessResult(checkID int64, responseTime int64, statusCode *int, message string) *storage.CheckHistory {
	responseTimeMs := int(responseTime)
	errorMsg := message

	return &storage.CheckHistory{
		CheckID:        checkID,
		Status:         storage.CheckStatusUp,
		ResponseTimeMs: &responseTimeMs,
		StatusCode:     statusCode,
		ErrorMessage:   &errorMsg,
		CheckedAt:      time.Now(),
	}
}

// DetermineErrorStatus determines the appropriate error status based on the error type.
// This provides consistent error categorization across all checker types.
//
// Parameters:
//   - err: Error that occurred
//
// Returns:
//   - string: Appropriate error status (timeout, down, error)
func (b *BaseChecker) DetermineErrorStatus(err error) string {
	errStr := err.Error()

	// Check for timeout errors
	if contains(errStr, "timeout") || contains(errStr, "deadline exceeded") {
		return storage.CheckStatusTimeout
	}

	// Check for connection errors (host down)
	if contains(errStr, "connection refused") ||
		contains(errStr, "no route to host") ||
		contains(errStr, "host unreachable") ||
		contains(errStr, "100% packet loss") {
		return storage.CheckStatusDown
	}

	// Default to generic error
	return storage.CheckStatusError
}

// contains checks if a string contains a substring (case-insensitive helper).
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOf(s, substr) >= 0)))
}

// indexOf finds the index of substring in string.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
