// Package checks implements HTTP handlers for check management operations.
//
// This package provides a complete RESTful API for CRUD operations on monitoring checks,
// including history tracking, status aggregation, and runtime metrics. All handlers
// follow production-grade patterns:
//   - Comprehensive input validation
//   - Proper error handling with contextual messages
//   - Pagination for list endpoints
//   - Security-conscious responses (no sensitive data leakage)
//   - Idempotent operations where applicable
package checks

import (
	"dideban/internal/core"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"dideban/internal/api/types"
	"dideban/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Handler manages all check-related HTTP endpoints.
//
// The handler is stateless and delegates all business logic to the storage layer.
// It focuses on:
//   - Request parsing and validation
//   - Error translation (database errors → HTTP status codes)
//   - Response formatting (DTO mapping, pagination)
//   - Security (input sanitization, output masking)
type Handler struct {
	storage *storage.Storage
	engine  *core.Engine
}

// NewHandler creates a new check handler instance.
//
// Parameters:
//   - storage: Database abstraction layer for CRUD operations
//
// Returns:
//   - Pointer to initialized Handler
func NewHandler(storage *storage.Storage, engine *core.Engine) *Handler {
	return &Handler{
		storage: storage,
		engine:  engine,
	}
}

// List handles GET /api/v1/checks
//
// Returns a paginated list of checks ordered by ID in descending order.
// Each check includes its current runtime status and last checked timestamp,
// computed from the most recent check_history record.
//
// Query parameters:
//   - page (default: 1, min: 1)
//   - page_size (default: 50, max: 500)
//   - enabled (optional boolean filter)
//
// Security considerations:
//   - No sensitive data exposed (config is returned as-is, assumed safe)
//   - Pagination prevents large result sets
//
// Returns:
//   - 200 OK with paginated check list and summary statistics
//   - 400 Bad Request for invalid pagination parameters
//   - 500 Internal Server Error on storage failure
func (h *Handler) List(c *gin.Context) {
	var pagination types.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		types.AbortWithError(c, types.ValidationError(err.Error()))
		return
	}

	// Apply optional enabled filter
	enabledFilter := c.Query("enabled")
	query := h.storage.DB().Model(&storage.Check{})

	if enabledFilter != "" {
		enabled, err := strconv.ParseBool(enabledFilter)
		if err == nil {
			query = query.Where("enabled = ?", enabled)
		}
	}

	// Count total checks matching filters
	var total int64
	if err := query.Count(&total).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to count checks", err))
		return
	}

	// Calculate offset and retrieve paginated checks
	offset := (pagination.Page - 1) * pagination.PageSize
	var checks []storage.Check
	if err := query.
		Order("id DESC").
		Limit(pagination.PageSize).
		Offset(offset).
		Find(&checks).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to retrieve checks", err))
		return
	}

	// Map to response DTOs and enrich with runtime status
	responses := make([]CheckResponse, 0, len(checks))
	for _, check := range checks {
		status, lastChecked := h.resolveCheckStatus(check.ID)

		responses = append(responses, CheckResponse{
			ID:              check.ID,
			Name:            check.Name,
			Type:            check.Type,
			Target:          check.Target,
			Config:          parseConfigToRawMessage(check.Config),
			Enabled:         check.Enabled,
			IntervalSeconds: check.IntervalSeconds,
			TimeoutSeconds:  check.TimeoutSeconds,
			Status:          status,
			LastCheckedAt:   lastChecked,
			CreatedAt:       check.CreatedAt,
			UpdatedAt:       check.UpdatedAt,
		})
	}

	// Calculate total pages
	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))
	paginationResponse := &types.PaginationResponse{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, types.SuccessResponseWithPagination(responses, paginationResponse))
}

