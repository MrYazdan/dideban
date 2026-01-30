// Package checks defines API request/response types for check management endpoints.
//
// This package provides strongly-typed DTOs (Data Transfer Objects) for all
// check-related API operations, ensuring type safety, validation, and clear
// contract between frontend and backend.
package checks

import (
	"encoding/json"
	"time"
)

// CheckRequest represents the request payload for creating or updating a check.
//
// Required fields:
//   - name: Unique identifier for the check (max 100 chars)
//   - type: Monitoring type (e.g., "http", "ping")
//   - target: URL or host to monitor
//   - interval_seconds: Execution frequency in seconds
//
// Optional fields:
//   - enabled: Whether the check is active (default: true)
//   - timeout_seconds: Maximum execution time (default: 30)
//   - config: JSON-encoded type-specific configuration
type CheckRequest struct {
	Name            *string          `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Type            *string          `json:"type,omitempty" binding:"omitempty,oneof=http ping"`
	Target          *string          `json:"target,omitempty" binding:"omitempty"`
	Config          *json.RawMessage `json:"config,omitempty"`
	IntervalSeconds *int             `json:"interval_seconds,omitempty" binding:"omitempty,min=10,max=86400"`
	TimeoutSeconds  *int             `json:"timeout_seconds,omitempty" binding:"omitempty,min=5,max=300"`
	Enabled         *bool            `json:"enabled,omitempty"`
}

// CheckResponse represents a check entity in API responses.
//
// This struct is used for all read operations (GET /list, GET /:id) and includes
// both configuration metadata and runtime state (status, last checked timestamp).
type CheckResponse struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	Target          string          `json:"target"`
	Config          json.RawMessage `json:"config"`
	Enabled         bool            `json:"enabled"`
	IntervalSeconds int             `json:"interval_seconds"`
	TimeoutSeconds  int             `json:"timeout_seconds"`
	Status          string          `json:"status"`       // Computed: "up", "down", "error", "timeout"
	LastCheckedAt   *time.Time      `json:"last_checked"` // Latest check_history.checked_at
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// CheckHistoryRequest represents the request payload for submitting a check result.
//
// This payload is typically sent by the monitoring engine after executing a check.
// It captures the outcome, performance metrics, and any errors encountered.
type CheckHistoryRequest struct {
	Status         string  `json:"status" binding:"required,oneof=up down error timeout"`
	ResponseTimeMs *int    `json:"response_time_ms,omitempty" binding:"omitempty,min=0"`
	StatusCode     *int    `json:"status_code,omitempty" binding:"omitempty,min=100,max=599"`
	ErrorMessage   *string `json:"error_message,omitempty" binding:"omitempty,max=1000"`
}

// CheckHistoryResponse represents a historical check execution record.
//
// This struct is used for retrieving past check results, enabling trend analysis,
// performance monitoring, and incident investigation.
type CheckHistoryResponse struct {
	ID             int64     `json:"id"`
	CheckID        int64     `json:"check_id"`
	Status         string    `json:"status"`
	ResponseTimeMs *int      `json:"response_time_ms,omitempty"`
	StatusCode     *int      `json:"status_code,omitempty"`
	ErrorMessage   *string   `json:"error_message,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

// CheckStatusSummary provides aggregated statistics for a check's performance.
//
// This struct is used in dashboard views and analytics endpoints to show
// high-level metrics without fetching full history records.
type CheckStatusSummary struct {
	TotalChecks   int64 `json:"total_checks"`
	SuccessCount  int64 `json:"success_count"`   // Status = "up"
	FailureCount  int64 `json:"failure_count"`   // Status = "down"
	ErrorCount    int64 `json:"error_count"`     // Status = "error" or "timeout"
	AvgResponseMs int64 `json:"avg_response_ms"` // Average of non-null response times
}
