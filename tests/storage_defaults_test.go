package tests

import (
	"encoding/json"
	"testing"

	"dideban/internal/storage"
)

// TestSetValidationDefaults tests the SetValidationDefaults function and global defaults behavior
func TestSetValidationDefaults(t *testing.T) {
	// Reset defaults before each test
	storage.SetValidationDefaults(nil, nil)

	t.Run("SetValidationDefaults with nil values", func(t *testing.T) {
		storage.SetValidationDefaults(nil, nil)

		// Test HTTP check with empty config - should use hardcoded fallbacks
		check := &storage.Check{
			Name:            "Test HTTP Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Parse the resulting config to verify hardcoded defaults
		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(check.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		// Verify hardcoded defaults
		if httpConfig.Method != "GET" {
			t.Errorf("Expected method 'GET', got '%s'", httpConfig.Method)
		}
		if httpConfig.ExpectedStatus != 200 {
			t.Errorf("Expected status 200, got %d", httpConfig.ExpectedStatus)
		}
		if !httpConfig.FollowRedirects {
			t.Error("Expected FollowRedirects to be true")
		}
		if !httpConfig.VerifySSL {
			t.Error("Expected VerifySSL to be true")
		}
	})

	t.Run("SetValidationDefaults with custom HTTP defaults", func(t *testing.T) {
		customHTTPDefaults := &storage.HTTPDefaultsConfig{
			Method:          "POST",
			Headers:         map[string]string{"Authorization": "Bearer token", "Content-Type": "application/json"},
			Body:            "custom body",
			TimeoutSeconds:  25,
			ExpectedStatus:  201,
			ExpectedContent: "success response",
			FollowRedirects: false,
			VerifySSL:       false,
		}

		storage.SetValidationDefaults(customHTTPDefaults, nil)

		check := &storage.Check{
			Name:            "Test HTTP Check with Custom Defaults",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 120,
			TimeoutSeconds:  60,
			Config:          `{}`, // Empty config should use custom defaults
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(check.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		// Verify custom defaults were applied
		if httpConfig.Method != "POST" {
			t.Errorf("Expected method 'POST', got '%s'", httpConfig.Method)
		}
		if httpConfig.ExpectedStatus != 201 {
			t.Errorf("Expected status 201, got %d", httpConfig.ExpectedStatus)
		}
		if httpConfig.Body != "custom body" {
			t.Errorf("Expected body 'custom body', got '%s'", httpConfig.Body)
		}
		if httpConfig.ExpectedContent != "success response" {
			t.Errorf("Expected content 'success response', got '%s'", httpConfig.ExpectedContent)
		}
		if httpConfig.FollowRedirects {
			t.Error("Expected FollowRedirects to be false")
		}
		if httpConfig.VerifySSL {
			t.Error("Expected VerifySSL to be false")
		}
		if len(httpConfig.Headers) != 2 {
			t.Errorf("Expected 2 headers, got %d", len(httpConfig.Headers))
		}
		if httpConfig.Headers["Authorization"] != "Bearer token" {
			t.Errorf("Expected Authorization header 'Bearer token', got '%s'", httpConfig.Headers["Authorization"])
		}
		if httpConfig.Headers["Content-Type"] != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got '%s'", httpConfig.Headers["Content-Type"])
		}
	})

	t.Run("SetValidationDefaults with custom Ping defaults", func(t *testing.T) {
		customPingDefaults := &storage.PingDefaultsConfig{
			Count:          5,
			IntervalMs:     500,
			PacketSize:     128,
			TimeoutSeconds: 3,
		}

		storage.SetValidationDefaults(nil, customPingDefaults)

		check := &storage.Check{
			Name:            "Test Ping Check with Custom Defaults",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8",
			IntervalSeconds: 90,
			TimeoutSeconds:  45,
			Config:          `{}`, // Empty config should use custom defaults
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		var pingConfig storage.PingCheckConfig
		err = json.Unmarshal([]byte(check.Config), &pingConfig)
		if err != nil {
			t.Errorf("Failed to parse Ping config: %v", err)
		}

		// Verify custom defaults were applied
		if pingConfig.Count != 5 {
			t.Errorf("Expected count 5, got %d", pingConfig.Count)
		}
		if pingConfig.Interval != 500 {
			t.Errorf("Expected interval 500, got %d", pingConfig.Interval)
		}
		if pingConfig.Size != 128 {
			t.Errorf("Expected size 128, got %d", pingConfig.Size)
		}
	})

	t.Run("SetValidationDefaults with both HTTP and Ping defaults", func(t *testing.T) {
		customHTTPDefaults := &storage.HTTPDefaultsConfig{
			Method:          "PUT",
			ExpectedStatus:  204,
			FollowRedirects: false,
			VerifySSL:       true,
		}

		customPingDefaults := &storage.PingDefaultsConfig{
			Count:      10,
			IntervalMs: 200,
			PacketSize: 64,
		}

		storage.SetValidationDefaults(customHTTPDefaults, customPingDefaults)

		// Test HTTP check
		httpCheck := &storage.Check{
			Name:            "Test HTTP with Both Defaults",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheck(httpCheck)
		if err != nil {
			t.Errorf("Expected no error for HTTP check, got: %v", err)
		}

		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(httpCheck.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		if httpConfig.Method != "PUT" {
			t.Errorf("Expected HTTP method 'PUT', got '%s'", httpConfig.Method)
		}
		if httpConfig.ExpectedStatus != 204 {
			t.Errorf("Expected HTTP status 204, got %d", httpConfig.ExpectedStatus)
		}

		// Test Ping check
		pingCheck := &storage.Check{
			Name:            "Test Ping with Both Defaults",
			Type:            storage.CheckTypePing,
			Target:          "1.1.1.1",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err = storage.ValidateCheck(pingCheck)
		if err != nil {
			t.Errorf("Expected no error for Ping check, got: %v", err)
		}

		var pingConfig storage.PingCheckConfig
		err = json.Unmarshal([]byte(pingCheck.Config), &pingConfig)
		if err != nil {
			t.Errorf("Failed to parse Ping config: %v", err)
		}

		if pingConfig.Count != 10 {
			t.Errorf("Expected Ping count 10, got %d", pingConfig.Count)
		}
		if pingConfig.Interval != 200 {
			t.Errorf("Expected Ping interval 200, got %d", pingConfig.Interval)
		}
		if pingConfig.Size != 64 {
			t.Errorf("Expected Ping size 64, got %d", pingConfig.Size)
		}
	})
}

// TestDefaultsOverrideByUserConfig tests that user-provided config overrides defaults
func TestDefaultsOverrideByUserConfig(t *testing.T) {
	// Set custom defaults
	customHTTPDefaults := &storage.HTTPDefaultsConfig{
		Method:          "POST",
		ExpectedStatus:  201,
		FollowRedirects: false,
		VerifySSL:       false,
		Headers:         map[string]string{"Default": "Header"},
		Body:            "default body",
		ExpectedContent: "default content",
	}

	customPingDefaults := &storage.PingDefaultsConfig{
		Count:      5,
		IntervalMs: 500,
		PacketSize: 128,
	}

	storage.SetValidationDefaults(customHTTPDefaults, customPingDefaults)

	t.Run("User HTTP config overrides defaults", func(t *testing.T) {
		userConfig := `{
			"method": "PATCH",
			"expected_status": 202,
			"follow_redirects": true,
			"verify_ssl": true,
			"headers": {"User": "Header", "Custom": "Value"},
			"body": "user body",
			"expected_content": "user content"
		}`

		check := &storage.Check{
			Name:            "Test HTTP Override",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          userConfig,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(check.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		// Verify user config overrode defaults
		if httpConfig.Method != "PATCH" {
			t.Errorf("Expected method 'PATCH', got '%s'", httpConfig.Method)
		}
		if httpConfig.ExpectedStatus != 202 {
			t.Errorf("Expected status 202, got %d", httpConfig.ExpectedStatus)
		}
		if !httpConfig.FollowRedirects {
			t.Error("Expected FollowRedirects to be true (user override)")
		}
		if !httpConfig.VerifySSL {
			t.Error("Expected VerifySSL to be true (user override)")
		}
		if httpConfig.Body != "user body" {
			t.Errorf("Expected body 'user body', got '%s'", httpConfig.Body)
		}
		if httpConfig.ExpectedContent != "user content" {
			t.Errorf("Expected content 'user content', got '%s'", httpConfig.ExpectedContent)
		}

		// Headers: user config merges with defaults, so we expect 3 headers total
		// Default: "Default": "Header" + User: "User": "Header", "Custom": "Value"
		if len(httpConfig.Headers) != 3 {
			t.Errorf("Expected 3 headers (1 default + 2 user), got %d", len(httpConfig.Headers))
		}
		if httpConfig.Headers["Default"] != "Header" {
			t.Errorf("Expected Default header 'Header' (from defaults), got '%s'", httpConfig.Headers["Default"])
		}
		if httpConfig.Headers["User"] != "Header" {
			t.Errorf("Expected User header 'Header', got '%s'", httpConfig.Headers["User"])
		}
		if httpConfig.Headers["Custom"] != "Value" {
			t.Errorf("Expected Custom header 'Value', got '%s'", httpConfig.Headers["Custom"])
		}
	})

	t.Run("User Ping config overrides defaults", func(t *testing.T) {
		userConfig := `{
			"count": 8,
			"interval": 750,
			"size": 256
		}`

		check := &storage.Check{
			Name:            "Test Ping Override",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          userConfig,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		var pingConfig storage.PingCheckConfig
		err = json.Unmarshal([]byte(check.Config), &pingConfig)
		if err != nil {
			t.Errorf("Failed to parse Ping config: %v", err)
		}

		// Verify user config overrode defaults
		if pingConfig.Count != 8 {
			t.Errorf("Expected count 8, got %d", pingConfig.Count)
		}
		if pingConfig.Interval != 750 {
			t.Errorf("Expected interval 750, got %d", pingConfig.Interval)
		}
		if pingConfig.Size != 256 {
			t.Errorf("Expected size 256, got %d", pingConfig.Size)
		}
	})

	t.Run("Partial user config merges with defaults", func(t *testing.T) {
		// User only specifies method, other fields should use defaults
		userConfig := `{"method": "DELETE"}`

		check := &storage.Check{
			Name:            "Test HTTP Partial Override",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          userConfig,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(check.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		// User override
		if httpConfig.Method != "DELETE" {
			t.Errorf("Expected method 'DELETE', got '%s'", httpConfig.Method)
		}

		// Should use defaults for other fields
		if httpConfig.ExpectedStatus != 201 {
			t.Errorf("Expected status 201 (from defaults), got %d", httpConfig.ExpectedStatus)
		}
		if httpConfig.FollowRedirects {
			t.Error("Expected FollowRedirects to be false (from defaults)")
		}
		if httpConfig.VerifySSL {
			t.Error("Expected VerifySSL to be false (from defaults)")
		}
		if httpConfig.Body != "default body" {
			t.Errorf("Expected body 'default body' (from defaults), got '%s'", httpConfig.Body)
		}
		if httpConfig.ExpectedContent != "default content" {
			t.Errorf("Expected content 'default content' (from defaults), got '%s'", httpConfig.ExpectedContent)
		}
		if len(httpConfig.Headers) != 1 {
			t.Errorf("Expected 1 header (from defaults), got %d", len(httpConfig.Headers))
		}
		if httpConfig.Headers["Default"] != "Header" {
			t.Errorf("Expected Default header 'Header' (from defaults), got '%s'", httpConfig.Headers["Default"])
		}
	})
}

// TestDefaultsWithValidateCheckWithDefaults tests the new ValidateCheckWithDefaults function
func TestValidateCheckWithDefaults(t *testing.T) {
	t.Run("ValidateCheckWithDefaults with custom defaults", func(t *testing.T) {
		httpDefaults := &storage.HTTPDefaultsConfig{
			Method:          "PATCH",
			ExpectedStatus:  206,
			FollowRedirects: false,
			VerifySSL:       true,
		}

		pingDefaults := &storage.PingDefaultsConfig{
			Count:      7,
			IntervalMs: 300,
			PacketSize: 96,
		}

		// Test HTTP check
		httpCheck := &storage.Check{
			Name:            "Test HTTP with ValidateCheckWithDefaults",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheckWithDefaults(httpCheck, httpDefaults, pingDefaults)
		if err != nil {
			t.Errorf("Expected no error for HTTP check, got: %v", err)
		}

		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(httpCheck.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		if httpConfig.Method != "PATCH" {
			t.Errorf("Expected method 'PATCH', got '%s'", httpConfig.Method)
		}
		if httpConfig.ExpectedStatus != 206 {
			t.Errorf("Expected status 206, got %d", httpConfig.ExpectedStatus)
		}

		// Test Ping check
		pingCheck := &storage.Check{
			Name:            "Test Ping with ValidateCheckWithDefaults",
			Type:            storage.CheckTypePing,
			Target:          "1.1.1.1",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err = storage.ValidateCheckWithDefaults(pingCheck, httpDefaults, pingDefaults)
		if err != nil {
			t.Errorf("Expected no error for Ping check, got: %v", err)
		}

		var pingConfig storage.PingCheckConfig
		err = json.Unmarshal([]byte(pingCheck.Config), &pingConfig)
		if err != nil {
			t.Errorf("Failed to parse Ping config: %v", err)
		}

		if pingConfig.Count != 7 {
			t.Errorf("Expected count 7, got %d", pingConfig.Count)
		}
		if pingConfig.Interval != 300 {
			t.Errorf("Expected interval 300, got %d", pingConfig.Interval)
		}
		if pingConfig.Size != 96 {
			t.Errorf("Expected size 96, got %d", pingConfig.Size)
		}
	})

	t.Run("ValidateCheckWithDefaults with nil defaults", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test with nil defaults",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheckWithDefaults(check, nil, nil)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Should use hardcoded fallback defaults
		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(check.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		if httpConfig.Method != "GET" {
			t.Errorf("Expected method 'GET' (hardcoded default), got '%s'", httpConfig.Method)
		}
		if httpConfig.ExpectedStatus != 200 {
			t.Errorf("Expected status 200 (hardcoded default), got %d", httpConfig.ExpectedStatus)
		}
	})
}

// TestDefaultsEdgeCases tests edge cases and error conditions
func TestDefaultsEdgeCases(t *testing.T) {
	t.Run("Empty headers in defaults", func(t *testing.T) {
		httpDefaults := &storage.HTTPDefaultsConfig{
			Method:  "GET",
			Headers: map[string]string{}, // Empty headers map
		}

		storage.SetValidationDefaults(httpDefaults, nil)

		check := &storage.Check{
			Name:            "Test Empty Headers",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		var httpConfig storage.HTTPCheckConfig
		err = json.Unmarshal([]byte(check.Config), &httpConfig)
		if err != nil {
			t.Errorf("Failed to parse HTTP config: %v", err)
		}

		if httpConfig.Headers == nil {
			t.Error("Expected headers to be initialized")
		}
		if len(httpConfig.Headers) != 0 {
			t.Errorf("Expected empty headers map, got %d headers", len(httpConfig.Headers))
		}
	})

	t.Run("Zero values in defaults", func(t *testing.T) {
		httpDefaults := &storage.HTTPDefaultsConfig{
			Method:          "", // Empty method
			ExpectedStatus:  0,  // Zero status
			FollowRedirects: false,
			VerifySSL:       false,
		}

		pingDefaults := &storage.PingDefaultsConfig{
			Count:      0, // Zero count
			IntervalMs: 0, // Zero interval
			PacketSize: 0, // Zero size
		}

		storage.SetValidationDefaults(httpDefaults, pingDefaults)

		// Test HTTP - should handle zero/empty values gracefully
		httpCheck := &storage.Check{
			Name:            "Test Zero HTTP Defaults",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheck(httpCheck)
		if err != nil {
			t.Errorf("Expected no error for HTTP check, got: %v", err)
		}

		// Test Ping - should handle zero values gracefully
		pingCheck := &storage.Check{
			Name:            "Test Zero Ping Defaults",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err = storage.ValidateCheck(pingCheck)
		if err != nil {
			t.Errorf("Expected no error for Ping check, got: %v", err)
		}
	})

	t.Run("Unsupported check type with defaults", func(t *testing.T) {
		storage.SetValidationDefaults(&storage.HTTPDefaultsConfig{}, &storage.PingDefaultsConfig{})

		check := &storage.Check{
			Name:            "Test Unsupported Type",
			Type:            "unsupported",
			Target:          "example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Config:          `{}`,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for unsupported check type")
		}
	})
}
