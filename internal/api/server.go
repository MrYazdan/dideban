// Package api provides HTTP API functionality for the Dideban monitoring system.
// This package implements a RESTFul API using Gin framework
//
// Example usage:
//
//	server := api.NewServer(cfg.Server, engine, logger)
//	err := server.Start()
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"dideban/internal/api/auth"
	"dideban/internal/config"
	"dideban/internal/core"
	"dideban/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Server represents the HTTP API server.
type Server struct {
	config      config.ServerConfig
	engine      *core.Engine
	storage     *storage.Storage
	router      *gin.Engine
	server      *http.Server
	authHandler *auth.Handler
}

// NewServer creates a new HTTP API server instance.
//
// Parameters:
//   - cfg: Server configuration containing address and timeout settings
//   - engine: Core monitoring engine instance
//   - storage: Storage instance for database operations
//
// Returns:
//   - *Server: Initialized server instance
func NewServer(cfg config.ServerConfig, engine *core.Engine, storage *storage.Storage) *Server {
	// Set Gin mode based on configuration
	gin.SetMode(gin.ReleaseMode)

	// Initialize auth handler with JWT configuration
	authHandler := auth.NewHandler(storage, []byte(cfg.JWT.Secret), cfg.JWT.TTL)

	server := &Server{
		config:      cfg,
		engine:      engine,
		storage:     storage,
		router:      gin.New(),
		authHandler: authHandler,
	}

	// Setup middleware and routes
	server.setupMiddleware()
	server.setupRoutes()

	// Create HTTP server
	server.server = &http.Server{
		Addr:         cfg.Addr,
		Handler:      server.router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return server
}

// Start starts the HTTP server and begins listening for requests.
//
// Returns:
//   - error: Any error that occurred during server startup
func (s *Server) Start() error {
	log.Info().Str("addr", s.config.Addr).Msg("Starting HTTP server")

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server.
//
// Parameters:
//   - ctx: Context for shutdown timeout
//
// Returns:
//   - error: Any error that occurred during shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server")

	// Stop auth handler cleanup goroutines
	if s.authHandler != nil {
		s.authHandler.Stop()
	}

	return s.server.Shutdown(ctx)
}

// setupMiddleware configures middleware for the Gin router.
func (s *Server) setupMiddleware() {
	// Request ID middleware (should be first)
	s.router.Use(RequestID())

	// Custom panic recovery middleware
	s.router.Use(PanicRecovery())

	// Request timeout middleware
	s.router.Use(TimeoutMiddleware(30 * time.Second))

	// Security headers
	s.router.Use(SecurityHeaders())

	// Content type validation
	s.router.Use(ContentType())

	// Rate limiting
	s.router.Use(RateLimit())

	// Custom logger middleware
	s.router.Use(LoggerMiddleware())

	// Error handling middleware
	s.router.Use(ErrorHandler())
}
