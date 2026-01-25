package tests

import (
	"testing"

	"dideban/internal/storage"
)

func TestValidateAlert(t *testing.T) {
	checkID := int64(1)
	agentID := int64(2)

	t.Run("Valid check alert", func(t *testing.T) {
		alert := &storage.Alert{
			CheckID:       &checkID,
			Type:          storage.AlertTypeTelegram,
			ConditionType: "status_down",
			Config:        `{"token": "test_token", "chat_id": "123456"}`,
			Enabled:       true,
		}

		err := storage.ValidateAlert(alert)
		if err != nil {
			t.Errorf("Expected no error for valid check alert, got: %v", err)
		}

		if alert.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
	})

	t.Run("Valid agent alert with threshold", func(t *testing.T) {
		conditionValue := float64(80.0)
		alert := &storage.Alert{
			AgentID:        &agentID,
			Type:           storage.AlertTypeTelegram,
			ConditionType:  "cpu_usage_high",
			ConditionValue: &conditionValue,
			Config:         `{"token": "test_token", "chat_id": "123456"}`,
			Enabled:        true,
		}

		err := storage.ValidateAlert(alert)
		if err != nil {
			t.Errorf("Expected no error for valid agent alert, got: %v", err)
		}
	})

	t.Run("No CheckID or AgentID should fail", func(t *testing.T) {
		alert := &storage.Alert{
			Type:          storage.AlertTypeTelegram,
			ConditionType: "status_down",
			Enabled:       true,
		}

		err := storage.ValidateAlert(alert)
		if err == nil {
			t.Error("Expected error for alert with no CheckID or AgentID")
		}
	})

	t.Run("Both CheckID and AgentID should fail", func(t *testing.T) {
		alert := &storage.Alert{
			CheckID:       &checkID,
			AgentID:       &agentID,
			Type:          storage.AlertTypeTelegram,
			ConditionType: "status_down",
			Enabled:       true,
		}

		err := storage.ValidateAlert(alert)
		if err == nil {
			t.Error("Expected error for alert with both CheckID and AgentID")
		}
	})

	t.Run("Invalid alert type should fail", func(t *testing.T) {
		alert := &storage.Alert{
			CheckID:       &checkID,
			Type:          "invalid_type",
			ConditionType: "status_down",
			Enabled:       true,
		}

		err := storage.ValidateAlert(alert)
		if err == nil {
			t.Error("Expected error for invalid alert type")
		}
	})

	t.Run("Invalid condition type for check alert should fail", func(t *testing.T) {
		alert := &storage.Alert{
			CheckID:       &checkID,
			Type:          storage.AlertTypeTelegram,
			ConditionType: "cpu_usage_high", // Invalid for check alert
			Enabled:       true,
		}

		err := storage.ValidateAlert(alert)
		if err == nil {
			t.Error("Expected error for invalid condition type for check alert")
		}
	})

	t.Run("Missing condition value for threshold condition should fail", func(t *testing.T) {
		alert := &storage.Alert{
			AgentID:       &agentID,
			Type:          storage.AlertTypeTelegram,
			ConditionType: "cpu_usage_high", // Requires condition value
			Enabled:       true,
		}

		err := storage.ValidateAlert(alert)
		if err == nil {
			t.Error("Expected error for missing condition value")
		}
	})

	t.Run("Invalid condition value should fail", func(t *testing.T) {
		conditionValue := float64(150.0) // Invalid percentage
		alert := &storage.Alert{
			AgentID:        &agentID,
			Type:           storage.AlertTypeTelegram,
			ConditionType:  "cpu_usage_high",
			ConditionValue: &conditionValue,
			Enabled:        true,
		}

		err := storage.ValidateAlert(alert)
		if err == nil {
			t.Error("Expected error for invalid condition value")
		}
	})
}

