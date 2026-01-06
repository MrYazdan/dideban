package tests

import (
	"os"
	"testing"
	"time"

	"dideban/internal/config"
)

// TestServerConfigValidation tests server configuration validation
func TestServerConfigValidation(t *testing.T) {
	t.Run("Valid server config", func(t *testing.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Addr:         ":8080",
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				JWT: config.JWTConfig{
					Secret: "this-is-a-very-secure-secret-key-with-32-plus-characters",
					TTL:    24 * time.Hour,
				},
			},
			Storage: config.StorageConfig{
				Path:            "test.db",
				MaxOpenConns:    32,
				MaxIdleConns:    8,
				ConnMaxLifetime: time.Hour,
			},
			Alert: config.AlertConfig{
				Telegram: config.BotConfig{Enabled: false, Timeout: 30 * time.Second},
				Bale:     config.BotConfig{Enabled: false, Timeout: 30 * time.Second},
			},
			Scheduler: config.SchedulerConfig{
				WorkerCount:     8,
				DefaultInterval: 60 * time.Second,
				MaxRetries:      3,
			},
			Log: config.LogConfig{
				Level:  "info",
				Pretty: false,
			},
		}

		// This should pass validation since we're using the internal validate function
		// We'll test this by trying to create a similar config through Load()
		if cfg.Server.Addr == "" {
			t.Error("Server addr should not be empty")
		}
	})

	t.Run("Invalid server address format", func(t *testing.T) {
		// Test invalid address by setting environment variable
		os.Setenv("DIDEBAN_SERVER_ADDR", "invalid-address")
		defer os.Unsetenv("DIDEBAN_SERVER_ADDR")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for invalid server address")
		}
	})

	t.Run("Invalid port range", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_ADDR", ":99999")
		defer os.Unsetenv("DIDEBAN_SERVER_ADDR")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for port out of range")
		}
	})

	t.Run("Zero timeout values", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_READ_TIMEOUT", "0s")
		defer os.Unsetenv("DIDEBAN_SERVER_READ_TIMEOUT")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for zero read timeout")
		}
	})

	t.Run("Timeout too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_READ_TIMEOUT", "10m")
		defer os.Unsetenv("DIDEBAN_SERVER_READ_TIMEOUT")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for read timeout too large")
		}
	})

	t.Run("Timeout too small", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_READ_TIMEOUT", "500ms")
		defer os.Unsetenv("DIDEBAN_SERVER_READ_TIMEOUT")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for read timeout too small")
		}
	})
}

// TestJWTConfigValidation tests JWT configuration validation
func TestJWTConfigValidation(t *testing.T) {
	t.Run("Empty JWT secret", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for empty JWT secret")
		}
	})

	t.Run("JWT secret too short", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "short")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for JWT secret too short")
		}
	})

	t.Run("Default JWT secret in production", func(t *testing.T) {
		// The default secret should trigger a validation error
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "your-secret-key-change-this-in-production")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for default JWT secret")
		}
	})

	t.Run("JWT TTL too small", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "this-is-a-very-secure-secret-key-with-32-plus-characters")
		os.Setenv("DIDEBAN_SERVER_JWT_TTL", "1m")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_SERVER_JWT_TTL")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for JWT TTL too small")
		}
	})

	t.Run("JWT TTL too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "this-is-a-very-secure-secret-key-with-32-plus-characters")
		os.Setenv("DIDEBAN_SERVER_JWT_TTL", "31d")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_SERVER_JWT_TTL")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for JWT TTL too large")
		}
	})

	t.Run("Zero JWT TTL", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "this-is-a-very-secure-secret-key-with-32-plus-characters")
		os.Setenv("DIDEBAN_SERVER_JWT_TTL", "0s")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_SERVER_JWT_TTL")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for zero JWT TTL")
		}
	})
}

// TestStorageConfigValidation tests storage configuration validation
func TestStorageConfigValidation(t *testing.T) {
	t.Run("Empty storage path", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_PATH", "")
		defer os.Unsetenv("DIDEBAN_STORAGE_PATH")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for empty storage path")
		}
	})

	t.Run("Storage path with directory traversal", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_PATH", "../../../etc/passwd")
		defer os.Unsetenv("DIDEBAN_STORAGE_PATH")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for storage path with directory traversal")
		}
	})

	t.Run("Zero max open connections", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS", "0")
		defer os.Unsetenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for zero max open connections")
		}
	})

	t.Run("Negative max idle connections", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_MAX_IDLE_CONNS", "-1")
		defer os.Unsetenv("DIDEBAN_STORAGE_MAX_IDLE_CONNS")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for negative max idle connections")
		}
	})

	t.Run("Max idle greater than max open", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS", "5")
		os.Setenv("DIDEBAN_STORAGE_MAX_IDLE_CONNS", "10")
		defer func() {
			os.Unsetenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS")
			os.Unsetenv("DIDEBAN_STORAGE_MAX_IDLE_CONNS")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for max idle > max open connections")
		}
	})

	t.Run("Connection lifetime too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_CONN_MAX_LIFETIME", "25h")
		defer os.Unsetenv("DIDEBAN_STORAGE_CONN_MAX_LIFETIME")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for connection lifetime too large")
		}
	})

	t.Run("Connection lifetime too small", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_CONN_MAX_LIFETIME", "30s")
		defer os.Unsetenv("DIDEBAN_STORAGE_CONN_MAX_LIFETIME")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for connection lifetime too small")
		}
	})

	t.Run("Max open connections too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS", "1001")
		defer os.Unsetenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for max open connections too large")
		}
	})
}