// Stats handles GET /api/v1/checks/stats
//
// Returns a summary of all checks, including total count, enabled/disabled breakdown,
// and status distribution (only for enabled checks).
//
// The status counts (up, down, error, timeout) are computed only for **enabled** checks,
// as disabled checks are not scheduled and thus have no meaningful runtime status.
//
// Returns:
//   - 200 OK with check statistics
//   - 500 Internal Server Error on storage failure
func (h *Handler) Stats(c *gin.Context) {
	db := h.storage.DB()

	// 1. Total count
	var total int64
	if err := db.Model(&storage.Check{}).Count(&total).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to count total checks", err))
		return
	}

	// 2. Enabled vs Disabled
	var enabledCount, disabledCount int64
	if err := db.Model(&storage.Check{}).Where("enabled = ?", true).Count(&enabledCount).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to count enabled checks", err))
		return
	}
	disabledCount = total - enabledCount

	// 3. Status distribution for enabled checks using latest history
	// We use a CTE or subquery to get the latest history per check
	// This version works reliably on both SQLite and PostgreSQL

	query := `
WITH latest_history AS (
    SELECT DISTINCT ON (check_id) 
        id, check_id, status, checked_at
    FROM check_history
    ORDER BY check_id, checked_at DESC
)
SELECT
    COALESCE(SUM(CASE WHEN lh.status = 'up' THEN 1 ELSE 0 END), 0) AS up,
    COALESCE(SUM(CASE WHEN lh.status = 'down' THEN 1 ELSE 0 END), 0) AS down,
    COALESCE(SUM(CASE WHEN lh.status = 'error' THEN 1 ELSE 0 END), 0) AS error,
    COALESCE(SUM(CASE WHEN lh.status = 'timeout' THEN 1 ELSE 0 END), 0) AS timeout
FROM checks c
LEFT JOIN latest_history lh ON c.id = lh.check_id
WHERE c.enabled = ?
`

	// For SQLite, DISTINCT ON is not supported → fallback to correlated subquery
	if h.storage.DB().Dialector.Name() == "sqlite" {
		query = `
SELECT
    COALESCE(SUM(CASE WHEN ch.status = 'up' THEN 1 ELSE 0 END), 0) AS up,
    COALESCE(SUM(CASE WHEN ch.status = 'down' THEN 1 ELSE 0 END), 0) AS down,
    COALESCE(SUM(CASE WHEN ch.status = 'error' THEN 1 ELSE 0 END), 0) AS error,
    COALESCE(SUM(CASE WHEN ch.status = 'timeout' THEN 1 ELSE 0 END), 0) AS timeout
FROM checks c
LEFT JOIN check_history ch ON c.id = ch.check_id
AND ch.id = (
    SELECT MAX(ch2.id)
    FROM check_history ch2
    WHERE ch2.check_id = c.id
)
WHERE c.enabled = ?
`
	}

	var statusCounts struct {
		Up      int64 `json:"up"`
		Down    int64 `json:"down"`
		Error   int64 `json:"error"`
		Timeout int64 `json:"timeout"`
	}

	if err := db.Raw(query, true).Scan(&statusCounts).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to compute status distribution", err))
		return
	}

	response := gin.H{
		"total":    total,
		"enabled":  enabledCount,
		"disabled": disabledCount,
		"status": gin.H{
			"up":      statusCounts.Up,
			"down":    statusCounts.Down,
			"error":   statusCounts.Error,
			"timeout": statusCounts.Timeout,
		},
	}

	c.JSON(http.StatusOK, types.SuccessResponse(response))
}

// resolveCheckStatus retrieves the most recent status and timestamp for a check.
//
// This helper method queries check_history to determine:
//   - Current status (from latest record)
//   - Last checked timestamp
//
// If no history exists, returns empty status and nil timestamp.
func (h *Handler) resolveCheckStatus(checkID int64) (string, *time.Time) {
	var history storage.CheckHistory
	err := h.storage.DB().
		Where("check_id = ?", checkID).
		Order("checked_at DESC").
		First(&history).Error

	if err != nil {
		return storage.CheckStatusDown, nil
	}

	return history.Status, &history.CheckedAt
}

// translateCheckDBError converts database errors to appropriate HTTP errors.
// This helper is used by both Create and Update handlers for consistent error handling.
func translateCheckDBError(err error) *types.ErrorWithContext {
	if err == nil {
		return nil
	}

	// Extract error message
	errMsg := err.Error()

	// Conflict errors (409)
	if strings.Contains(errMsg, "name already exists") ||
		strings.Contains(errMsg, "UNIQUE constraint failed") ||
		strings.Contains(errMsg, "duplicate key") {
		return types.ConflictError("check with this name already exists")
	}

	// Validation errors (400)
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "unsupported") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "too long") ||
		strings.Contains(errMsg, "too short") ||
		strings.Contains(errMsg, "must be between") ||
		strings.Contains(errMsg, "must be less than") {
		return types.ValidationError(errMsg)
	}

	// Not found errors (404)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return types.NotFoundError("check")
	}

	// Default to internal server error (500)
	return types.InternalError("database operation failed", err)
}

