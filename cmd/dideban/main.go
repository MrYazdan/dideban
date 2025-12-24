// Package main provides the entry point for the Dideban monitoring system.
//
// Dideban (دیدبان) is a lightweight, fast, and self-hosted monitoring system
// built for private infrastructures, VPCs, and production-grade web applications.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"dideban/internal/config"
	"dideban/internal/logger"

	"github.com/rs/zerolog/log"
)

// main is the entry point of the Dideban monitoring system.
//
// The startup sequence is as follows:
//  1. Load configuration
//  2. Initialize logger
//  3. Setup graceful shutdown handling
//  4. Start the main server
func main() {
	// Load application configuration (fails fast on error)
	cfg := loadConfig()

	// Initialize structured logger based on configuration
	logger.Init(cfg)

	// Log basic startup information for observability
	logStartup(cfg)

	// Create root context used across the entire application lifecycle.
	// This context is cancelled on shutdown signals (SIGINT, SIGTERM, SIGQUIT).
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register OS signal handlers for graceful shutdown
	setupSignalHandlers(cancel)

	log.Info().Msg("Application shutdown complete")
}

// setupSignalHandlers configures OS signal handling
// to enable graceful shutdown of the application.
//
// Upon receiving a shutdown signal, the provided cancel function is invoked,
// which propagates cancellation through the entire application via context.
func setupSignalHandlers(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // kill <pid>
		syscall.SIGQUIT, // quit
	)

	go func() {
		sig := <-sigChan

		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		// Trigger graceful shutdown
		cancel()
	}()
}

// loadConfig loads application configuration and terminates the program
// immediately if configuration cannot be loaded.
func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to load configuration")
	}
	return cfg
}

// logStartup logs essential startup metadata such as version,
// build information and configuration details.
func logStartup(cfg *config.Config) {
	log.Info().
		Str("version", "0.1.0").
		Str("log_level", cfg.Log.Level).
		Msg("🚀 Starting Dideban")
}
