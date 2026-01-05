package tests

import (
	"context"
	"testing"

	"dideban/internal/checks"
	"dideban/internal/config"
	"dideban/internal/storage"
)

func TestHTTPCheckerConfigIntegration(t *testing.T) {
	// Create test application config
	cfg := &config.Config{
		Checks: config.ChecksConfig{
			HTTP: config.HTTPDefaultsConfig{
				Method:          "GET",
				Headers:         map[string]string{"User-Agent": "Dideban", "Accept": "application/json"},
				Body:            "",
				TimeoutSeconds:  10,
				ExpectedStatus:  200,
				ExpectedContent: "",
				FollowRedirects: true,
				VerifySSL:       true,
			},
		},
	}

	checker := checks.NewHTTPChecker(cfg)

	t.Run("HTTP checker can be created with config", func(t *testing.T) {
		if checker == nil {
			t.Error("Expected checker to be created")
		}

		if checker.Type() != "http" {
			t.Errorf("Expected type 'http', got: %s", checker.Type())
		}
	})

	// Test actual check execution with different configs
	t.Run("HTTP check with empty config", func(t *testing.T) {
		check := &storage.Check{
			ID:              1,
			Name:            "Test HTTP Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://mockhttp.org/get",
			Config:          "",
			IntervalSeconds: 10,
			TimeoutSeconds:  5,
			Enabled:         true,
		}

		// This tests the full flow including config parsing
		result, err := checker.Check(context.Background(), check)
		if err != nil {
			t.Logf("Check failed (expected for network test): %v", err)
			// Network tests might fail in CI, so we just log
		} else {
			if result == nil {
				t.Error("Expected result to be returned")
			}
			if result.CheckID != check.ID {
				t.Errorf("Expected CheckID %d, got %d", check.ID, result.CheckID)
			}
		}
	})

	t.Run("HTTP check with invalid config should fail", func(t *testing.T) {
		check := &storage.Check{
			ID:              3,
			Name:            "Test Invalid HTTP Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://mockhttp.org/get",
			Config:          `{"method": "INVALID", "expected_status": 999}`,
			IntervalSeconds: 10,
			TimeoutSeconds:  5,
			Enabled:         true,
		}

		_, err := checker.Check(context.Background(), check)
		if err == nil {
			t.Error("Expected error for invalid config")
		}
	})
}

func TestPingCheckerConfigIntegration(t *testing.T) {
	// Create test application config
	cfg := &config.Config{
		Checks: config.ChecksConfig{
			Ping: config.PingDefaultsConfig{
				Count:          5,
				IntervalMs:     2000,
				PacketSize:     64,
				TimeoutSeconds: 10,
			},
		},
	}

	checker := checks.NewPingChecker(cfg)

	t.Run("Ping checker can be created with config", func(t *testing.T) {
		if checker == nil {
			t.Error("Expected checker to be created")
		}

		if checker.Type() != "ping" {
			t.Errorf("Expected type 'ping', got: %s", checker.Type())
		}
	})

	t.Run("Ping check with empty config", func(t *testing.T) {
		check := &storage.Check{
			ID:              1,
			Name:            "Test Ping Check",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8", // Google DNS
			Config:          "",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Enabled:         true,
		}

		// This tests the full flow including config parsing
		result, err := checker.Check(context.Background(), check)
		if err != nil {
			t.Logf("Ping failed (might be expected in some environments): %v", err)
		} else {
			if result == nil {
				t.Error("Expected result to be returned")
			}
			if result.CheckID != check.ID {
				t.Errorf("Expected CheckID %d, got %d", check.ID, result.CheckID)
			}
		}
	})

	t.Run("Ping check with custom config", func(t *testing.T) {
		check := &storage.Check{
			ID:     2,
			Name:   "Test Ping Check with Config",
			Type:   storage.CheckTypePing,
			Target: "1.1.1.1", // Cloudflare DNS
			Config: `{
				"count": 3,
				"interval": 1000,
				"size": 56
			}`,
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Enabled:         true,
		}

		result, err := checker.Check(context.Background(), check)
		if err != nil {
			t.Logf("Ping failed (might be expected in some environments): %v", err)
		} else if result != nil {
			if result.CheckID != check.ID {
				t.Errorf("Expected CheckID %d, got %d", check.ID, result.CheckID)
			}
		}
	})

	t.Run("Ping check with invalid config should fail", func(t *testing.T) {
		check := &storage.Check{
			ID:              3,
			Name:            "Test Invalid Ping Check",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8",
			Config:          `{"count": 0, "interval": 50, "size": 5}`, // All invalid
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Enabled:         true,
		}

		_, err := checker.Check(context.Background(), check)
		if err == nil {
			t.Error("Expected error for invalid config")
		}
	})
}

func TestChecksManagerIntegration(t *testing.T) {
	// Create test application config
	cfg := &config.Config{
		Checks: config.ChecksConfig{
			HTTP: config.HTTPDefaultsConfig{
				Method:         "GET",
				TimeoutSeconds: 30,
				ExpectedStatus: 200,
			},
			Ping: config.PingDefaultsConfig{
				Count:          3,
				IntervalMs:     1000,
				PacketSize:     56,
				TimeoutSeconds: 5,
			},
		},
	}

	manager := checks.NewManager(cfg)

	t.Run("Manager should register all checkers", func(t *testing.T) {
		supportedTypes := manager.GetSupportedTypes()

		expectedTypes := []string{"http", "ping"}
		if len(supportedTypes) != len(expectedTypes) {
			t.Errorf("Expected %d supported types, got %d", len(expectedTypes), len(supportedTypes))
		}

		for _, expectedType := range expectedTypes {
			found := false
			for _, supportedType := range supportedTypes {
				if supportedType == expectedType {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected type %s to be supported", expectedType)
			}
		}
	})

	t.Run("Manager should handle unsupported check type", func(t *testing.T) {
		check := &storage.Check{
			ID:              1,
			Name:            "Unsupported Check",
			Type:            "unsupported",
			Target:          "example.com",
			Config:          "{}",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		_, err := manager.ExecuteCheck(context.Background(), check)
		if err == nil {
			t.Error("Expected error for unsupported check type")
		}
	})

	t.Run("Manager should execute HTTP check", func(t *testing.T) {
		check := &storage.Check{
			ID:              2,
			Name:            "HTTP Check via Manager",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://mockhttp.org/get",
			Config:          `{"expected_status": 200}`,
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		result, err := manager.ExecuteCheck(context.Background(), check)
		if err != nil {
			t.Logf("HTTP check failed (might be expected in some environments): %v", err)
		} else if result != nil {
			if result.CheckID != check.ID {
				t.Errorf("Expected CheckID %d, got %d", check.ID, result.CheckID)
			}
		}
	})

	t.Run("Manager should execute Ping check", func(t *testing.T) {
		check := &storage.Check{
			ID:              3,
			Name:            "Ping Check via Manager",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8",
			Config:          `{"count": 2}`,
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		result, err := manager.ExecuteCheck(context.Background(), check)
		if err != nil {
			t.Logf("Ping check failed (might be expected in some environments): %v", err)
		} else if result != nil {
			if result.CheckID != check.ID {
				t.Errorf("Expected CheckID %d, got %d", check.ID, result.CheckID)
			}
		}
	})
}