// parseConfigToRawMessage converts a JSON config string to json.RawMessage.
// If the string is empty or invalid, returns an empty JSON object {}.
func parseConfigToRawMessage(configStr string) json.RawMessage {
	if configStr == "" || configStr == "{}" {
		return json.RawMessage("{}")
	}

	// Parse and re-marshal to ensure valid JSON format
	var tmp interface{}
	if err := json.Unmarshal([]byte(configStr), &tmp); err != nil {
		return json.RawMessage("{}")
	}

	raw, err := json.Marshal(tmp)
	if err != nil {
		return json.RawMessage("{}")
	}

	return raw
}

// Create handles POST /api/v1/checks
//
// Creates a new monitoring check with the provided configuration.
// The check name must be unique within the system. Upon creation,
// the check is immediately validated against business rules.
//
// Required fields:
//   - name: Unique identifier (max 100 chars)
//   - type: Monitoring type ("http" or "ping")
//   - target: URL or host to monitor
//   - interval_seconds: Execution frequency (5-86400 seconds)
//
// Optional fields:
//   - enabled (default: true)
//   - timeout_seconds (default: 30)
//   - config: Type-specific JSON configuration (can be sent as object or string)
//
// Validation:
//   - Name uniqueness enforced via database constraint
//   - Target format validated per check type (URL for HTTP, hostname/IP for ping)
//   - Config JSON validated and normalized with defaults
//
// Returns:
//   - 201 Created with check details
//   - 400 Bad Request for invalid input or validation failure
//   - 409 Conflict if check name already exists
//   - 500 Internal Server Error on storage failure
func (h *Handler) Create(c *gin.Context) {
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		types.AbortWithError(c, types.ValidationError(err.Error()))
		return
	}

	// Validate required fields
	if req.Name == nil || req.Type == nil || req.Target == nil || req.IntervalSeconds == nil {
		types.AbortWithError(c, types.ValidationError(
			"name, type, target, and interval_seconds are required",
		))
		return
	}

	// Ensure name uniqueness
	var existing storage.Check
	err := h.storage.DB().Where("name = ?", *req.Name).First(&existing).Error
	switch {
	case err == nil:
		types.AbortWithError(c, types.ConflictError("check with this name already exists"))
		return
	case !errors.Is(err, gorm.ErrRecordNotFound):
		types.AbortWithError(c, types.InternalError("failed to check check name uniqueness", err))
		return
	}

	// Convert config from json.RawMessage to string
	configStr := "{}"
	if req.Config != nil {
		// Validate that it's valid JSON
		var tmp interface{}
		if err := json.Unmarshal(*req.Config, &tmp); err != nil {
			types.AbortWithError(c, types.ValidationError("config must be valid JSON"))
			return
		}
		configStr = string(*req.Config)
	}

	// Build check entity
	now := time.Now()
	check := &storage.Check{
		Name:            *req.Name,
		Type:            *req.Type,
		Target:          *req.Target,
		Enabled:         true,
		IntervalSeconds: *req.IntervalSeconds,
		TimeoutSeconds:  5, // Default
		Config:          configStr,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Apply optional fields
	if req.Enabled != nil {
		check.Enabled = *req.Enabled
	}
	if req.TimeoutSeconds != nil {
		check.TimeoutSeconds = *req.TimeoutSeconds
	}

	// Use engine to persist and schedule (if running)
	if h.engine != nil && h.engine.IsRunning() {
		// Engine handles both DB save and scheduling
		if err := h.engine.AddCheck(check); err != nil {
			// Translate common errors to HTTP status codes
			if strings.Contains(err.Error(), "name already exists") ||
				strings.Contains(err.Error(), "UNIQUE constraint failed") {
				types.AbortWithError(c, types.ConflictError("check with this name already exists"))
			} else if strings.Contains(err.Error(), "invalid") ||
				strings.Contains(err.Error(), "unsupported") {
				types.AbortWithError(c, types.ValidationError(err.Error()))
			} else {
				types.AbortWithError(c, types.InternalError("failed to create and schedule check", err))
			}
			return
		}
	} else {
		// Engine not running: save directly to DB (e.g., during startup or tests)
		if err := h.storage.DB().Create(check).Error; err != nil {
			if strings.Contains(err.Error(), "name already exists") ||
				strings.Contains(err.Error(), "UNIQUE constraint failed") {
				types.AbortWithError(c, types.ConflictError("check with this name already exists"))
			} else {
				types.AbortWithError(c, types.InternalError("failed to create check", err))
			}
			return
		}
	}

	// Resolve runtime status from history (likely empty at creation)
	status, lastChecked := h.resolveCheckStatus(check.ID)

	response := CheckResponse{
		ID:              check.ID,
		Name:            check.Name,
		Type:            check.Type,
		Target:          check.Target,
		Config:          parseConfigToRawMessage(check.Config),
		Enabled:         check.Enabled,
		IntervalSeconds: check.IntervalSeconds,
		TimeoutSeconds:  check.TimeoutSeconds,
		Status:          status,
		LastCheckedAt:   lastChecked,
		CreatedAt:       check.CreatedAt,
		UpdatedAt:       check.UpdatedAt,
	}

	c.JSON(http.StatusCreated, types.SuccessResponse(response))
}

