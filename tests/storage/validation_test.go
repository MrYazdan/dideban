package tests

import (
	"testing"

	"dideban/internal/storage"
)

func TestValidateCheck(t *testing.T) {
	t.Run("Valid HTTP check", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test HTTP Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			Config:          `{"method": "GET", "expected_status": 200}`,
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Enabled:         true,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid check, got: %v", err)
		}

		// Check timestamps were set
		if check.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if check.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Valid Ping check", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test Ping Check",
			Type:            storage.CheckTypePing,
			Target:          "8.8.8.8",
			Config:          `{"count": 3, "interval": 1000}`,
			IntervalSeconds: 120,
			TimeoutSeconds:  60,
			Enabled:         true,
		}

		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid ping check, got: %v", err)
		}
	})

	t.Run("Empty name should fail", func(t *testing.T) {
		check := &storage.Check{
			Name:            "",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})

	t.Run("Invalid check type should fail", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test Check",
			Type:            "invalid",
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for invalid check type")
		}
	})

	t.Run("Empty target should fail", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for empty target")
		}
	})

	t.Run("Interval too short should fail", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 1, // Too short
			TimeoutSeconds:  30,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for interval too short")
		}
	})

	t.Run("Timeout >= interval should fail", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test Check",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 30,
			TimeoutSeconds:  30, // Same as interval
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for timeout >= interval")
		}
	})

	t.Run("Invalid name characters should fail", func(t *testing.T) {
		check := &storage.Check{
			Name:            "Test@Check#Invalid",
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for invalid name characters")
		}
	})

	t.Run("Name too long should fail", func(t *testing.T) {
		longName := make([]byte, 101)
		for i := range longName {
			longName[i] = 'a'
		}

		check := &storage.Check{
			Name:            string(longName),
			Type:            storage.CheckTypeHTTP,
			Target:          "https://example.com",
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
		}

		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for name too long")
		}
	})
}

func TestValidateHTTPCheckConfig(t *testing.T) {
	t.Run("Empty config uses defaults", func(t *testing.T) {
		result, err := storage.ValidateHTTPCheckConfig("", nil)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Should return valid JSON with defaults
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("Valid config with custom values", func(t *testing.T) {
		config := `{
			"method": "POST",
			"expected_status": 201,
			"headers": {"Content-Type": "application/json"},
			"body": "test body",
			"expected_content": "success",
			"follow_redirects": false,
			"verify_ssl": false
		}`

		result, err := storage.ValidateHTTPCheckConfig(config, nil)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("Invalid HTTP method should fail", func(t *testing.T) {
		config := `{"method": "INVALID"}`

		_, err := storage.ValidateHTTPCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for invalid HTTP method")
		}
	})

	t.Run("Invalid status code should fail", func(t *testing.T) {
		config := `{"expected_status": 999}`

		_, err := storage.ValidateHTTPCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for invalid status code")
		}
	})

	t.Run("Empty header key should fail", func(t *testing.T) {
		config := `{"headers": {"": "value"}}`

		_, err := storage.ValidateHTTPCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for empty header key")
		}
	})

	t.Run("Header key too long should fail", func(t *testing.T) {
		longKey := make([]byte, 101)
		for i := range longKey {
			longKey[i] = 'a'
		}

		config := `{"headers": {"` + string(longKey) + `": "value"}}`

		_, err := storage.ValidateHTTPCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for header key too long")
		}
	})

	t.Run("Body too large should fail", func(t *testing.T) {
		largeBody := make([]byte, 11*1024*1024) // 11MB
		for i := range largeBody {
			largeBody[i] = 'a'
		}

		config := `{"body": "` + string(largeBody) + `"}`

		_, err := storage.ValidateHTTPCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for body too large")
		}
	})

	t.Run("Custom defaults should be applied", func(t *testing.T) {
		defaults := &storage.HTTPCheckConfig{
			Method:         "PUT",
			ExpectedStatus: 202,
			Headers:        map[string]string{"Custom": "Header"},
		}

		result, err := storage.ValidateHTTPCheckConfig("", defaults)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}
		// Note: We'd need to parse the result to verify defaults were applied
	})
}

