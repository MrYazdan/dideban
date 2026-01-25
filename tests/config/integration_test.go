package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dideban/internal/config"
)

// cleanupDidebanEnvVars removes all DIDEBAN_* environment variables
func cleanupDidebanEnvVars() {
	for _, env := range os.Environ() {
		if len(env) >= 8 && env[:8] == "DIDEBAN_" {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				os.Unsetenv(parts[0])
			}
		}
	}
}

// TestConfigIntegration tests complete configuration scenarios
func TestConfigIntegration(t *testing.T) {
	t.Run("Production-like configuration", func(t *testing.T) {
		// Clean up all DIDEBAN environment variables first
		cleanupDidebanEnvVars()

		// Set up a production-like configuration
		envVars := map[string]string{
			"DIDEBAN_SERVER_ADDR":                  ":443",
			"DIDEBAN_SERVER_READ_TIMEOUT":          "60s",
			"DIDEBAN_SERVER_WRITE_TIMEOUT":         "60s",
			"DIDEBAN_SERVER_IDLE_TIMEOUT":          "120s",
			"DIDEBAN_SERVER_JWT_SECRET":            "production-super-secure-jwt-secret-key-with-64-characters-minimum",
			"DIDEBAN_SERVER_JWT_TTL":               "8h",
			"DIDEBAN_STORAGE_PATH":                 "/var/lib/dideban/production.db",
			"DIDEBAN_STORAGE_MAX_OPEN_CONNS":       "100",
			"DIDEBAN_STORAGE_MAX_IDLE_CONNS":       "25",
			"DIDEBAN_STORAGE_CONN_MAX_LIFETIME":    "4h",
			"DIDEBAN_ALERT_TELEGRAM_ENABLED":       "false", // Disable telegram to avoid token requirement
			"DIDEBAN_ALERT_BALE_ENABLED":           "false", // Disable bale to avoid token requirement
			"DIDEBAN_SCHEDULER_WORKER_COUNT":       "16",
			"DIDEBAN_SCHEDULER_DEFAULT_INTERVAL":   "30s",
			"DIDEBAN_SCHEDULER_MAX_RETRIES":        "5",
			"DIDEBAN_CHECKS_HTTP_TIMEOUT_SECONDS":  "30",
			"DIDEBAN_CHECKS_HTTP_EXPECTED_STATUS":  "200",
			"DIDEBAN_CHECKS_HTTP_FOLLOW_REDIRECTS": "true",
			"DIDEBAN_CHECKS_HTTP_VERIFY_SSL":       "true",
			"DIDEBAN_CHECKS_PING_COUNT":            "5",
			"DIDEBAN_CHECKS_PING_INTERVAL_MS":      "2000",
			"DIDEBAN_CHECKS_PING_PACKET_SIZE":      "64",
			"DIDEBAN_CHECKS_PING_TIMEOUT_SECONDS":  "10",
			"DIDEBAN_LOG_LEVEL":                    "warn",
			"DIDEBAN_LOG_PRETTY":                   "false",
		}

		// Set environment variables
		for key, value := range envVars {
			os.Setenv(key, value)
		}

		// Clean up after test
		defer cleanupDidebanEnvVars()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load production config: %v", err)
		}

		// Verify production configuration
		if cfg.Server.Addr != ":443" {
			t.Errorf("Expected server addr ':443', got '%s'", cfg.Server.Addr)
		}
		if cfg.Server.ReadTimeout != 60*time.Second {
			t.Errorf("Expected read timeout 60s, got %v", cfg.Server.ReadTimeout)
		}
		if cfg.Server.JWT.TTL != 8*time.Hour {
			t.Errorf("Expected JWT TTL 8h, got %v", cfg.Server.JWT.TTL)
		}
		if cfg.Storage.MaxOpenConns != 100 {
			t.Errorf("Expected max open conns 100, got %d", cfg.Storage.MaxOpenConns)
		}
		if cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be disabled")
		}
		if cfg.Alert.Bale.Enabled {
			t.Error("Expected Bale to be disabled")
		}
		if cfg.Scheduler.WorkerCount != 16 {
			t.Errorf("Expected worker count 16, got %d", cfg.Scheduler.WorkerCount)
		}
		if cfg.Log.Level != "warn" {
			t.Errorf("Expected log level 'warn', got '%s'", cfg.Log.Level)
		}
	})

	t.Run("Development configuration", func(t *testing.T) {
		// Create a development config file
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")

		devConfig := `
server:
  addr: ":3000"
  read_timeout: "10s"
  write_timeout: "10s"
  idle_timeout: "30s"
  jwt:
    secret: "development-jwt-secret-key-not-for-production-use-only"
    ttl: "1h"

storage:
  path: "dev.db"
  max_open_conns: 5
  max_idle_conns: 2
  conn_max_lifetime: "10m"

alert:
  telegram:
    enabled: false
  bale:
    enabled: false

scheduler:
  worker_count: 2
  default_interval: "10s"
  max_retries: 1

checks:
  http:
    method: "GET"
    timeout_seconds: 5
    expected_status: 200
    follow_redirects: true
    verify_ssl: false  # Allow self-signed certs in dev
  ping:
    count: 1
    interval_ms: 500
    packet_size: 32
    timeout_seconds: 2

log:
  level: "debug"
  pretty: true
`

		err := os.WriteFile(configPath, []byte(devConfig), 0644)
		if err != nil {
			t.Fatalf("Failed to write dev config: %v", err)
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

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load dev config: %v", err)
		}

		// Verify development configuration
		if cfg.Server.Addr != ":3000" {
			t.Errorf("Expected server addr ':3000', got '%s'", cfg.Server.Addr)
		}
		if cfg.Server.ReadTimeout != 10*time.Second {
			t.Errorf("Expected read timeout 10s, got %v", cfg.Server.ReadTimeout)
		}
		if cfg.Storage.Path != "dev.db" {
			t.Errorf("Expected storage path 'dev.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be disabled in dev")
		}
		if cfg.Alert.Bale.Enabled {
			t.Error("Expected Bale to be disabled in dev")
		}
		if cfg.Scheduler.WorkerCount != 2 {
			t.Errorf("Expected worker count 2, got %d", cfg.Scheduler.WorkerCount)
		}
		if cfg.Checks.HTTP.VerifySSL {
			t.Error("Expected SSL verification to be disabled in dev")
		}
		if cfg.Log.Level != "debug" {
			t.Errorf("Expected log level 'debug', got '%s'", cfg.Log.Level)
		}
		if !cfg.Log.Pretty {
			t.Error("Expected pretty logging in dev")
		}
	})

	t.Run("Minimal configuration", func(t *testing.T) {
		// Test with only required overrides
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "minimal-but-secure-jwt-secret-key-with-32-characters")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load minimal config: %v", err)
		}

		// Should use defaults for everything else
		if cfg.Server.Addr != ":8080" {
			t.Errorf("Expected default server addr ':8080', got '%s'", cfg.Server.Addr)
		}
		if cfg.Storage.Path != "dideban.db" {
			t.Errorf("Expected default storage path 'dideban.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be disabled by default")
		}
		if cfg.Log.Level != "info" {
			t.Errorf("Expected default log level 'info', got '%s'", cfg.Log.Level)
		}
	})

	t.Run("High-performance configuration", func(t *testing.T) {
		// Configuration optimized for high performance
		envVars := map[string]string{
			"DIDEBAN_SERVER_JWT_SECRET":           "high-performance-jwt-secret-key-with-32-plus-characters",
			"DIDEBAN_STORAGE_MAX_OPEN_CONNS":      "200",
			"DIDEBAN_STORAGE_MAX_IDLE_CONNS":      "50",
			"DIDEBAN_STORAGE_CONN_MAX_LIFETIME":   "30m",
			"DIDEBAN_SCHEDULER_WORKER_COUNT":      "32",
			"DIDEBAN_SCHEDULER_DEFAULT_INTERVAL":  "15s",
			"DIDEBAN_CHECKS_HTTP_TIMEOUT_SECONDS": "5",
			"DIDEBAN_CHECKS_PING_COUNT":           "1",
			"DIDEBAN_CHECKS_PING_INTERVAL_MS":     "100",
			"DIDEBAN_CHECKS_PING_TIMEOUT_SECONDS": "2",
			"DIDEBAN_LOG_LEVEL":                   "error", // Minimal logging
		}

		for key, value := range envVars {
			os.Setenv(key, value)
		}

		defer func() {
			for key := range envVars {
				os.Unsetenv(key)
			}
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load high-performance config: %v", err)
		}

		// Verify high-performance settings
		if cfg.Storage.MaxOpenConns != 200 {
			t.Errorf("Expected max open conns 200, got %d", cfg.Storage.MaxOpenConns)
		}
		if cfg.Scheduler.WorkerCount != 32 {
			t.Errorf("Expected worker count 32, got %d", cfg.Scheduler.WorkerCount)
		}
		if cfg.Scheduler.DefaultInterval != 15*time.Second {
			t.Errorf("Expected default interval 15s, got %v", cfg.Scheduler.DefaultInterval)
		}
		if cfg.Checks.HTTP.TimeoutSeconds != 5 {
			t.Errorf("Expected HTTP timeout 5s, got %d", cfg.Checks.HTTP.TimeoutSeconds)
		}
		if cfg.Checks.Ping.Count != 1 {
			t.Errorf("Expected ping count 1, got %d", cfg.Checks.Ping.Count)
		}
		if cfg.Log.Level != "error" {
			t.Errorf("Expected log level 'error', got '%s'", cfg.Log.Level)
		}
	})

	t.Run("Security-focused configuration", func(t *testing.T) {
		// Configuration with security best practices
		envVars := map[string]string{
			"DIDEBAN_SERVER_ADDR":                  "127.0.0.1:8080", // Bind to localhost only
			"DIDEBAN_SERVER_READ_TIMEOUT":          "30s",
			"DIDEBAN_SERVER_WRITE_TIMEOUT":         "30s",
			"DIDEBAN_SERVER_IDLE_TIMEOUT":          "60s",
			"DIDEBAN_SERVER_JWT_SECRET":            "ultra-secure-jwt-secret-key-with-64-characters-for-maximum-security",
			"DIDEBAN_SERVER_JWT_TTL":               "2h",                  // Shorter TTL for security
			"DIDEBAN_STORAGE_PATH":                 "./secure/dideban.db", // Specific directory
			"DIDEBAN_CHECKS_HTTP_VERIFY_SSL":       "true",
			"DIDEBAN_CHECKS_HTTP_FOLLOW_REDIRECTS": "false", // Don't follow redirects for security
			"DIDEBAN_LOG_LEVEL":                    "info",
		}

		for key, value := range envVars {
			os.Setenv(key, value)
		}

		defer func() {
			for key := range envVars {
				os.Unsetenv(key)
			}
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load security-focused config: %v", err)
		}

		// Verify security settings
		if cfg.Server.Addr != "127.0.0.1:8080" {
			t.Errorf("Expected server addr '127.0.0.1:8080', got '%s'", cfg.Server.Addr)
		}
		if cfg.Server.JWT.TTL != 2*time.Hour {
			t.Errorf("Expected JWT TTL 2h, got %v", cfg.Server.JWT.TTL)
		}
		if cfg.Storage.Path != "./secure/dideban.db" {
			t.Errorf("Expected storage path './secure/dideban.db', got '%s'", cfg.Storage.Path)
		}
		if !cfg.Checks.HTTP.VerifySSL {
			t.Error("Expected SSL verification to be enabled")
		}
		if cfg.Checks.HTTP.FollowRedirects {
			t.Error("Expected redirect following to be disabled for security")
		}
	})
}