// Get handles GET /api/v1/checks/:id
//
// Retrieves detailed information about a single check by its ID,
// including its current runtime status and last checked timestamp.
//
// Security considerations:
//   - Returns 404 for non-existent checks (no enumeration vulnerability)
//   - No sensitive data exposed
//
// Returns:
//   - 200 OK with check details
//   - 400 Bad Request for invalid check ID format
//   - 404 Not Found if check does not exist
//   - 500 Internal Server Error on storage failure
func (h *Handler) Get(c *gin.Context) {
	// Parse and validate check ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		types.AbortWithError(c, types.ValidationError("invalid check ID"))
		return
	}

	// Retrieve check from database
	var check storage.Check
	if err := h.storage.DB().Where("id = ?", id).First(&check).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			types.AbortWithError(c, types.NotFoundError("check"))
		} else {
			types.AbortWithError(c, types.InternalError("failed to retrieve check", err))
		}
		return
	}

	// Resolve runtime status from history
	status, lastChecked := h.resolveCheckStatus(check.ID)

	response := CheckResponse{
		ID:              check.ID,
		Name:            check.Name,
		Type:            check.Type,
		Target:          check.Target,
		Config:          parseConfigToRawMessage(check.Config),
		Enabled:         check.Enabled,
		IntervalSeconds: check.IntervalSeconds,
		TimeoutSeconds:  check.TimeoutSeconds,
		Status:          status,
		LastCheckedAt:   lastChecked,
		CreatedAt:       check.CreatedAt,
		UpdatedAt:       check.UpdatedAt,
	}

	c.JSON(http.StatusOK, types.SuccessResponse(response))
}

// Update handles PATCH /api/v1/checks/:id
//
// Partially updates a check configuration. Only provided fields are modified;
// all others remain unchanged. The update triggers full validation to ensure
// the modified check remains valid.
//
// Updatable fields:
//   - name (must remain unique)
//   - target (validated per check type)
//   - config (validated and normalized)
//   - interval_seconds (5-86400)
//   - timeout_seconds (1-300, must be < interval)
//   - enabled
//
// Validation:
//   - All modified fields re-validated via BeforeUpdate hook
//   - Name uniqueness enforced (cannot conflict with other checks)
//   - Business rules enforced (timeout < interval, valid target format, etc.)
//
// Returns:
//   - 200 OK with updated check details
//   - 400 Bad Request for invalid input or empty payload
//   - 404 Not Found if check does not exist
//   - 409 Conflict if name conflicts with existing check
//   - 500 Internal Server Error on storage failure
func (h *Handler) Update(c *gin.Context) {
	// Parse and validate check ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		types.AbortWithError(c, types.ValidationError("invalid check ID"))
		return
	}

	// Bind request body
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		types.AbortWithError(c, types.ValidationError(err.Error()))
		return
	}

	// Reject empty PATCH payload
	if req.Name == nil && req.Type == nil && req.Target == nil &&
		req.Config == nil && req.IntervalSeconds == nil &&
		req.TimeoutSeconds == nil && req.Enabled == nil {
		types.AbortWithError(c, types.ValidationError("no fields to update"))
		return
	}

	// Retrieve existing check
	var check storage.Check
	if err := h.storage.DB().Where("id = ?", id).First(&check).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			types.AbortWithError(c, types.NotFoundError("check"))
		} else {
			types.AbortWithError(c, types.InternalError("failed to retrieve check", err))
		}
		return
	}

	// Apply only provided fields
	if req.Name != nil {
		check.Name = *req.Name
	}
	if req.Type != nil {
		check.Type = *req.Type
	}
	if req.Target != nil {
		check.Target = *req.Target
	}
	if req.IntervalSeconds != nil {
		check.IntervalSeconds = *req.IntervalSeconds
	}
	if req.TimeoutSeconds != nil {
		check.TimeoutSeconds = *req.TimeoutSeconds
	}
	if req.Enabled != nil {
		check.Enabled = *req.Enabled
	}

	// Convert config from json.RawMessage to string (if provided)
	if req.Config != nil {
		// Validate that it's valid JSON
		var tmp interface{}
		if err := json.Unmarshal(*req.Config, &tmp); err != nil {
			types.AbortWithError(c, types.ValidationError("config must be valid JSON"))
			return
		}
		check.Config = string(*req.Config)
	}

	// Save with validation (triggers BeforeUpdate hook)
	if err := h.storage.DB().Save(&check).Error; err != nil {
		types.AbortWithError(c, translateCheckDBError(err))
		return
	}

	// Resolve updated runtime status
	status, lastChecked := h.resolveCheckStatus(check.ID)
	response := CheckResponse{
		ID:              check.ID,
		Name:            check.Name,
		Type:            check.Type,
		Target:          check.Target,
		Config:          parseConfigToRawMessage(check.Config),
		Enabled:         check.Enabled,
		IntervalSeconds: check.IntervalSeconds,
		TimeoutSeconds:  check.TimeoutSeconds,
		Status:          status,
		LastCheckedAt:   lastChecked,
		CreatedAt:       check.CreatedAt,
		UpdatedAt:       check.UpdatedAt,
	}

	c.JSON(http.StatusOK, types.SuccessResponse(response))
}

