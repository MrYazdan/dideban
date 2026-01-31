package api

import (
	v1 "dideban/internal/api/v1"

	"github.com/gin-gonic/gin"
)

// setupRoutes configures API routes.
func (s *Server) setupRoutes() {
	// Initialize handlers
	baseHandler := NewHandler(s.engine, s.storage)

	// Base api router group
	apiGroup := s.router.Group("/api")

	// Base endpoints (no authentication required)
	apiGroup.GET("/ping", baseHandler.Ping)
	apiGroup.GET("/health", baseHandler.Health)

	// Authentication endpoints
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.POST("/login", s.authHandler.Login)
		authGroup.POST("/logout", s.authHandler.Logout)
		authGroup.GET("/me", s.authHandler.Me)
		authGroup.POST("/refresh", s.authHandler.Refresh)
	}

	// Initialize auth middleware for v1 endpoints
	//authMiddleware := v1.NewAuthMiddleware(s.authHandler.GetTokenManager())

	// API v1 routes (protected with authentication)
	v1Group := apiGroup.Group("/v1")
	//v1Group.Use(authMiddleware.RequireAuth())

	v1.SetupRoutes(v1Group, s.engine, s.storage)

	// Serve static files (frontend) from root
	s.router.Static("/static", "./web/static")
	s.router.StaticFile("/", "./web/index.html")
	s.router.NoRoute(func(c *gin.Context) {
		// Serve index.html for SPA routes
		c.File("./web/index.html")
	})
}