func TestValidateAlertHistory(t *testing.T) {
	t.Run("Valid alert history", func(t *testing.T) {
		history := &storage.AlertHistory{
			AlertID: 1,
			Title:   "Test Alert",
			Message: "This is a test alert message",
			Status:  storage.AlertStatusSent,
		}

		err := storage.ValidateAlertHistory(history)
		if err != nil {
			t.Errorf("Expected no error for valid alert history, got: %v", err)
		}

		if history.SentAt.IsZero() {
			t.Error("Expected SentAt to be set")
		}
		if history.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
	})

	t.Run("Empty AlertID should fail", func(t *testing.T) {
		history := &storage.AlertHistory{
			AlertID: 0,
			Title:   "Test Alert",
			Message: "Test message",
			Status:  storage.AlertStatusSent,
		}

		err := storage.ValidateAlertHistory(history)
		if err == nil {
			t.Error("Expected error for empty AlertID")
		}
	})

	t.Run("Empty title should fail", func(t *testing.T) {
		history := &storage.AlertHistory{
			AlertID: 1,
			Title:   "",
			Message: "Test message",
			Status:  storage.AlertStatusSent,
		}

		err := storage.ValidateAlertHistory(history)
		if err == nil {
			t.Error("Expected error for empty title")
		}
	})

	t.Run("Title too long should fail", func(t *testing.T) {
		longTitle := make([]byte, 201)
		for i := range longTitle {
			longTitle[i] = 'a'
		}

		history := &storage.AlertHistory{
			AlertID: 1,
			Title:   string(longTitle),
			Message: "Test message",
			Status:  storage.AlertStatusSent,
		}

		err := storage.ValidateAlertHistory(history)
		if err == nil {
			t.Error("Expected error for title too long")
		}
	})

	t.Run("Empty message should fail", func(t *testing.T) {
		history := &storage.AlertHistory{
			AlertID: 1,
			Title:   "Test Alert",
			Message: "",
			Status:  storage.AlertStatusSent,
		}

		err := storage.ValidateAlertHistory(history)
		if err == nil {
			t.Error("Expected error for empty message")
		}
	})

	t.Run("Message too long should fail", func(t *testing.T) {
		longMessage := make([]byte, 5001)
		for i := range longMessage {
			longMessage[i] = 'a'
		}

		history := &storage.AlertHistory{
			AlertID: 1,
			Title:   "Test Alert",
			Message: string(longMessage),
			Status:  storage.AlertStatusSent,
		}

		err := storage.ValidateAlertHistory(history)
		if err == nil {
			t.Error("Expected error for message too long")
		}
	})

	t.Run("Invalid status should fail", func(t *testing.T) {
		history := &storage.AlertHistory{
			AlertID: 1,
			Title:   "Test Alert",
			Message: "Test message",
			Status:  "invalid_status",
		}

		err := storage.ValidateAlertHistory(history)
		if err == nil {
			t.Error("Expected error for invalid status")
		}
	})
}

func TestValidateAdmin(t *testing.T) {
	t.Run("Valid admin", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "testuser",
			Password: "hashed_password_12345678",
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err != nil {
			t.Errorf("Expected no error for valid admin, got: %v", err)
		}
	})

	t.Run("Empty username should fail", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "",
			Password: "hashed_password",
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for empty username")
		}
	})

	t.Run("Username too short should fail", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "ab", // Too short
			Password: "hashed_password",
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for username too short")
		}
	})

	t.Run("Username too long should fail", func(t *testing.T) {
		longUsername := make([]byte, 51)
		for i := range longUsername {
			longUsername[i] = 'a'
		}

		admin := &storage.Admin{
			Username: string(longUsername),
			Password: "hashed_password",
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for username too long")
		}
	})

	t.Run("Invalid username characters should fail", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "test@user", // Invalid characters
			Password: "hashed_password",
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for invalid username characters")
		}
	})

	t.Run("Empty password should fail", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "testuser",
			Password: "",
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for empty password")
		}
	})

	t.Run("Password too short should fail", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "testuser",
			Password: "short", // Too short
			FullName: "Test User",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for password too short")
		}
	})

	t.Run("Empty full name should fail", func(t *testing.T) {
		admin := &storage.Admin{
			Username: "testuser",
			Password: "hashed_password",
			FullName: "",
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for empty full name")
		}
	})

	t.Run("Full name too long should fail", func(t *testing.T) {
		longName := make([]byte, 101)
		for i := range longName {
			longName[i] = 'a'
		}

		admin := &storage.Admin{
			Username: "testuser",
			Password: "hashed_password",
			FullName: string(longName),
		}

		err := storage.ValidateAdmin(admin)
		if err == nil {
			t.Error("Expected error for full name too long")
		}
	})
}

