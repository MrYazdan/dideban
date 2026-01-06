package tests

import (
	"testing"

	"dideban/internal/storage"
)

func TestValidateHTTPCheckTarget(t *testing.T) {
	// Create a test check to use the validation
	createHTTPCheck := func(target string) *storage.Check {
		return &storage.Check{
			Name:            "Test HTTP Check",
			Type:            storage.CheckTypeHTTP,
			Target:          target,
			Config:          `{}`,
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Enabled:         true,
		}
	}

	t.Run("Valid HTTPS URL", func(t *testing.T) {
		check := createHTTPCheck("https://example.com")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid HTTPS URL, got: %v", err)
		}
	})

	t.Run("Valid HTTP URL", func(t *testing.T) {
		check := createHTTPCheck("http://example.com")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid HTTP URL, got: %v", err)
		}
	})

	t.Run("URL without scheme should work (auto HTTPS)", func(t *testing.T) {
		check := createHTTPCheck("example.com")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for URL without scheme, got: %v", err)
		}
	})

	t.Run("URL with path", func(t *testing.T) {
		check := createHTTPCheck("https://example.com/api/health")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for URL with path, got: %v", err)
		}
	})

	t.Run("URL with query parameters", func(t *testing.T) {
		check := createHTTPCheck("https://example.com/api?param=value")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for URL with query parameters, got: %v", err)
		}
	})

	t.Run("URL with port", func(t *testing.T) {
		check := createHTTPCheck("https://example.com:8080")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for URL with port, got: %v", err)
		}
	})

	t.Run("IP address URL", func(t *testing.T) {
		check := createHTTPCheck("https://192.168.1.1")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for IP address URL, got: %v", err)
		}
	})

	t.Run("IPv6 URL", func(t *testing.T) {
		check := createHTTPCheck("https://[2001:db8::1]")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for IPv6 URL, got: %v", err)
		}
	})

	t.Run("Localhost URL", func(t *testing.T) {
		check := createHTTPCheck("http://localhost:3000")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for localhost URL, got: %v", err)
		}
	})

	t.Run("Invalid scheme should fail", func(t *testing.T) {
		check := createHTTPCheck("ftp://example.com")
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for invalid scheme")
		}
	})

	t.Run("Empty host should fail", func(t *testing.T) {
		check := createHTTPCheck("https://")
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for empty host")
		}
	})

	t.Run("Invalid URL format should fail", func(t *testing.T) {
		check := createHTTPCheck("not-a-url")
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for invalid URL format")
		}
	})

	t.Run("Hostname too long should fail", func(t *testing.T) {
		longHostname := make([]byte, 254)
		for i := range longHostname {
			longHostname[i] = 'a'
		}
		longHostname[253] = 'm' // .com
		longHostname[252] = 'o'
		longHostname[251] = 'c'
		longHostname[250] = '.'

		check := createHTTPCheck("https://" + string(longHostname))
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for hostname too long")
		}
	})

	t.Run("Invalid hostname format should fail", func(t *testing.T) {
		check := createHTTPCheck("https://example..com")
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for invalid hostname format")
		}
	})
}

func TestValidatePingCheckTarget(t *testing.T) {
	// Create a test check to use the validation
	createPingCheck := func(target string) *storage.Check {
		return &storage.Check{
			Name:            "Test Ping Check",
			Type:            storage.CheckTypePing,
			Target:          target,
			Config:          `{}`,
			IntervalSeconds: 60,
			TimeoutSeconds:  30,
			Enabled:         true,
		}
	}

	t.Run("Valid IPv4 address", func(t *testing.T) {
		check := createPingCheck("127.0.0.1")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid IPv4 address, got: %v", err)
		}
	})

	t.Run("Valid IPv6 address", func(t *testing.T) {
		check := createPingCheck("2001:db8::1")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid IPv6 address, got: %v", err)
		}
	})

	t.Run("Valid hostname", func(t *testing.T) {
		check := createPingCheck("example.com")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid hostname, got: %v", err)
		}
	})

	t.Run("Valid subdomain", func(t *testing.T) {
		check := createPingCheck("api.example.com")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for valid subdomain, got: %v", err)
		}
	})

	t.Run("Valid hostname with hyphens", func(t *testing.T) {
		check := createPingCheck("my-server.example-domain.com")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for hostname with hyphens, got: %v", err)
		}
	})

	t.Run("Localhost", func(t *testing.T) {
		check := createPingCheck("localhost")
		err := storage.ValidateCheck(check)
		if err != nil {
			t.Errorf("Expected no error for localhost, got: %v", err)
		}
	})

	t.Run("Public DNS servers", func(t *testing.T) {
		targets := []string{
			"8.8.8.8",        // Google DNS
			"1.1.1.1",        // Cloudflare DNS
			"208.67.222.222", // OpenDNS
		}

		for _, target := range targets {
			check := createPingCheck(target)
			err := storage.ValidateCheck(check)
			if err != nil {
				t.Errorf("Expected no error for %s, got: %v", target, err)
			}
		}
	})

	t.Run("Empty target should fail", func(t *testing.T) {
		check := createPingCheck("")
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for empty target")
		}
	})

	t.Run("Hostname too long should fail", func(t *testing.T) {
		longHostname := make([]byte, 254)
		for i := range longHostname {
			longHostname[i] = 'a'
		}

		check := createPingCheck(string(longHostname))
		err := storage.ValidateCheck(check)
		if err == nil {
			t.Error("Expected error for hostname too long")
		}
	})

	t.Run("Invalid hostname format should fail", func(t *testing.T) {
		invalidHostnames := []string{
			"example..com", // Double dots
			".example.com", // Starting with dot
			"example.com.", // Ending with dot
			"exam@ple.com", // Invalid character
			"example .com", // Space
			"example#com",  // Hash
		}

		for _, hostname := range invalidHostnames {
			check := createPingCheck(hostname)
			err := storage.ValidateCheck(check)
			if err == nil {
				t.Errorf("Expected error for invalid hostname: %s", hostname)
			}
		}
	})

	t.Run("Invalid IP addresses should fail", func(t *testing.T) {
		invalidIPs := []string{
			"256.1.1.1",      // Invalid IPv4 octet
			"192.168.1",      // Incomplete IPv4
			"192.168.1.1.1",  // Too many octets
			"gggg::1",        // Invalid IPv6 hex
			"2001:db8::gggg", // Invalid IPv6 hex
		}

		for _, ip := range invalidIPs {
			check := createPingCheck(ip)
			err := storage.ValidateCheck(check)
			if err == nil {
				t.Errorf("Expected error for invalid IP: %s", ip)
			}
		}
	})
}

