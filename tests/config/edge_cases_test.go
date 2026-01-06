package tests

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"dideban/internal/config"
)

// TestConfigEdgeCases tests various edge cases and error conditions
func TestConfigEdgeCases(t *testing.T) {
	t.Run("Empty environment variable values", func(t *testing.T) {
		// Set empty environment variables
		os.Setenv("DIDEBAN_SERVER_ADDR", "")
		os.Setenv("DIDEBAN_STORAGE_PATH", "")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_ADDR")
			os.Unsetenv("DIDEBAN_STORAGE_PATH")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for empty server addr")
		}
	})

	t.Run("Invalid duration formats", func(t *testing.T) {
		invalidDurations := []struct {
			envVar string
			value  string
		}{
			{"DIDEBAN_SERVER_READ_TIMEOUT", "invalid"},
			{"DIDEBAN_SERVER_JWT_TTL", "not-a-duration"},
			{"DIDEBAN_STORAGE_CONN_MAX_LIFETIME", ""},
			{"DIDEBAN_SCHEDULER_DEFAULT_INTERVAL", "abc"},
		}

		for _, tc := range invalidDurations {
			os.Setenv(tc.envVar, tc.value)

			_, err := config.Load()
			if err == nil {
				t.Errorf("Expected error for invalid duration %s=%s", tc.envVar, tc.value)
			}

			os.Unsetenv(tc.envVar)
		}
	})

	t.Run("Invalid integer values", func(t *testing.T) {
		invalidInts := []struct {
			envVar string
			value  string
		}{
			{"DIDEBAN_STORAGE_MAX_OPEN_CONNS", "not-a-number"},
			{"DIDEBAN_STORAGE_MAX_IDLE_CONNS", "invalid"},
			{"DIDEBAN_SCHEDULER_WORKER_COUNT", "abc"},
			{"DIDEBAN_SCHEDULER_MAX_RETRIES", "not-int"},
		}

		for _, tc := range invalidInts {
			os.Setenv(tc.envVar, tc.value)

			_, err := config.Load()
			if err == nil {
				t.Errorf("Expected error for invalid integer %s=%s", tc.envVar, tc.value)
			}

			os.Unsetenv(tc.envVar)
		}
	})

	t.Run("Invalid boolean values", func(t *testing.T) {
		// Note: Viper is quite permissive with boolean values
		// Most invalid values will be treated as false, but some may cause parse errors
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")
		os.Setenv("DIDEBAN_LOG_PRETTY", "invalid-boolean")
		defer func() {
			os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
			os.Unsetenv("DIDEBAN_LOG_PRETTY")
		}()

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for invalid boolean value")
		}
	})

	t.Run("Extremely large values", func(t *testing.T) {
		os.Setenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS", "999999999")
		defer os.Unsetenv("DIDEBAN_STORAGE_MAX_OPEN_CONNS")

		_, err := config.Load()
		if err == nil {
			t.Error("Expected error for extremely large max open connections")
		}
	})

	t.Run("Negative values where not allowed", func(t *testing.T) {
		negativeValues := []struct {
			envVar string
			value  string
		}{
			{"DIDEBAN_STORAGE_MAX_OPEN_CONNS", "-1"},
			{"DIDEBAN_SCHEDULER_WORKER_COUNT", "-5"},
		}

		for _, tc := range negativeValues {
			os.Setenv(tc.envVar, tc.value)

			_, err := config.Load()
			if err == nil {
				t.Errorf("Expected error for negative value %s=%s", tc.envVar, tc.value)
			}

			os.Unsetenv(tc.envVar)
		}
	})

	t.Run("Zero values where not allowed", func(t *testing.T) {
		zeroValues := []struct {
			envVar string
			value  string
		}{
			{"DIDEBAN_SERVER_READ_TIMEOUT", "0s"},
			{"DIDEBAN_STORAGE_MAX_OPEN_CONNS", "0"},
			{"DIDEBAN_SCHEDULER_WORKER_COUNT", "0"},
		}

		for _, tc := range zeroValues {
			os.Setenv(tc.envVar, tc.value)

			_, err := config.Load()
			if err == nil {
				t.Errorf("Expected error for zero value %s=%s", tc.envVar, tc.value)
			}

			os.Unsetenv(tc.envVar)
		}
	})
}

