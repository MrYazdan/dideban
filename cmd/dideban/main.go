// Package main provides the entry point for the Dideban monitoring system.
//
// Dideban (دیدبان) is a lightweight, fast, and self-hosted monitoring system
// built for private infrastructures, VPCs, and production-grade web applications.
package main

import (
	"fmt"

	"dideban/internal/config"

	"github.com/rs/zerolog/log"
)

// Version information set during build time
var (
	Version   = "0.1.0"
	GitCommit = "unknown"
	BuildTime = "unknown"
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

	fmt.Print(cfg)
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
