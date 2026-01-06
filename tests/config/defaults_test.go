package tests

import (
	"os"
	"testing"
	"time"

	"dideban/internal/config"
)

// TestConfigDefaults tests that default values are properly set
func TestConfigDefaults(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Run("Server defaults", func(t *testing.T) {
		if cfg.Server.Addr != ":8080" {
			t.Errorf("Expected server addr ':8080', got '%s'", cfg.Server.Addr)
		}
		if cfg.Server.ReadTimeout != 30*time.Second {
			t.Errorf("Expected read timeout 30s, got %v", cfg.Server.ReadTimeout)
		}
		if cfg.Server.WriteTimeout != 30*time.Second {
			t.Errorf("Expected write timeout 30s, got %v", cfg.Server.WriteTimeout)
		}
		if cfg.Server.IdleTimeout != 60*time.Second {
			t.Errorf("Expected idle timeout 60s, got %v", cfg.Server.IdleTimeout)
		}
	})

	t.Run("JWT defaults", func(t *testing.T) {
		if cfg.Server.JWT.Secret != "test-jwt-secret-key-with-32-plus-characters-for-testing" {
			t.Errorf("Expected test JWT secret, got '%s'", cfg.Server.JWT.Secret)
		}
		if cfg.Server.JWT.TTL != 24*time.Hour {
			t.Errorf("Expected JWT TTL 24h, got %v", cfg.Server.JWT.TTL)
		}
	})

	t.Run("Storage defaults", func(t *testing.T) {
		if cfg.Storage.Path != "dideban.db" {
			t.Errorf("Expected storage path 'dideban.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Storage.MaxOpenConns != 32 {
			t.Errorf("Expected max open conns 32, got %d", cfg.Storage.MaxOpenConns)
		}
		if cfg.Storage.MaxIdleConns != 8 {
			t.Errorf("Expected max idle conns 8, got %d", cfg.Storage.MaxIdleConns)
		}
		if cfg.Storage.ConnMaxLifetime != time.Hour {
			t.Errorf("Expected conn max lifetime 1h, got %v", cfg.Storage.ConnMaxLifetime)
		}
	})

	t.Run("Alert defaults", func(t *testing.T) {
		if cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be disabled by default")
		}
		if cfg.Alert.Telegram.Timeout != 30*time.Second {
			t.Errorf("Expected Telegram timeout 30s, got %v", cfg.Alert.Telegram.Timeout)
		}
		if cfg.Alert.Bale.Enabled {
			t.Error("Expected Bale to be disabled by default")
		}
		if cfg.Alert.Bale.Timeout != 30*time.Second {
			t.Errorf("Expected Bale timeout 30s, got %v", cfg.Alert.Bale.Timeout)
		}
	})

	t.Run("Scheduler defaults", func(t *testing.T) {
		if cfg.Scheduler.WorkerCount != 8 {
			t.Errorf("Expected worker count 8, got %d", cfg.Scheduler.WorkerCount)
		}
		if cfg.Scheduler.DefaultInterval != 60*time.Second {
			t.Errorf("Expected default interval 60s, got %v", cfg.Scheduler.DefaultInterval)
		}
		if cfg.Scheduler.MaxRetries != 3 {
			t.Errorf("Expected max retries 3, got %d", cfg.Scheduler.MaxRetries)
		}
	})

	t.Run("HTTP check defaults", func(t *testing.T) {
		if cfg.Checks.HTTP.Method != "GET" {
			t.Errorf("Expected HTTP method 'GET', got '%s'", cfg.Checks.HTTP.Method)
		}
		if cfg.Checks.HTTP.TimeoutSeconds != 10 {
			t.Errorf("Expected HTTP timeout 10s, got %d", cfg.Checks.HTTP.TimeoutSeconds)
		}
		if cfg.Checks.HTTP.ExpectedStatus != 200 {
			t.Errorf("Expected HTTP status 200, got %d", cfg.Checks.HTTP.ExpectedStatus)
		}
		if !cfg.Checks.HTTP.FollowRedirects {
			t.Error("Expected HTTP follow redirects to be true")
		}
		if !cfg.Checks.HTTP.VerifySSL {
			t.Error("Expected HTTP verify SSL to be true")
		}
		if cfg.Checks.HTTP.Headers == nil {
			t.Error("Expected HTTP headers to be initialized")
		}
		if cfg.Checks.HTTP.Body != "" {
			t.Errorf("Expected empty HTTP body, got '%s'", cfg.Checks.HTTP.Body)
		}
		if cfg.Checks.HTTP.ExpectedContent != "" {
			t.Errorf("Expected empty HTTP expected content, got '%s'", cfg.Checks.HTTP.ExpectedContent)
		}
	})

	t.Run("Ping check defaults", func(t *testing.T) {
		if cfg.Checks.Ping.Count != 3 {
			t.Errorf("Expected ping count 3, got %d", cfg.Checks.Ping.Count)
		}
		if cfg.Checks.Ping.IntervalMs != 1000 {
			t.Errorf("Expected ping interval 1000ms, got %d", cfg.Checks.Ping.IntervalMs)
		}
		if cfg.Checks.Ping.PacketSize != 56 {
			t.Errorf("Expected ping packet size 56, got %d", cfg.Checks.Ping.PacketSize)
		}
		if cfg.Checks.Ping.TimeoutSeconds != 5 {
			t.Errorf("Expected ping timeout 5s, got %d", cfg.Checks.Ping.TimeoutSeconds)
		}
	})

	t.Run("Log defaults", func(t *testing.T) {
		if cfg.Log.Level != "info" {
			t.Errorf("Expected log level 'info', got '%s'", cfg.Log.Level)
		}
		if cfg.Log.Pretty {
			t.Error("Expected log pretty to be false")
		}
	})
}