// TestConfigCompatibility tests backward compatibility scenarios
func TestConfigCompatibility(t *testing.T) {
	t.Run("Legacy environment variable names", func(t *testing.T) {
		// Test that the current environment variable naming works
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "legacy-compatible-jwt-secret-key-with-32-characters")
		os.Setenv("DIDEBAN_SERVER_ADDR", ":8081")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_SERVER_ADDR")
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config with current env vars: %v", err)
		}

		if cfg.Server.Addr != ":8081" {
			t.Errorf("Expected server addr ':8081', got '%s'", cfg.Server.Addr)
		}
	})

	t.Run("Mixed case environment variables", func(t *testing.T) {
		// Environment variables should be case-insensitive due to OS differences
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "mixed-case-jwt-secret-key-with-32-characters")
		os.Setenv("dideban_server_addr", ":8082") // lowercase
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("dideban_server_addr")
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// The lowercase env var might not work depending on the system
		// This test documents the expected behavior
		if cfg.Server.Addr == ":8082" {
			t.Log("Lowercase environment variables are supported")
		} else {
			t.Log("Lowercase environment variables are not supported (expected)")
		}
	})
}

// TestConfigPerformance tests configuration loading performance
func TestConfigPerformance(t *testing.T) {
	t.Run("Config loading performance", func(t *testing.T) {
		// Set up a complex configuration
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "performance-test-jwt-secret-key-with-32-characters")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		// Measure time to load config multiple times
		const iterations = 100
		start := time.Now()

		for i := 0; i < iterations; i++ {
			_, err := config.Load()
			if err != nil {
				t.Fatalf("Failed to load config iteration %d: %v", i, err)
			}
		}

		elapsed := time.Since(start)
		avgTime := elapsed / iterations

		t.Logf("Average config load time: %v", avgTime)

		// Config loading should be reasonably fast (less than 10ms per load)
		if avgTime > 10*time.Millisecond {
			t.Errorf("Config loading too slow: %v per load", avgTime)
		}
	})
}

