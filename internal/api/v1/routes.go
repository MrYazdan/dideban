package v1

import (
	"dideban/internal/core"
	"dideban/internal/storage"

	"dideban/internal/api/v1/agents"
	"dideban/internal/api/v1/checks"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures API routes.
func SetupRoutes(routerGroup *gin.RouterGroup, engine *core.Engine, storage *storage.Storage) {
	// Initialize handlers
	checksHandler := checks.NewHandler(storage, engine)
	agentsHandler := agents.NewHandler(storage)

	// Checks management
	checksGroup := routerGroup.Group("/checks")
	{
		checksGroup.GET("", checksHandler.List)
		checksGroup.POST("", checksHandler.Create)
		checksGroup.GET("/stats", checksHandler.Stats)
		checksGroup.GET("/:id", checksHandler.Get)
		checksGroup.PATCH("/:id", checksHandler.Update)
		checksGroup.DELETE("/:id", checksHandler.Delete)
		checksGroup.GET("/:id/history", checksHandler.History)
		checksGroup.GET("/:id/history/:history_id", checksHandler.GetHistoryByID)
	}

	// Agents management
	agentsGroup := routerGroup.Group("/agents")
	{
		agentsGroup.GET("", agentsHandler.List)
		agentsGroup.POST("", agentsHandler.Create)
		agentsGroup.GET("/stats", agentsHandler.Stats)
		agentsGroup.GET("/:id", agentsHandler.Get)
		agentsGroup.PATCH("/:id", agentsHandler.Update)
		agentsGroup.DELETE("/:id", agentsHandler.Delete)
		agentsGroup.POST("/:id/regenerate", agentsHandler.RegenerateToken)
		agentsGroup.GET("/:id/history", agentsHandler.History)
		agentsGroup.POST("/:id/history", agentsHandler.AgentAuthMiddleware(), agentsHandler.CreateHistory)
		agentsGroup.GET("/:id/history/:history_id", agentsHandler.GetHistoryByID)
	}
}
