package config

import (
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Package-level constants for performance optimization
var (
	validLogLevels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
)

// validateConfig validates the configuration and returns an error if invalid.
func validateConfig(c *Config) error {
	for _, validate := range []func() error{
		func() error { return validateServerConfig(c.Server) },
		func() error { return validateStorageConfig(c.Storage) },
		func() error { return validateAlertConfig(c.Alert) },
		func() error { return validateSchedulerConfig(c.Scheduler) },
		func() error { return validateLogConfig(c.Log) },
	} {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}

// validateServerConfig validates server configuration.
func validateServerConfig(s ServerConfig) error {
	if s.Addr == "" {
		return fmt.Errorf("server.addr cannot be empty")
	}

	// Validate address format
	host, portStr, err := net.SplitHostPort(s.Addr)
	if err != nil {
		return fmt.Errorf("server.addr invalid format: %w", err)
	}

	// Validate port range
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("server.addr invalid port: %w", err)
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("server.addr port out of range (1-65535)")
		}
	}

	// Validate host if specified
	if host != "" && host != "0.0.0.0" && host != "localhost" {
		if ip := net.ParseIP(host); ip == nil {
			// Try to resolve hostname
			if _, err := net.LookupHost(host); err != nil {
				return fmt.Errorf("server.addr invalid host: %s", host)
			}
		}
	}

	// Validate timeouts
	if s.ReadTimeout <= 0 {
		return fmt.Errorf("server.read_timeout must be greater than 0")
	}
	if s.WriteTimeout <= 0 {
		return fmt.Errorf("server.write_timeout must be greater than 0")
	}
	if s.IdleTimeout <= 0 {
		return fmt.Errorf("server.idle_timeout must be greater than 0")
	}

	// Validate timeout ranges (reasonable limits)
	if s.ReadTimeout > 5*time.Minute {
		return fmt.Errorf("server.read_timeout too large (max 5m)")
	}
	if s.WriteTimeout > 5*time.Minute {
		return fmt.Errorf("server.write_timeout too large (max 5m)")
	}
	if s.IdleTimeout > 30*time.Minute {
		return fmt.Errorf("server.idle_timeout too large (max 30m)")
	}

	// Minimum timeout validation
	if s.ReadTimeout < time.Second {
		return fmt.Errorf("server.read_timeout too small (min 1s)")
	}
	if s.WriteTimeout < time.Second {
		return fmt.Errorf("server.write_timeout too small (min 1s)")
	}

	// Validate JWT configuration
	if err := validateJWTConfig(s.JWT); err != nil {
		return fmt.Errorf("server.jwt: %w", err)
	}

	return nil
}

// validateStorageConfig validates storage configuration.
func validateStorageConfig(s StorageConfig) error {
	if s.Path == "" {
		return fmt.Errorf("storage.path cannot be empty")
	}

	// Validate path format (basic check)
	if strings.Contains(s.Path, "..") {
		return fmt.Errorf("storage.path cannot contain '..' for security")
	}

	// Validate connection pool settings
	if s.MaxOpenConns <= 0 {
		return fmt.Errorf("storage.max_open_conns must be greater than 0")
	}
	if s.MaxIdleConns < 0 {
		return fmt.Errorf("storage.max_idle_conns cannot be negative")
	}
	if s.MaxIdleConns > s.MaxOpenConns {
		return fmt.Errorf("storage.max_idle_conns cannot be greater than max_open_conns")
	}
	if s.ConnMaxLifetime <= 0 {
		return fmt.Errorf("storage.conn_max_lifetime must be greater than 0")
	}

	// Validate reasonable limits
	if s.MaxOpenConns > 1000 {
		return fmt.Errorf("storage.max_open_conns too large (max 1000)")
	}
	if s.ConnMaxLifetime > 24*time.Hour {
		return fmt.Errorf("storage.conn_max_lifetime too large (max 24h)")
	}
	if s.ConnMaxLifetime < time.Minute {
		return fmt.Errorf("storage.conn_max_lifetime too small (min 1m)")
	}

	return nil
}