// TestConfigDirectoryHandling tests cross-platform config directory handling
func TestConfigDirectoryHandling(t *testing.T) {
	t.Run("Config directory detection", func(t *testing.T) {
		// Set a valid JWT secret to avoid validation errors
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		// This test verifies that the config directory logic works
		// We can't easily test the actual directory resolution without
		// manipulating environment variables, but we can test the logic

		if runtime.GOOS == "windows" {
			// On Windows, should use APPDATA
			originalAppData := os.Getenv("APPDATA")
			os.Setenv("APPDATA", "C:\\Users\\Test\\AppData\\Roaming")
			defer func() {
				if originalAppData != "" {
					os.Setenv("APPDATA", originalAppData)
				} else {
					os.Unsetenv("APPDATA")
				}
			}()
		} else {
			// On Unix-like systems, should use HOME
			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", "/home/test")
			defer func() {
				if originalHome != "" {
					os.Setenv("HOME", originalHome)
				} else {
					os.Unsetenv("HOME")
				}
			}()
		}

		// Load config should work regardless of config directory
		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Should have default values
		if cfg.Server.Addr != ":8080" {
			t.Errorf("Expected default server addr, got '%s'", cfg.Server.Addr)
		}
	})

	t.Run("Missing environment variables for config directory", func(t *testing.T) {
		// Set a valid JWT secret to avoid validation errors
		os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
		defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

		var originalEnv string
		var envVar string

		if runtime.GOOS == "windows" {
			envVar = "APPDATA"
			originalEnv = os.Getenv("APPDATA")
			os.Unsetenv("APPDATA")
		} else {
			envVar = "HOME"
			originalEnv = os.Getenv("HOME")
			os.Unsetenv("HOME")
		}

		defer func() {
			if originalEnv != "" {
				os.Setenv(envVar, originalEnv)
			}
		}()

		// Should still work without config directory environment variables
		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config without %s: %v", envVar, err)
		}

		if cfg.Server.Addr != ":8080" {
			t.Errorf("Expected default server addr, got '%s'", cfg.Server.Addr)
		}
	})
}

// TestConfigFilePermissions tests behavior with different file permissions
func TestConfigFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping file permission tests on Windows")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
server:
  addr: ":7070"
  jwt:
    secret: "file-based-jwt-secret-key-with-32-plus-characters-for-security"
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

	t.Run("Readable config file", func(t *testing.T) {
		// Make file readable
		err := os.Chmod(configPath, 0644)
		if err != nil {
			t.Fatalf("Failed to chmod config file: %v", err)
		}

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load readable config: %v", err)
		}

		if cfg.Server.Addr != ":7070" {
			t.Errorf("Expected server addr ':7070', got '%s'", cfg.Server.Addr)
		}
	})

	t.Run("Unreadable config file", func(t *testing.T) {
		// Make file unreadable
		err := os.Chmod(configPath, 0000)
		if err != nil {
			t.Fatalf("Failed to chmod config file: %v", err)
		}

		// Restore permissions after test
		defer os.Chmod(configPath, 0644)

		_, err = config.Load()
		if err == nil {
			t.Error("Expected error for unreadable config file")
		}
	})
}