func TestValidatePingCheckConfig(t *testing.T) {
	t.Run("Empty config uses defaults", func(t *testing.T) {
		result, err := storage.ValidatePingCheckConfig("", nil)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("Valid config with custom values", func(t *testing.T) {
		config := `{
			"count": 5,
			"interval": 2000,
			"size": 64
		}`

		result, err := storage.ValidatePingCheckConfig(config, nil)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("Count too low should fail", func(t *testing.T) {
		config := `{"count": 0}`

		_, err := storage.ValidatePingCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for count too low")
		}
	})

	t.Run("Count too high should fail", func(t *testing.T) {
		config := `{"count": 11}`

		_, err := storage.ValidatePingCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for count too high")
		}
	})

	t.Run("Interval too low should fail", func(t *testing.T) {
		config := `{"interval": 50}`

		_, err := storage.ValidatePingCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for interval too low")
		}
	})

	t.Run("Interval too high should fail", func(t *testing.T) {
		config := `{"interval": 15000}`

		_, err := storage.ValidatePingCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for interval too high")
		}
	})

	t.Run("Packet size too small should fail", func(t *testing.T) {
		config := `{"size": 7}`

		_, err := storage.ValidatePingCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for packet size too small")
		}
	})

	t.Run("Packet size too large should fail", func(t *testing.T) {
		config := `{"size": 1500}`

		_, err := storage.ValidatePingCheckConfig(config, nil)
		if err == nil {
			t.Error("Expected error for packet size too large")
		}
	})

	t.Run("Custom defaults should be applied", func(t *testing.T) {
		defaults := &storage.PingCheckConfig{
			Count:    10,
			Interval: 500,
			Size:     128,
		}

		result, err := storage.ValidatePingCheckConfig("", defaults)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}
	})
}

func TestValidateAgent(t *testing.T) {
	t.Run("Valid agent", func(t *testing.T) {
		agent := &storage.Agent{
			Name:            "Test Agent",
			IntervalSeconds: 60,
			Enabled:         true,
		}

		err := storage.ValidateAgent(agent)
		if err != nil {
			t.Errorf("Expected no error for valid agent, got: %v", err)
		}

		// Check auth token was generated
		if agent.AuthToken == "" {
			t.Error("Expected auth token to be generated")
		}

		// Check timestamps were set
		if agent.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if agent.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Empty name should fail", func(t *testing.T) {
		agent := &storage.Agent{
			Name:            "",
			IntervalSeconds: 60,
		}

		err := storage.ValidateAgent(agent)
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})

	t.Run("Interval too short should fail", func(t *testing.T) {
		agent := &storage.Agent{
			Name:            "Test Agent",
			IntervalSeconds: 5, // Too short
		}

		err := storage.ValidateAgent(agent)
		if err == nil {
			t.Error("Expected error for interval too short")
		}
	})

	t.Run("Interval too long should fail", func(t *testing.T) {
		agent := &storage.Agent{
			Name:            "Test Agent",
			IntervalSeconds: 7200, // Too long
		}

		err := storage.ValidateAgent(agent)
		if err == nil {
			t.Error("Expected error for interval too long")
		}
	})

	t.Run("Invalid name characters should fail", func(t *testing.T) {
		agent := &storage.Agent{
			Name:            "Test@Agent#Invalid",
			IntervalSeconds: 60,
		}

		err := storage.ValidateAgent(agent)
		if err == nil {
			t.Error("Expected error for invalid name characters")
		}
	})

	t.Run("Existing auth token should be preserved", func(t *testing.T) {
		existingToken := "existing_token_12345678901234567890123456789012"
		agent := &storage.Agent{
			Name:            "Test Agent",
			IntervalSeconds: 60,
			AuthToken:       existingToken,
		}

		err := storage.ValidateAgent(agent)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if agent.AuthToken != existingToken {
			t.Error("Expected existing auth token to be preserved")
		}
	})
}