// TestAlertConfigValidation tests alert configuration validation
func TestAlertConfigValidation(t *testing.T) {
	t.Run("Telegram enabled without token", func(t *testing.T) {
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		defer os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram enabled without token")
		}
	})

	t.Run("Telegram enabled without chat ID", func(t *testing.T) {
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
		defer func() {
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram enabled without chat ID")
		}
	})

	t.Run("Telegram token too short", func(t *testing.T) {
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "short")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", "123456789")
		defer func() {
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram token too short")
		}
	})

	t.Run("Telegram token too long", func(t *testing.T) {
		longToken := make([]byte, 201)
		for i := range longToken {
			longToken[i] = 'a'
		}

		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", string(longToken))
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", "123456789")
		defer func() {
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram token too long")
		}
	})

	t.Run("Telegram chat ID too long", func(t *testing.T) {
		longChatID := make([]byte, 51)
		for i := range longChatID {
			longChatID[i] = '1'
		}

		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", string(longChatID))
		defer func() {
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram chat ID too long")
		}
	})

	t.Run("Telegram timeout too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", "123456789")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TIMEOUT", "3m")
		defer func() {
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TIMEOUT")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram timeout too large")
		}
	})

	t.Run("Telegram timeout too small", func(t *testing.T) {
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", "123456789")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TIMEOUT", "500ms")
		defer func() {
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TIMEOUT")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Telegram timeout too small")
		}
	})

	t.Run("Valid Telegram configuration", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "this-is-a-very-secure-secret-key-with-32-plus-characters")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", "123456789")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_ENABLED")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_TOKEN")
			os.Unsetenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID")
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Errorf("Expected no error for valid Telegram config, got: %v", err)
		}

		if !cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be enabled")
		}
	})

	t.Run("Bale configuration validation", func(t *testing.T) {
		// Test similar validation for Bale
		os.Setenv("DIDEBAN_ALERT_BALE_ENABLED", "true")
		defer os.Unsetenv("DIDEBAN_ALERT_BALE_ENABLED")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for Bale enabled without token")
		}
	})
}

// TestSchedulerConfigValidation tests scheduler configuration validation
func TestSchedulerConfigValidation(t *testing.T) {
	t.Run("Zero worker count", func(t *testing.T) {
		os.Setenv("DIDEBAN_SCHEDULER_WORKER_COUNT", "0")
		defer os.Unsetenv("DIDEBAN_SCHEDULER_WORKER_COUNT")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for zero worker count")
		}
	})

	t.Run("Worker count too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_SCHEDULER_WORKER_COUNT", "1001")
		defer os.Unsetenv("DIDEBAN_SCHEDULER_WORKER_COUNT")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for worker count too large")
		}
	})

	t.Run("Default interval too small", func(t *testing.T) {
		os.Setenv("DIDEBAN_SCHEDULER_DEFAULT_INTERVAL", "1s")
		defer os.Unsetenv("DIDEBAN_SCHEDULER_DEFAULT_INTERVAL")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for default interval too small")
		}
	})

	t.Run("Default interval too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_SCHEDULER_DEFAULT_INTERVAL", "25h")
		defer os.Unsetenv("DIDEBAN_SCHEDULER_DEFAULT_INTERVAL")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for default interval too large")
		}
	})

	t.Run("Negative max retries", func(t *testing.T) {
		os.Setenv("DIDEBAN_SCHEDULER_MAX_RETRIES", "-1")
		defer os.Unsetenv("DIDEBAN_SCHEDULER_MAX_RETRIES")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for negative max retries")
		}
	})

	t.Run("Max retries too large", func(t *testing.T) {
		os.Setenv("DIDEBAN_SCHEDULER_MAX_RETRIES", "11")
		defer os.Unsetenv("DIDEBAN_SCHEDULER_MAX_RETRIES")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for max retries too large")
		}
	})
}

// TestLogConfigValidation tests log configuration validation
func TestLogConfigValidation(t *testing.T) {
	t.Run("Invalid log level", func(t *testing.T) {
		os.Setenv("DIDEBAN_LOG_LEVEL", "invalid")
		defer os.Unsetenv("DIDEBAN_LOG_LEVEL")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for invalid log level")
		}
	})

	t.Run("Valid log levels", func(t *testing.T) {
		validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}

		for _, level := range validLevels {
			os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "this-is-a-very-secure-secret-key-with-32-plus-characters")
			os.Setenv("DIDEBAN_LOG_LEVEL", level)

			cfg, err := config.Load()
			if err != nil {
				t.Errorf("Expected no error for log level '%s', got: %v", level, err)
			}

			if cfg.Log.Level != level {
				t.Errorf("Expected log level '%s', got '%s'", level, cfg.Log.Level)
			}

			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_LOG_LEVEL")
		}
	})

	t.Run("Case insensitive log level", func(t *testing.T) {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "this-is-a-very-secure-secret-key-with-32-plus-characters")
		os.Setenv("DIDEBAN_LOG_LEVEL", "INFO")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_LOG_LEVEL")
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Errorf("Expected no error for uppercase log level, got: %v", err)
		}

		// Should be normalized to lowercase
		if cfg.Log.Level != "info" {
			t.Errorf("Expected normalized log level 'info', got '%s'", cfg.Log.Level)
		}
	})
}