func TestValidateCheckHistory(t *testing.T) {
	responseTime := 150
	statusCode := 200

	t.Run("Valid check history", func(t *testing.T) {
		history := &storage.CheckHistory{
			CheckID:        1,
			Status:         storage.CheckStatusUp,
			ResponseTimeMs: &responseTime,
			StatusCode:     &statusCode,
		}

		err := storage.ValidateCheckHistory(history)
		if err != nil {
			t.Errorf("Expected no error for valid check history, got: %v", err)
		}

		if history.CheckedAt.IsZero() {
			t.Error("Expected CheckedAt to be set")
		}
	})

	t.Run("Empty CheckID should fail", func(t *testing.T) {
		history := &storage.CheckHistory{
			CheckID: 0,
			Status:  storage.CheckStatusUp,
		}

		err := storage.ValidateCheckHistory(history)
		if err == nil {
			t.Error("Expected error for empty CheckID")
		}
	})

	t.Run("Invalid status should fail", func(t *testing.T) {
		history := &storage.CheckHistory{
			CheckID: 1,
			Status:  "invalid_status",
		}

		err := storage.ValidateCheckHistory(history)
		if err == nil {
			t.Error("Expected error for invalid status")
		}
	})

	t.Run("Negative response time should fail", func(t *testing.T) {
		negativeTime := -10
		history := &storage.CheckHistory{
			CheckID:        1,
			Status:         storage.CheckStatusUp,
			ResponseTimeMs: &negativeTime,
		}

		err := storage.ValidateCheckHistory(history)
		if err == nil {
			t.Error("Expected error for negative response time")
		}
	})

	t.Run("Invalid status code should fail", func(t *testing.T) {
		invalidCode := 999
		history := &storage.CheckHistory{
			CheckID:    1,
			Status:     storage.CheckStatusUp,
			StatusCode: &invalidCode,
		}

		err := storage.ValidateCheckHistory(history)
		if err == nil {
			t.Error("Expected error for invalid status code")
		}
	})

	t.Run("Error message too long should fail", func(t *testing.T) {
		longError := make([]byte, 1001)
		for i := range longError {
			longError[i] = 'a'
		}
		longErrorStr := string(longError)

		history := &storage.CheckHistory{
			CheckID:      1,
			Status:       storage.CheckStatusError,
			ErrorMessage: &longErrorStr,
		}

		err := storage.ValidateCheckHistory(history)
		if err == nil {
			t.Error("Expected error for error message too long")
		}
	})
}

func TestValidateAgentHistory(t *testing.T) {
	t.Run("Valid agent history", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:            1,
			CollectDurationMs:  1500,
			CPULoad1:           1.5,
			CPULoad5:           2.0,
			CPULoad15:          1.8,
			CPUUsagePercent:    75.5,
			MemoryTotalMB:      8192,
			MemoryUsedMB:       6144,
			MemoryAvailableMB:  2048,
			MemoryUsagePercent: 75.0,
			DiskTotalGB:        500,
			DiskUsedGB:         350,
			DiskUsagePercent:   70.0,
		}

		err := storage.ValidateAgentHistory(history)
		if err != nil {
			t.Errorf("Expected no error for valid agent history, got: %v", err)
		}

		if history.CollectedAt.IsZero() {
			t.Error("Expected CollectedAt to be set")
		}
	})

	t.Run("Empty AgentID should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:           0,
			CollectDurationMs: 1000,
			MemoryTotalMB:     8192,
			DiskTotalGB:       500,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for empty AgentID")
		}
	})

	t.Run("Negative collect duration should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:           1,
			CollectDurationMs: -100,
			MemoryTotalMB:     8192,
			DiskTotalGB:       500,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for negative collect duration")
		}
	})

	t.Run("Collect duration too long should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:           1,
			CollectDurationMs: 400000, // > 5 minutes
			MemoryTotalMB:     8192,
			DiskTotalGB:       500,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for collect duration too long")
		}
	})

	t.Run("Invalid CPU usage percent should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:         1,
			CPUUsagePercent: 150.0, // > 100%
			MemoryTotalMB:   8192,
			DiskTotalGB:     500,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for invalid CPU usage percent")
		}
	})

	t.Run("Zero memory total should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:       1,
			MemoryTotalMB: 0,
			DiskTotalGB:   500,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for zero memory total")
		}
	})

	t.Run("Invalid memory usage percent should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:            1,
			MemoryTotalMB:      8192,
			MemoryUsagePercent: -10.0, // Negative
			DiskTotalGB:        500,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for invalid memory usage percent")
		}
	})

	t.Run("Zero disk total should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:       1,
			MemoryTotalMB: 8192,
			DiskTotalGB:   0,
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for zero disk total")
		}
	})

	t.Run("Invalid disk usage percent should fail", func(t *testing.T) {
		history := &storage.AgentHistory{
			AgentID:          1,
			MemoryTotalMB:    8192,
			DiskTotalGB:      500,
			DiskUsagePercent: 120.0, // > 100%
		}

		err := storage.ValidateAgentHistory(history)
		if err == nil {
			t.Error("Expected error for invalid disk usage percent")
		}
	})
}