// Delete handles DELETE /api/v1/checks/:id
//
// Permanently deletes a check and all its associated history records.
// Database foreign key constraints with CASCADE ensure related check_history
// records are automatically removed.
//
// Security considerations:
//   - Proper 404 semantics (no enumeration)
//   - Atomic operation (all or nothing)
//
// Returns:
//   - 200 OK on successful deletion
//   - 400 Bad Request for invalid check ID format
//   - 404 Not Found if check does not exist
//   - 500 Internal Server Error on storage failure
func (h *Handler) Delete(c *gin.Context) {
	// Parse and validate check ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		types.AbortWithError(c, types.ValidationError("invalid check ID"))
		return
	}

	// Verify existence for proper 404 semantics
	var count int64
	if err := h.storage.DB().
		Model(&storage.Check{}).
		Where("id = ?", id).
		Limit(1).
		Count(&count).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to check check existence", err))
		return
	}

	if count == 0 {
		types.AbortWithError(c, types.NotFoundError("check"))
		return
	}

	// Delete check (CASCADE removes related check_history records)
	if err := h.storage.DB().Delete(&storage.Check{}, id).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to delete check", err))
		return
	}

	// Remove from scheduler if engine is running
	if h.engine != nil && h.engine.IsRunning() {
		if err := h.engine.RemoveCheckFromScheduler(id); err != nil {
			log.Warn().Int64("check_id", id).Err(err).Msg("Failed to remove check from scheduler")
		}
	}

	c.JSON(http.StatusOK, types.SuccessResponse(gin.H{
		"message": "check deleted successfully",
	}))
}

