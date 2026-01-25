package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"dideban/internal/config"
)

// TestConfigLoadingFromEnvironment tests loading configuration from environment variables
func TestConfigLoadingFromEnvironment(t *testing.T) {
	// Clean up environment after test
	defer func() {
		envVars := []string{
			"DIDEBAN_SERVER_ADDR",
			"DIDEBAN_SERVER_READ_TIMEOUT",
			"DIDEBAN_SERVER_JWT_SECRET",
			"DIDEBAN_SERVER_JWT_TTL",
			"DIDEBAN_STORAGE_PATH",
			"DIDEBAN_STORAGE_MAX_OPEN_CONNS",
			"DIDEBAN_LOG_LEVEL",
			"DIDEBAN_LOG_PRETTY",
			"DIDEBAN_ALERT_TELEGRAM_ENABLED",
			"DIDEBAN_ALERT_TELEGRAM_TOKEN",
			"DIDEBAN_ALERT_TELEGRAM_CHAT_ID",
			"DIDEBAN_ALERT_TELEGRAM_TIMEOUT",
		}
		for _, env := range envVars {
			os.Unsetenv(env)
		}
	}()

	t.Run("Environment variables override defaults", func(t *testing.T) {
		// Set environment variables
		os.Setenv("DIDEBAN_SERVER_ADDR", ":9090")
		os.Setenv("DIDEBAN_SERVER_READ_TIMEOUT", "45s")
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")
		os.Setenv("DIDEBAN_SERVER_JWT_TTL", "12h")
		os.Setenv("DIDEBAN_STORAGE_PATH", "custom.db")
		os.Setenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS", "64")
		os.Setenv("DIDEBAN_LOG_LEVEL", "debug")
		os.Setenv("DIDEBAN_LOG_PRETTY", "true")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify environment variables were applied
		if cfg.Server.Addr != ":9090" {
			t.Errorf("Expected server addr ':9090', got '%s'", cfg.Server.Addr)
		}
		if cfg.Server.ReadTimeout != 45*time.Second {
			t.Errorf("Expected read timeout 45s, got %v", cfg.Server.ReadTimeout)
		}
		if cfg.Server.JWT.Secret != "custom-jwt-secret-key-with-32-plus-characters-for-security" {
			t.Errorf("Expected custom JWT secret, got '%s'", cfg.Server.JWT.Secret)
		}
		if cfg.Server.JWT.TTL != 12*time.Hour {
			t.Errorf("Expected JWT TTL 12h, got %v", cfg.Server.JWT.TTL)
		}
		if cfg.Storage.Path != "custom.db" {
			t.Errorf("Expected storage path 'custom.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Storage.MaxOpenConns != 64 {
			t.Errorf("Expected max open conns 64, got %d", cfg.Storage.MaxOpenConns)
		}
		if cfg.Log.Level != "debug" {
			t.Errorf("Expected log level 'debug', got '%s'", cfg.Log.Level)
		}
		if !cfg.Log.Pretty {
			t.Error("Expected log pretty to be true")
		}
	})

	t.Run("Nested environment variables", func(t *testing.T) {
		// Test nested configuration with underscores
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_ENABLED", "true")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TOKEN", "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_CHAT_ID", "123456789")
		os.Setenv("DIDEBAN_ALERT_TELEGRAM_TIMEOUT", "45s")
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be enabled")
		}
		if cfg.Alert.Telegram.Token != "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz" {
			t.Errorf("Expected Telegram token, got '%s'", cfg.Alert.Telegram.Token)
		}
		if cfg.Alert.Telegram.ChatID != "123456789" {
			t.Errorf("Expected Telegram chat ID '123456789', got '%s'", cfg.Alert.Telegram.ChatID)
		}
		if cfg.Alert.Telegram.Timeout != 45*time.Second {
			t.Errorf("Expected Telegram timeout 45s, got %v", cfg.Alert.Telegram.Timeout)
		}
	})

	t.Run("Boolean environment variables", func(t *testing.T) {
		testCases := []struct {
			value    string
			expected bool
		}{
			{"true", true},
			{"TRUE", true},
			{"True", true},
			{"1", true},
			{"false", false},
			{"FALSE", false},
			{"False", false},
			{"0", false},
			{"", false},
		}

		for _, tc := range testCases {
			os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")
			os.Setenv("DIDEBAN_LOG_PRETTY", tc.value)

			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("Failed to load config with pretty='%s': %v", tc.value, err)
			}

			if cfg.Log.Pretty != tc.expected {
				t.Errorf("Expected log pretty %v for value '%s', got %v", tc.expected, tc.value, cfg.Log.Pretty)
			}

			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_LOG_PRETTY")
		}
	})
}

