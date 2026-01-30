// Package api provides public endpoints for system health and connectivity.
//
// These endpoints are designed to be lightweight, fast, and reliable for external monitoring
// systems (e.g., load balancers, uptime monitors, observability tools).
package api

import (
	"context"
	"net/http"
	"time"

	"dideban/internal/core"
	"dideban/internal/storage"

	"github.com/gin-gonic/gin"
)

// Handler manages public endpoints.
//
// It provides essential system-level information,
// making it suitable for health checks, liveness probes, and basic diagnostics.
type Handler struct {
	engine    *core.Engine
	storage   *storage.Storage
	startTime time.Time
}

// NewHandler initializes a new public API handler.
//
// Parameters:
//   - engine: Core monitoring engine (maybe nil in test environments)
//   - storage: Database storage layer (maybe nil in test environments)
//
// Returns a fully initialized handler ready for HTTP routing.
func NewHandler(engine *core.Engine, storage *storage.Storage) *Handler {
	return &Handler{
		engine:    engine,
		storage:   storage,
		startTime: time.Now(),
	}
}

// Ping handles GET /ping
//
// A lightweight endpoint for basic connectivity verification.
// Returns a simple "pong" response with minimal overhead.
//
// Use cases:
//   - Load balancer health checks
//   - Network connectivity tests
//   - Quick service availability confirmation
//
// Response:
//   - 200 OK with {"message": "pong"}
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Health handles GET /health
//
// Provides comprehensive system health status including database, engine,
// and scheduler components. Designed for Kubernetes liveness/readiness probes
// and observability dashboards.
//
// The endpoint aggregates real-time metrics from all core subsystems:
//   - Database: connection and query latency
//   - Engine: running state and active checks
//   - Scheduler: worker activity and job count
//
// Overall status is "healthy" only if all components report healthy;
// otherwise, it returns "degraded".
//
// Response:
//   - 200 OK with detailed health report
//   - All fields are always present (never null)
func (h *Handler) Health(c *gin.Context) {
	ctx := c.Request.Context()

	// Calculate system uptime
	uptime := time.Since(h.startTime)

	// Perform component health checks
	dbStatus, dbResponseTime := h.checkDatabaseHealth(ctx)
	engineStatus, checksRunning := h.checkEngineHealth()
	schedulerStatus, workersActive := h.checkSchedulerHealth()

	// Determine overall system status
	overallStatus := "healthy"
	if dbStatus != "healthy" || engineStatus != "healthy" || schedulerStatus != "healthy" {
		overallStatus = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    overallStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    uptime.String(),
		"version":   "0.1.0", // TODO: inject from build info
		"components": gin.H{
			"database": gin.H{
				"status":           dbStatus,
				"response_time_ms": dbResponseTime,
			},
			"engine": gin.H{
				"status":         engineStatus,
				"checks_running": checksRunning,
			},
			"scheduler": gin.H{
				"status":         schedulerStatus,
				"workers_active": workersActive,
			},
		},
	})
}

// checkDatabaseHealth performs a database connectivity and latency test.
//
// It attempts to ping the underlying SQL database and measures response time.
// Returns:
//   - status: "healthy" or "unhealthy"
//   - response_time_ms: round-trip time in milliseconds (0 if unhealthy)
func (h *Handler) checkDatabaseHealth(ctx context.Context) (string, int64) {
	// Handle nil storage (e.g., during testing)
	if h.storage == nil {
		return "unhealthy", 0
	}

	start := time.Now()
	sqlDB, err := h.storage.DB().DB()
	if err != nil {
		return "unhealthy", time.Since(start).Milliseconds()
	}

	err = sqlDB.PingContext(ctx)
	responseTime := time.Since(start).Milliseconds()
	if err != nil {
		return "unhealthy", responseTime
	}

	return "healthy", responseTime
}

// checkEngineHealth verifies the monitoring engine's operational status.
//
// It checks if the engine is running and counts enabled checks as a proxy
// for active monitoring workload.
//
// Returns:
//   - status: "healthy", "degraded", or "unhealthy"
//   - checks_running: number of enabled checks (0 if engine unavailable)
func (h *Handler) checkEngineHealth() (string, int) {
	if h.engine == nil {
		return "unhealthy", 0
	}

	if !h.engine.IsRunning() {
		return "unhealthy", 0
	}

	// Handle nil storage gracefully
	if h.storage == nil {
		return "degraded", 0
	}

	var checksCount int64
	if err := h.storage.DB().
		Model(&storage.Check{}).
		Where("enabled = ?", true).
		Count(&checksCount).Error; err != nil {
		return "degraded", 0
	}

	return "healthy", int(checksCount)
}

// checkSchedulerHealth assesses the background job scheduler's health.
//
// Since the scheduler is tightly coupled with the engine, this method
// reuses the engine's running state and enabled check count as indicators
// of active scheduled work.
//
// Returns:
//   - status: "healthy", "degraded", or "unhealthy"
//   - workers_active: estimated number of active scheduled jobs
func (h *Handler) checkSchedulerHealth() (string, int) {
	if h.engine == nil {
		return "unhealthy", 0
	}

	if !h.engine.IsRunning() {
		return "unhealthy", 0
	}

	// Handle nil storage gracefully
	if h.storage == nil {
		return "degraded", 0
	}

	var checksCount int64
	if err := h.storage.DB().
		Model(&storage.Check{}).
		Where("enabled = ?", true).
		Count(&checksCount).Error; err != nil {
		return "degraded", 0
	}

	return "healthy", int(checksCount)
}