func TestCheckTypeValidation(t *testing.T) {
	t.Run("Valid check types", func(t *testing.T) {
		validTypes := []string{
			storage.CheckTypeHTTP,
			storage.CheckTypePing,
		}

		for _, checkType := range validTypes {
			if !storage.IsValidCheckType(checkType) {
				t.Errorf("Expected %s to be valid check type", checkType)
			}
		}
	})

	t.Run("Invalid check types", func(t *testing.T) {
		invalidTypes := []string{
			"",
			"invalid",
			"HTTP", // Wrong case
			"PING", // Wrong case
			"tcp",
			"udp",
			"dns",
		}

		for _, checkType := range invalidTypes {
			if storage.IsValidCheckType(checkType) {
				t.Errorf("Expected %s to be invalid check type", checkType)
			}
		}
	})
}

func TestCheckStatusValidation(t *testing.T) {
	t.Run("Valid check statuses", func(t *testing.T) {
		validStatuses := []string{
			storage.CheckStatusUp,
			storage.CheckStatusDown,
			storage.CheckStatusError,
			storage.CheckStatusTimeout,
		}

		for _, status := range validStatuses {
			if !storage.IsValidCheckStatus(status) {
				t.Errorf("Expected %s to be valid check status", status)
			}
		}
	})

	t.Run("Invalid check statuses", func(t *testing.T) {
		invalidStatuses := []string{
			"",
			"invalid",
			"UP",   // Wrong case
			"DOWN", // Wrong case
			"ok",
			"fail",
			"pending",
		}

		for _, status := range invalidStatuses {
			if storage.IsValidCheckStatus(status) {
				t.Errorf("Expected %s to be invalid check status", status)
			}
		}
	})
}

func TestAlertTypeValidation(t *testing.T) {
	t.Run("Valid alert types", func(t *testing.T) {
		validTypes := []string{
			storage.AlertTypeTelegram,
			storage.AlertTypeBale,
			storage.AlertTypeEmail,
			storage.AlertTypeWebhook,
		}

		for _, alertType := range validTypes {
			if !storage.IsValidAlertType(alertType) {
				t.Errorf("Expected %s to be valid alert type", alertType)
			}
		}
	})

	t.Run("Invalid alert types", func(t *testing.T) {
		invalidTypes := []string{
			"",
			"invalid",
			"TELEGRAM", // Wrong case
			"sms",
			"slack",
			"discord",
		}

		for _, alertType := range invalidTypes {
			if storage.IsValidAlertType(alertType) {
				t.Errorf("Expected %s to be invalid alert type", alertType)
			}
		}
	})
}

func TestAlertStatusValidation(t *testing.T) {
	t.Run("Valid alert statuses", func(t *testing.T) {
		validStatuses := []string{
			storage.AlertStatusSent,
			storage.AlertStatusFailed,
			storage.AlertStatusPending,
		}

		for _, status := range validStatuses {
			if !storage.IsValidAlertStatus(status) {
				t.Errorf("Expected %s to be valid alert status", status)
			}
		}
	})

	t.Run("Invalid alert statuses", func(t *testing.T) {
		invalidStatuses := []string{
			"",
			"invalid",
			"SENT", // Wrong case
			"success",
			"error",
			"queued",
		}

		for _, status := range invalidStatuses {
			if storage.IsValidAlertStatus(status) {
				t.Errorf("Expected %s to be invalid alert status", status)
			}
		}
	})
}

func TestAgentStatusValidation(t *testing.T) {
	t.Run("Valid agent statuses", func(t *testing.T) {
		validStatuses := []string{
			storage.AgentStatusOnline,
			storage.AgentStatusOffline,
		}

		for _, status := range validStatuses {
			if !storage.IsValidAgentStatus(status) {
				t.Errorf("Expected %s to be valid agent status", status)
			}
		}
	})

	t.Run("Invalid agent statuses", func(t *testing.T) {
		invalidStatuses := []string{
			"",
			"invalid",
			"ONLINE", // Wrong case
			"up",
			"down",
			"connected",
			"disconnected",
		}

		for _, status := range invalidStatuses {
			if storage.IsValidAgentStatus(status) {
				t.Errorf("Expected %s to be invalid agent status", status)
			}
		}
	})
}