// History handles GET /api/v1/checks/:id/history
//
// Returns a paginated list of historical execution records for a specific check,
// ordered by execution time in descending order (newest first).
//
// Query parameters:
//   - page (default: 1)
//   - page_size (default: 20, max: 500)
//   - short (optional boolean): if true, returns compact format [id, status]
//
// Compact format (when short=true):
//   - Response data is an array of arrays: [[id1, status1], [id2, status2], ...]
//   - Only two fields are included: ID and Status
//   - Pagination still applies
//
// Returns:
//   - 200 OK with paginated history records (full or compact)
//   - 400 Bad Request for invalid check ID or pagination parameters
//   - 404 Not Found if check does not exist
//   - 500 Internal Server Error on storage failure
func (h *Handler) History(c *gin.Context) {
	// Parse and validate check ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		types.AbortWithError(c, types.ValidationError("invalid check ID"))
		return
	}

	// Verify check exists
	var count int64
	if err := h.storage.DB().
		Model(&storage.Check{}).
		Where("id = ?", id).
		Limit(1).
		Count(&count).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to check check existence", err))
		return
	}

	if count == 0 {
		types.AbortWithError(c, types.NotFoundError("check"))
		return
	}

	// Bind pagination and short query parameters
	var pagination types.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		types.AbortWithError(c, types.ValidationError(err.Error()))
		return
	}

	page := pagination.Page
	pageSize := pagination.PageSize
	offset := (page - 1) * pageSize

	// Count total history records
	var total int64
	if err := h.storage.DB().
		Model(&storage.CheckHistory{}).
		Where("check_id = ?", id).
		Count(&total).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to count check history", err))
		return
	}

	if pagination.Short {
		// Compact mode: fetch only ID and Status
		type CompactRecord struct {
			ID     int64  `gorm:"column:id"`
			Status string `gorm:"column:status"`
		}

		var compactRecords []CompactRecord
		if err := h.storage.DB().
			Select("id, status").
			Model(&storage.CheckHistory{}).
			Where("check_id = ?", id).
			Order("checked_at DESC").
			Limit(pageSize).
			Offset(offset).
			Find(&compactRecords).Error; err != nil {
			types.AbortWithError(c, types.InternalError("failed to retrieve compact check history", err))
			return
		}

		// Convert to array of arrays: [[id, status], ...]
		compactResponse := make([][]interface{}, len(compactRecords))
		for i, r := range compactRecords {
			compactResponse[i] = []interface{}{r.ID, r.Status}
		}

		totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
		paginationResponse := &types.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		}

		c.JSON(http.StatusOK, types.SuccessResponseWithPagination(compactResponse, paginationResponse))
		return
	}

	// Full mode: fetch all fields
	var histories []storage.CheckHistory
	if err := h.storage.DB().
		Where("check_id = ?", id).
		Order("checked_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&histories).Error; err != nil {
		types.AbortWithError(c, types.InternalError("failed to retrieve check history", err))
		return
	}

	// Map to response DTOs
	responses := make([]CheckHistoryResponse, 0, len(histories))
	for _, history := range histories {
		responses = append(responses, CheckHistoryResponse{
			ID:             history.ID,
			CheckID:        history.CheckID,
			Status:         history.Status,
			ResponseTimeMs: history.ResponseTimeMs,
			StatusCode:     history.StatusCode,
			ErrorMessage:   history.ErrorMessage,
			CheckedAt:      history.CheckedAt,
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	paginationResponse := &types.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, types.SuccessResponseWithPagination(responses, paginationResponse))
}

// GetHistoryByID handles GET /api/v1/checks/:id/history/:history_id
//
// Retrieves a single historical execution record for a specific check by its history ID.
// This endpoint provides direct access to a specific observation without pagination overhead.
//
// Security considerations:
//   - Validates that history record belongs to the specified check
//   - Returns 404 for non-existent records (no enumeration)
//
// Use cases:
//   - Deep dive into specific check execution
//   - Debugging failed checks
//   - Retrieving detailed error messages
//
// Returns:
//   - 200 OK with the specific history record
//   - 400 Bad Request for invalid ID formats
//   - 404 Not Found if check or history record does not exist
//   - 500 Internal Server Error on storage failure
func (h *Handler) GetHistoryByID(c *gin.Context) {
	// Parse and validate check ID
	checkID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		types.AbortWithError(c, types.ValidationError("invalid check ID"))
		return
	}

	// Parse and validate history ID
	historyID, err := strconv.ParseInt(c.Param("history_id"), 10, 64)
	if err != nil {
		types.AbortWithError(c, types.ValidationError("invalid history ID"))
		return
	}

	// Retrieve specific history record with ownership validation
	var history storage.CheckHistory
	err = h.storage.DB().
		Where("id = ? AND check_id = ?", historyID, checkID).
		First(&history).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			types.AbortWithError(c, types.NotFoundError("check history record"))
		} else {
			types.AbortWithError(c, types.InternalError("failed to retrieve check history", err))
		}
		return
	}

	response := CheckHistoryResponse{
		ID:             history.ID,
		CheckID:        history.CheckID,
		Status:         history.Status,
		ResponseTimeMs: history.ResponseTimeMs,
		StatusCode:     history.StatusCode,
		ErrorMessage:   history.ErrorMessage,
		CheckedAt:      history.CheckedAt,
	}

	c.JSON(http.StatusOK, types.SuccessResponse(response))
}