// TestDefaultsConsistency tests that defaults are consistent across different loading methods
func TestDefaultsConsistency(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	// Load config multiple times to ensure consistency
	cfg1, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config first time: %v", err)
	}

	cfg2, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config second time: %v", err)
	}

	// Compare key default values
	if cfg1.Server.Addr != cfg2.Server.Addr {
		t.Error("Server addr defaults are inconsistent")
	}
	if cfg1.Storage.Path != cfg2.Storage.Path {
		t.Error("Storage path defaults are inconsistent")
	}
	if cfg1.Scheduler.WorkerCount != cfg2.Scheduler.WorkerCount {
		t.Error("Scheduler worker count defaults are inconsistent")
	}
	if cfg1.Log.Level != cfg2.Log.Level {
		t.Error("Log level defaults are inconsistent")
	}
}

// TestDefaultsWithEmptyConfig tests that defaults work with empty configuration
func TestDefaultsWithEmptyConfig(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	// This test ensures that even with no config file, defaults are applied
	// The Load() function should handle missing config files gracefully
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config with defaults: %v", err)
	}

	// Verify essential defaults are set
	if cfg.Server.Addr == "" {
		t.Error("Server addr should have default value")
	}
	if cfg.Storage.Path == "" {
		t.Error("Storage path should have default value")
	}
	if cfg.Log.Level == "" {
		t.Error("Log level should have default value")
	}
}

// TestTimeoutDefaults tests that all timeout values have reasonable defaults
func TestTimeoutDefaults(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Server timeouts
	if cfg.Server.ReadTimeout <= 0 {
		t.Error("Server read timeout should be positive")
	}
	if cfg.Server.WriteTimeout <= 0 {
		t.Error("Server write timeout should be positive")
	}
	if cfg.Server.IdleTimeout <= 0 {
		t.Error("Server idle timeout should be positive")
	}

	// JWT timeout
	if cfg.Server.JWT.TTL <= 0 {
		t.Error("JWT TTL should be positive")
	}

	// Storage timeout
	if cfg.Storage.ConnMaxLifetime <= 0 {
		t.Error("Storage connection max lifetime should be positive")
	}

	// Alert timeouts
	if cfg.Alert.Telegram.Timeout <= 0 {
		t.Error("Telegram timeout should be positive")
	}
	if cfg.Alert.Bale.Timeout <= 0 {
		t.Error("Bale timeout should be positive")
	}

	// Scheduler timeout
	if cfg.Scheduler.DefaultInterval <= 0 {
		t.Error("Scheduler default interval should be positive")
	}
}

// TestCountDefaults tests that all count/numeric values have reasonable defaults
func TestCountDefaults(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Storage connection counts
	if cfg.Storage.MaxOpenConns <= 0 {
		t.Error("Max open connections should be positive")
	}
	if cfg.Storage.MaxIdleConns < 0 {
		t.Error("Max idle connections should be non-negative")
	}

	// Scheduler counts
	if cfg.Scheduler.WorkerCount <= 0 {
		t.Error("Worker count should be positive")
	}
	if cfg.Scheduler.MaxRetries < 0 {
		t.Error("Max retries should be non-negative")
	}

	// HTTP check defaults
	if cfg.Checks.HTTP.TimeoutSeconds <= 0 {
		t.Error("HTTP timeout should be positive")
	}
	if cfg.Checks.HTTP.ExpectedStatus <= 0 {
		t.Error("HTTP expected status should be positive")
	}

	// Ping check defaults
	if cfg.Checks.Ping.Count <= 0 {
		t.Error("Ping count should be positive")
	}
	if cfg.Checks.Ping.IntervalMs <= 0 {
		t.Error("Ping interval should be positive")
	}
	if cfg.Checks.Ping.PacketSize <= 0 {
		t.Error("Ping packet size should be positive")
	}
	if cfg.Checks.Ping.TimeoutSeconds <= 0 {
		t.Error("Ping timeout should be positive")
	}
}

// TestBooleanDefaults tests that boolean values have correct defaults
func TestBooleanDefaults(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Alert defaults (should be disabled)
	if cfg.Alert.Telegram.Enabled {
		t.Error("Telegram should be disabled by default")
	}
	if cfg.Alert.Bale.Enabled {
		t.Error("Bale should be disabled by default")
	}

	// HTTP check defaults (should be enabled)
	if !cfg.Checks.HTTP.FollowRedirects {
		t.Error("HTTP follow redirects should be enabled by default")
	}
	if !cfg.Checks.HTTP.VerifySSL {
		t.Error("HTTP verify SSL should be enabled by default")
	}

	// Log defaults
	if cfg.Log.Pretty {
		t.Error("Log pretty should be disabled by default")
	}
}