// validateAlertConfig validates alert configuration.
func validateAlertConfig(a AlertConfig) error {
	// Validate Telegram configuration
	if err := validateBotConfig(a.Telegram, "telegram"); err != nil {
		return err
	}

	// Validate Bale configuration
	if err := validateBotConfig(a.Bale, "bale"); err != nil {
		return err
	}

	return nil
}

// validateBotConfig validates a bot configuration.
func validateBotConfig(bot BotConfig, botType string) error {
	if bot.Enabled {
		if bot.Token == "" {
			return fmt.Errorf("alert.%s.token is required when %s is enabled", botType, botType)
		}
		if bot.ChatID == "" {
			return fmt.Errorf("alert.%s.chat_id is required when %s is enabled", botType, botType)
		}

		// Validate token format (basic security check)
		if len(bot.Token) < 10 {
			return fmt.Errorf("alert.%s.token too short (min 10 chars)", botType)
		}
		if len(bot.Token) > 200 {
			return fmt.Errorf("alert.%s.token too long (max 200 chars)", botType)
		}

		// Validate chat ID format (basic check)
		if len(bot.ChatID) < 1 {
			return fmt.Errorf("alert.%s.chat_id too short", botType)
		}
		if len(bot.ChatID) > 50 {
			return fmt.Errorf("alert.%s.chat_id too long (max 50 chars)", botType)
		}

		// Validate timeout
		if bot.Timeout <= 0 {
			return fmt.Errorf("alert.%s.timeout must be greater than 0", botType)
		}
		if bot.Timeout > 2*time.Minute {
			return fmt.Errorf("alert.%s.timeout too large (max 2m)", botType)
		}
		if bot.Timeout < time.Second {
			return fmt.Errorf("alert.%s.timeout too small (min 1s)", botType)
		}
	}
	return nil
}

// validateSchedulerConfig validates scheduler configuration.
func validateSchedulerConfig(s SchedulerConfig) error {
	if s.WorkerCount <= 0 {
		return fmt.Errorf("scheduler.worker_count must be greater than 0")
	}
	if s.WorkerCount > 1000 {
		return fmt.Errorf("scheduler.worker_count too large (max 1000)")
	}

	if s.DefaultInterval <= 0 {
		return fmt.Errorf("scheduler.default_interval must be greater than 0")
	}
	if s.DefaultInterval < 5*time.Second {
		return fmt.Errorf("scheduler.default_interval too small (min 5s)")
	}
	if s.DefaultInterval > 24*time.Hour {
		return fmt.Errorf("scheduler.default_interval too large (max 24h)")
	}

	if s.MaxRetries < 0 {
		return fmt.Errorf("scheduler.max_retries cannot be negative")
	}
	if s.MaxRetries > 10 {
		return fmt.Errorf("scheduler.max_retries too large (max 10)")
	}

	return nil
}

// validateLogConfig validates log configuration.
func validateLogConfig(l LogConfig) error {
	if !slices.Contains(validLogLevels, strings.ToLower(l.Level)) {
		return fmt.Errorf("log.level must be one of: debug, info, warn, error, fatal, panic")
	}
	return nil
}

// validateJWTConfig validates JWT configuration.
func validateJWTConfig(j JWTConfig) error {
	// Validate JWT secret
	if j.Secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	// Security: JWT secret should be strong enough
	if len(j.Secret) < 32 {
		return fmt.Errorf("secret too short (minimum 32 characters for security)")
	}

	// Warn about default secret (security issue)
	if j.Secret == "your-secret-key-change-this-in-production" {
		return fmt.Errorf("default secret detected - change this in production for security")
	}

	// Validate TTL
	if j.TTL <= 0 {
		return fmt.Errorf("ttl must be greater than 0")
	}

	// Reasonable TTL limits
	if j.TTL < 5*time.Minute {
		return fmt.Errorf("ttl too small (minimum 5 minutes)")
	}

	if j.TTL > 30*24*time.Hour { // 30 days
		return fmt.Errorf("ttl too large (maximum 30 days)")
	}

	return nil
}