// TestConfigRealWorldScenarios tests real-world usage scenarios
func TestConfigRealWorldScenarios(t *testing.T) {
	t.Run("Docker container scenario", func(t *testing.T) {
		// Simulate Docker container environment
		envVars := map[string]string{
			"DIDEBAN_SERVER_ADDR":            "0.0.0.0:8080", // Bind to all interfaces
			"DIDEBAN_SERVER_JWT_SECRET":      "docker-container-jwt-secret-key-with-32-characters",
			"DIDEBAN_STORAGE_PATH":           "/data/dideban.db",
			"DIDEBAN_LOG_LEVEL":              "info",
			"DIDEBAN_LOG_PRETTY":             "false", // JSON logging for containers
			"DIDEBAN_ALERT_TELEGRAM_ENABLED": "true",
			"DIDEBAN_ALERT_TELEGRAM_TOKEN":   "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
			"DIDEBAN_ALERT_TELEGRAM_CHAT_ID": "-1001234567890",
			"DIDEBAN_ALERT_BALE_ENABLED":     "false", // Disable bale to avoid token requirement
		}

		// Clean up all DIDEBAN environment variables first
		cleanupDidebanEnvVars()

		for key, value := range envVars {
			os.Setenv(key, value)
		}

		defer cleanupDidebanEnvVars()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load Docker config: %v", err)
		}

		// Verify Docker-appropriate settings
		if cfg.Server.Addr != "0.0.0.0:8080" {
			t.Errorf("Expected server addr '0.0.0.0:8080', got '%s'", cfg.Server.Addr)
		}
		if cfg.Storage.Path != "/data/dideban.db" {
			t.Errorf("Expected storage path '/data/dideban.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Log.Pretty {
			t.Error("Expected JSON logging for containers")
		}
		if !cfg.Alert.Telegram.Enabled {
			t.Error("Expected Telegram to be enabled")
		}
	})

	t.Run("Kubernetes deployment scenario", func(t *testing.T) {
		// Simulate Kubernetes environment with secrets and config maps
		envVars := map[string]string{
			"DIDEBAN_SERVER_ADDR":            "0.0.0.0:8080",
			"DIDEBAN_SERVER_JWT_SECRET":      "k8s-secret-jwt-key-from-kubernetes-secret-with-32-chars",
			"DIDEBAN_STORAGE_PATH":           "/var/lib/dideban/data.db",
			"DIDEBAN_ALERT_TELEGRAM_TOKEN":   "k8s_telegram_token_from_secret",
			"DIDEBAN_ALERT_TELEGRAM_CHAT_ID": "k8s_chat_id_from_configmap",
			"DIDEBAN_ALERT_TELEGRAM_ENABLED": "true",
			"DIDEBAN_ALERT_BALE_ENABLED":     "false", // Disable bale to avoid token requirement
			"DIDEBAN_SCHEDULER_WORKER_COUNT": "8",
			"DIDEBAN_LOG_LEVEL":              "info",
			"DIDEBAN_LOG_PRETTY":             "false",
		}

		// Clean up all DIDEBAN environment variables first
		cleanupDidebanEnvVars()

		for key, value := range envVars {
			os.Setenv(key, value)
		}

		defer cleanupDidebanEnvVars()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load Kubernetes config: %v", err)
		}

		// Verify Kubernetes-appropriate settings
		if cfg.Server.Addr != "0.0.0.0:8080" {
			t.Errorf("Expected server addr '0.0.0.0:8080', got '%s'", cfg.Server.Addr)
		}
		if cfg.Storage.Path != "/var/lib/dideban/data.db" {
			t.Errorf("Expected storage path '/var/lib/dideban/data.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Alert.Telegram.Token != "k8s_telegram_token_from_secret" {
			t.Errorf("Expected Telegram token from secret, got '%s'", cfg.Alert.Telegram.Token)
		}
		if cfg.Scheduler.WorkerCount != 8 {
			t.Errorf("Expected worker count 8, got %d", cfg.Scheduler.WorkerCount)
		}
	})

	t.Run("Systemd service scenario", func(t *testing.T) {
		// Simulate systemd service environment
		envVars := map[string]string{
			"DIDEBAN_SERVER_ADDR":            "127.0.0.1:8080", // Local binding for systemd
			"DIDEBAN_SERVER_JWT_SECRET":      "systemd-service-jwt-secret-key-with-32-characters",
			"DIDEBAN_STORAGE_PATH":           "/var/lib/dideban/dideban.db",
			"DIDEBAN_LOG_LEVEL":              "info",
			"DIDEBAN_LOG_PRETTY":             "false", // Structured logging for journald
			"DIDEBAN_SCHEDULER_WORKER_COUNT": "4",
		}

		for key, value := range envVars {
			os.Setenv(key, value)
		}

		defer func() {
			for key := range envVars {
				os.Unsetenv(key)
			}
		}()

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load systemd config: %v", err)
		}

		// Verify systemd-appropriate settings
		if cfg.Server.Addr != "127.0.0.1:8080" {
			t.Errorf("Expected server addr '127.0.0.1:8080', got '%s'", cfg.Server.Addr)
		}
		if cfg.Storage.Path != "/var/lib/dideban/dideban.db" {
			t.Errorf("Expected storage path '/var/lib/dideban/dideban.db', got '%s'", cfg.Storage.Path)
		}
		if cfg.Log.Pretty {
			t.Error("Expected structured logging for systemd")
		}
	})
}