// TestConfigLoadingFromFile tests loading configuration from YAML file
func TestConfigLoadingFromFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
server:
  addr: ":7070"
  read_timeout: "25s"
  write_timeout: "25s"
  idle_timeout: "50s"
  jwt:
    secret: "file-based-jwt-secret-key-with-32-plus-characters-for-security"
    ttl: "6h"

storage:
  path: "file-test.db"
  max_open_conns: 16
  max_idle_conns: 4
  conn_max_lifetime: "30m"

alert:
  telegram:
    enabled: true
    token: "file_telegram_token_1234567890"
    chat_id: "file_chat_id_123"
    timeout: "20s"
  bale:
    enabled: false

scheduler:
  worker_count: 4
  default_interval: "30s"
  max_retries: 5

checks:
  http:
    method: "POST"
    headers:
      authorization: "Bearer token"
      contenttype: "application/json"
    body: "test body"
    timeout_seconds: 15
    expected_status: 201
    expected_content: "success"
    follow_redirects: false
    verify_ssl: false
  ping:
    count: 5
    interval_ms: 500
    packet_size: 128
    timeout_seconds: 3

log:
  level: "warn"
  pretty: true
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Change to temp directory so config.yaml is found
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	// Verify file values were loaded
	if cfg.Server.Addr != ":7070" {
		t.Errorf("Expected server addr ':7070', got '%s'", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 25*time.Second {
		t.Errorf("Expected read timeout 25s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.JWT.Secret != "file-based-jwt-secret-key-with-32-plus-characters-for-security" {
		t.Errorf("Expected file JWT secret, got '%s'", cfg.Server.JWT.Secret)
	}
	if cfg.Server.JWT.TTL != 6*time.Hour {
		t.Errorf("Expected JWT TTL 6h, got %v", cfg.Server.JWT.TTL)
	}
	if cfg.Storage.Path != "file-test.db" {
		t.Errorf("Expected storage path 'file-test.db', got '%s'", cfg.Storage.Path)
	}
	if cfg.Storage.MaxOpenConns != 16 {
		t.Errorf("Expected max open conns 16, got %d", cfg.Storage.MaxOpenConns)
	}
	if cfg.Storage.MaxIdleConns != 4 {
		t.Errorf("Expected max idle conns 4, got %d", cfg.Storage.MaxIdleConns)
	}
	if cfg.Storage.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("Expected conn max lifetime 30m, got %v", cfg.Storage.ConnMaxLifetime)
	}
	if !cfg.Alert.Telegram.Enabled {
		t.Error("Expected Telegram to be enabled")
	}
	if cfg.Alert.Telegram.Token != "file_telegram_token_1234567890" {
		t.Errorf("Expected file Telegram token, got '%s'", cfg.Alert.Telegram.Token)
	}
	if cfg.Scheduler.WorkerCount != 4 {
		t.Errorf("Expected worker count 4, got %d", cfg.Scheduler.WorkerCount)
	}
	if cfg.Scheduler.DefaultInterval != 30*time.Second {
		t.Errorf("Expected default interval 30s, got %v", cfg.Scheduler.DefaultInterval)
	}
	if cfg.Checks.HTTP.Method != "POST" {
		t.Errorf("Expected HTTP method 'POST', got '%s'", cfg.Checks.HTTP.Method)
	}
	if cfg.Checks.HTTP.ExpectedStatus != 201 {
		t.Errorf("Expected HTTP status 201, got %d", cfg.Checks.HTTP.ExpectedStatus)
	}
	if cfg.Checks.HTTP.FollowRedirects {
		t.Error("Expected HTTP follow redirects to be false")
	}
	if cfg.Checks.Ping.Count != 5 {
		t.Errorf("Expected ping count 5, got %d", cfg.Checks.Ping.Count)
	}
	if cfg.Log.Level != "warn" {
		t.Errorf("Expected log level 'warn', got '%s'", cfg.Log.Level)
	}
	if !cfg.Log.Pretty {
		t.Error("Expected log pretty to be true")
	}

	// Test headers map
	if len(cfg.Checks.HTTP.Headers) != 2 {
		t.Errorf("Expected 2 HTTP headers, got %d", len(cfg.Checks.HTTP.Headers))
	}

	if cfg.Checks.HTTP.Headers["authorization"] != "Bearer token" {
		t.Errorf("Expected Authorization header, got '%s'", cfg.Checks.HTTP.Headers["Authorization"])
	}
	if cfg.Checks.HTTP.Headers["contenttype"] != "application/json" {
		t.Errorf("Expected Content-Type header, got '%s'", cfg.Checks.HTTP.Headers["Content-Type"])
	}
}

// TestConfigPrecedence tests that environment variables override file configuration
func TestConfigPrecedence(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
server:
  addr: ":7070"
  jwt:
    secret: "file-based-jwt-secret-key-with-32-plus-characters-for-security"
    ttl: "6h"
storage:
  path: "file-test.db"
log:
  level: "warn"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Set environment variables that should override file values
	os.Setenv("DIDEBAN_SERVER_ADDR", ":8888")
	os.Setenv("DIDEBAN_STORAGE_PATH", "env-test.db")
	os.Setenv("DIDEBAN_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("DIDEBAN_SERVER_ADDR")
		os.Unsetenv("DIDEBAN_STORAGE_PATH")
		os.Unsetenv("DIDEBAN_LOG_LEVEL")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Environment variables should override file values
	if cfg.Server.Addr != ":8888" {
		t.Errorf("Expected server addr ':8888' (from env), got '%s'", cfg.Server.Addr)
	}
	if cfg.Storage.Path != "env-test.db" {
		t.Errorf("Expected storage path 'env-test.db' (from env), got '%s'", cfg.Storage.Path)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Expected log level 'debug' (from env), got '%s'", cfg.Log.Level)
	}

	// File values should be used where no env override exists
	if cfg.Server.JWT.Secret != "file-based-jwt-secret-key-with-32-plus-characters-for-security" {
		t.Errorf("Expected JWT secret from file, got '%s'", cfg.Server.JWT.Secret)
	}
	if cfg.Server.JWT.TTL != 6*time.Hour {
		t.Errorf("Expected JWT TTL 6h (from file), got %v", cfg.Server.JWT.TTL)
	}
}

// TestConfigFileNotFound tests behavior when config file is not found
func TestConfigFileNotFound(t *testing.T) {
	// Change to a directory without config.yaml
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Set a secure JWT secret to avoid validation error
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing-only")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	// Should still work with defaults
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error when config file not found, got: %v", err)
	}

	// Should use default values
	if cfg.Server.Addr != ":8080" {
		t.Errorf("Expected default server addr ':8080', got '%s'", cfg.Server.Addr)
	}
	if cfg.Storage.Path != "dideban.db" {
		t.Errorf("Expected default storage path 'dideban.db', got '%s'", cfg.Storage.Path)
	}
}

// TestInvalidConfigFile tests behavior with invalid YAML file
func TestInvalidConfigFile(t *testing.T) {
	// Create a temporary invalid config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidContent := `
server:
  addr: ":8080"
  invalid_yaml: [
    missing_closing_bracket
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Should return error for invalid YAML
	_, err = config.Load()
	if err == nil {
		t.Error("Expected error for invalid YAML config file")
	}
}

// TestConfigNormalization tests that configuration values are properly normalized
func TestConfigNormalization(t *testing.T) {
	// Test log level normalization (should be lowercase)
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")
	os.Setenv("DIDEBAN_LOG_LEVEL", "INFO")
	defer func() {
		os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
		os.Unsetenv("DIDEBAN_LOG_LEVEL")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Log level should be normalized to lowercase
	if cfg.Log.Level != "info" {
		t.Errorf("Expected normalized log level 'info', got '%s'", cfg.Log.Level)
	}
}

// TestConfigDurationParsing tests that duration strings are properly parsed
func TestConfigDurationParsing(t *testing.T) {
	testCases := []struct {
		envVar   string
		value    string
		expected time.Duration
	}{
		{"DIDEBAN_SERVER_READ_TIMEOUT", "45s", 45 * time.Second},
		{"DIDEBAN_SERVER_JWT_TTL", "2h", 2 * time.Hour},
		{"DIDEBAN_STORAGE_CONN_MAX_LIFETIME", "90m", 90 * time.Minute},
		{"DIDEBAN_SCHEDULER_DEFAULT_INTERVAL", "120s", 120 * time.Second},
	}

	for _, tc := range testCases {
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")
		os.Setenv(tc.envVar, tc.value)

		cfg, err := config.Load()
		if err != nil {
			t.Errorf("Failed to load config with %s=%s: %v", tc.envVar, tc.value, err)
			continue
		}

		var actual time.Duration
		switch tc.envVar {
		case "DIDEBAN_SERVER_READ_TIMEOUT":
			actual = cfg.Server.ReadTimeout
		case "DIDEBAN_SERVER_JWT_TTL":
			actual = cfg.Server.JWT.TTL
		case "DIDEBAN_STORAGE_CONN_MAX_LIFETIME":
			actual = cfg.Storage.ConnMaxLifetime
		case "DIDEBAN_SCHEDULER_DEFAULT_INTERVAL":
			actual = cfg.Scheduler.DefaultInterval
		}

		if actual != tc.expected {
			t.Errorf("Expected %s to be %v, got %v", tc.envVar, tc.expected, actual)
		}

		os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
		os.Unsetenv(tc.envVar)
	}
}