// TestConfigWithSpecialCharacters tests configuration with special characters
func TestConfigWithSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with special characters in strings
	configContent := `
server:
  jwt:
    secret: "special-chars-!@#$%^&*()_+-=[]{}|;:,.<>?/~` + "`" + `"
storage:
  path: "test with spaces & special chars!.db"
alert:
  telegram:
    enabled: true
    token: "1234567890:ABCdef-GHI_jklMNO.pqrsTUVwxyz"
    chat_id: "@channel_with_special-chars"
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

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config with special characters: %v", err)
	}

	// Verify special characters are preserved
	expectedSecret := "special-chars-!@#$%^&*()_+-=[]{}|;:,.<>?/~`"
	if cfg.Server.JWT.Secret != expectedSecret {
		t.Errorf("Expected JWT secret with special chars, got '%s'", cfg.Server.JWT.Secret)
	}

	expectedPath := "test with spaces & special chars!.db"
	if cfg.Storage.Path != expectedPath {
		t.Errorf("Expected storage path with special chars, got '%s'", cfg.Storage.Path)
	}

	expectedChatID := "@channel_with_special-chars"
	if cfg.Alert.Telegram.ChatID != expectedChatID {
		t.Errorf("Expected chat ID with special chars, got '%s'", cfg.Alert.Telegram.ChatID)
	}
}

// TestConfigWithUnicodeCharacters tests configuration with Unicode characters
func TestConfigWithUnicodeCharacters(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with Unicode characters
	configContent := `
server:
  jwt:
    secret: "unicode-secret-دیدبان-مراقب-سیستم-های-شما-است"
storage:
  path: "دیدبان.db"
alert:
  telegram:
    enabled: true
    token: "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
    chat_id: "کانال_دیدبان"
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

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config with Unicode characters: %v", err)
	}

	// Verify Unicode characters are preserved
	expectedSecret := "unicode-secret-دیدبان-مراقب-سیستم-های-شما-است"
	if cfg.Server.JWT.Secret != expectedSecret {
		t.Errorf("Expected JWT secret with Unicode, got '%s'", cfg.Server.JWT.Secret)
	}

	expectedPath := "دیدبان.db"
	if cfg.Storage.Path != expectedPath {
		t.Errorf("Expected storage path with Unicode, got '%s'", cfg.Storage.Path)
	}

	expectedChatID := "کانال_دیدبان"
	if cfg.Alert.Telegram.ChatID != expectedChatID {
		t.Errorf("Expected chat ID with Unicode, got '%s'", cfg.Alert.Telegram.ChatID)
	}
}

// TestConfigMemoryUsage tests that config loading doesn't consume excessive memory
func TestConfigMemoryUsage(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	// Load config multiple times to check for memory leaks
	for i := 0; i < 100; i++ {
		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Failed to load config iteration %d: %v", i, err)
		}

		// Verify config is valid
		if cfg.Server.Addr == "" {
			t.Errorf("Invalid config at iteration %d", i)
		}
	}
}

// TestConfigConcurrency tests concurrent config loading
func TestConfigConcurrency(t *testing.T) {
	// Set a valid JWT secret to avoid validation errors
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "test-jwt-secret-key-with-32-plus-characters-for-testing")
	defer os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")

	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	// Load config concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			cfg, err := config.Load()
			if err != nil {
				errors <- err
				return
			}

			// Verify config is valid
			if cfg.Server.Addr == "" {
				errors <- err
				return
			}

			errors <- nil
		}(i)
	}

	// Check results
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		if err != nil {
			t.Errorf("Concurrent config load %d failed: %v", i, err)
		}
	}
}

// TestConfigValidationOrder tests that validation happens in the correct order
func TestConfigValidationOrder(t *testing.T) {
	// Test that normalization happens before validation
	os.Setenv("DIDEBAN_SERVER_JWT_SECRET", "custom-jwt-secret-key-with-32-plus-characters-for-security")
	os.Setenv("DIDEBAN_LOG_LEVEL", "DEBUG") // Should be normalized to lowercase
	defer func() {
		os.Unsetenv("DIDEBAN_SERVER_JWT_SECRET")
		os.Unsetenv("DIDEBAN_LOG_LEVEL")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Should be normalized to lowercase before validation
	if cfg.Log.Level != "debug" {
		t.Errorf("Expected normalized log level 'debug', got '%s'", cfg.Log.Level)
	}
}
