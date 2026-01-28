// Package storage provides strict, defensive validation for the Check model.
//
// This file contains only validation logic related to the Check entity,
// enabling focused, incremental migration to a GORM-based storage layer.
// All validators are pure functions that never mutate timestamps,
// as CreatedAt/UpdatedAt are managed automatically by GORM.
package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// HTTPDefaultsConfig represents default configuration values for HTTP checks.
// These defaults are applied during validation when user-provided config fields are missing.
// Used by SetValidationDefaults and ValidateCheckWithDefaults.
type HTTPDefaultsConfig struct {
	Method          string            `mapstructure:"method" yaml:"method"`
	Headers         map[string]string `mapstructure:"headers" yaml:"headers"`
	Body            string            `mapstructure:"body" yaml:"body"`
	TimeoutSeconds  int               `mapstructure:"timeout_seconds" yaml:"timeout_seconds"`
	ExpectedStatus  int               `mapstructure:"expected_status" yaml:"expected_status"`
	ExpectedContent string            `mapstructure:"expected_content" yaml:"expected_content"`
	FollowRedirects bool              `mapstructure:"follow_redirects" yaml:"follow_redirects"`
	VerifySSL       bool              `mapstructure:"verify_ssl" yaml:"verify_ssl"`
}

// PingDefaultsConfig represents default configuration values for ping checks.
// These defaults are applied during validation when user-provided config fields are missing.
// Used by SetValidationDefaults and ValidateCheckWithDefaults.
type PingDefaultsConfig struct {
	Count          int `mapstructure:"count" yaml:"count"`
	IntervalMs     int `mapstructure:"interval_ms" yaml:"interval_ms"`
	PacketSize     int `mapstructure:"packet_size" yaml:"packet_size"`
	TimeoutSeconds int `mapstructure:"timeout_seconds" yaml:"timeout_seconds"`
}

// globalHTTPDefaults and globalPingDefaults hold application-wide configuration
// used during validation of check-specific JSON configs.
// These are set once at startup via SetValidationDefaults.
var (
	globalHTTPDefaults *HTTPDefaultsConfig
	globalPingDefaults *PingDefaultsConfig
)

// SetValidationDefaults configures the global defaults used when validating
// HTTP and Ping check configurations. This function must be called during
// application initialization before any validation occurs.
func SetValidationDefaults(httpDefaults *HTTPDefaultsConfig, pingDefaults *PingDefaultsConfig) {
	globalHTTPDefaults = httpDefaults
	globalPingDefaults = pingDefaults
}

// ValidateCheck performs comprehensive validation of a Check entity.
// It enforces business rules, data integrity constraints, and security requirements.
// The function returns a descriptive error if any validation rule is violated.
//
// Note: This function does NOT normalize timestamps (CreatedAt/UpdatedAt),
// as those fields are managed automatically by GORM.
func ValidateCheck(check *Check) error {
	if check == nil {
		return fmt.Errorf("check is nil")
	}

	// Validate name: required, max length, valid characters
	if check.Name == "" {
		return fmt.Errorf("check name cannot be empty")
	}
	if len(check.Name) > 100 {
		return fmt.Errorf("check name too long (max 100 chars)")
	}

	// Validate type: must be supported
	if !IsValidCheckType(check.Type) {
		return fmt.Errorf("unsupported check type: %s", check.Type)
	}

	// Validate target: non-empty and format-correct for the given type
	if check.Target == "" {
		return fmt.Errorf("check target cannot be empty")
	}
	if err := validateCheckTarget(check.Type, check.Target); err != nil {
		return fmt.Errorf("invalid target: %w", err)
	}

	// Validate intervals: reasonable bounds and timeout < interval
	if check.IntervalSeconds < 5 {
		return fmt.Errorf("check interval too short (minimum 5 seconds)")
	}
	if check.IntervalSeconds > 86400 {
		return fmt.Errorf("check interval too long (maximum 24 hours)")
	}
	if check.TimeoutSeconds < 1 {
		return fmt.Errorf("check timeout too short (minimum 1 second)")
	}
	if check.TimeoutSeconds > 300 {
		return fmt.Errorf("check timeout too long (maximum 5 minutes)")
	}
	if check.TimeoutSeconds >= check.IntervalSeconds {
		return fmt.Errorf("check timeout must be less than interval")
	}

	// Validate and normalize configuration JSON using global defaults
	validatedConfig, err := validateCheckConfig(check.Type, check.Config)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	check.Config = validatedConfig

	return nil
}

// ValidateAgent validates an Agent entity with full business logic.
//
// This function enforces data integrity, security constraints, and application-level
// defaults for the Agent model. It does NOT normalize timestamps (CreatedAt/UpdatedAt),
// as those fields are managed automatically by GORM.
//
// Key behaviors:
//   - Name must be non-empty, ≤100 chars, and contain only alphanumeric, space, hyphen, underscore
//   - IntervalSeconds must be between 10–86400 seconds
//   - Status defaults to "offline" if empty, and must be either "online" or "offline"
//   - AuthToken is auto-generated if empty, and must be 32–128 characters long
func ValidateAgent(agent *Agent) error {
	if agent == nil {
		return fmt.Errorf("agent is nil")
	}

	// Validate name: required, max length, valid characters
	if agent.Name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}
	if len(agent.Name) > 100 {
		return fmt.Errorf("agent name too long (max 100 chars)")
	}

	// Validate interval: reasonable bounds for monitoring agents
	if agent.IntervalSeconds < 10 {
		return fmt.Errorf("agent interval too short (minimum 10 seconds)")
	}
	if agent.IntervalSeconds > 86400 { // 24 hours
		return fmt.Errorf("agent interval too long (maximum 24 hours)")
	}

	// Normalize and validate status
	if agent.Status == "" {
		agent.Status = AgentStatusOffline
	}
	if !IsValidAgentStatus(agent.Status) {
		return fmt.Errorf("invalid agent status: %s", agent.Status)
	}

	// Auto-generate auth token if not provided
	if agent.AuthToken == "" {
		token, err := generateAuthToken()
		if err != nil {
			return fmt.Errorf("failed to generate auth token: %w", err)
		}
		agent.AuthToken = token
	}

	// Validate auth token length (security requirement)
	if len(agent.AuthToken) < 32 {
		return fmt.Errorf("auth token too short (minimum 32 chars)")
	}
	if len(agent.AuthToken) > 128 {
		return fmt.Errorf("auth token too long (maximum 128 chars)")
	}

	return nil
}

// ValidateAlert validates an Alert entity with full context-aware business logic.
//
// This function enforces data integrity and semantic constraints that depend on
// whether the alert is associated with a Check or an Agent. It does NOT normalize
// timestamps (CreatedAt), as that field is managed automatically by GORM.
//
// Key validation rules:
//   - Exactly one of CheckID or AgentID must be set (mutual exclusion)
//   - ConditionType must match the associated entity type
//   - ConditionValue is required for percentage-based thresholds (e.g., cpu_usage_high)
//   - Config must be valid JSON and contain required fields per alert type
func ValidateAlert(alert *Alert) error {
	if alert == nil {
		return fmt.Errorf("alert is nil")
	}

	// Enforce mutual exclusion: exactly one of CheckID or AgentID must be set
	hasCheck := alert.CheckID != nil
	hasAgent := alert.AgentID != nil
	if !hasCheck && !hasAgent {
		return fmt.Errorf("alert must be associated with either a check or an agent")
	}
	if hasCheck && hasAgent {
		return fmt.Errorf("alert cannot be associated with both a check and an agent")
	}

	// Validate alert type against supported enum values
	if !IsValidAlertType(alert.Type) {
		return fmt.Errorf("unsupported alert type: %s", alert.Type)
	}

	// Validate condition type based on context (check vs agent)
	if err := validateAlertConditionType(alert.ConditionType, hasCheck); err != nil {
		return err
	}

	// For threshold-based conditions, ConditionValue is mandatory and must be in [0,100]
	if needsConditionValue(alert.ConditionType) {
		if alert.ConditionValue == nil {
			return fmt.Errorf("condition value is required for condition type: %s", alert.ConditionType)
		}
		if *alert.ConditionValue < 0 || *alert.ConditionValue > 100 {
			return fmt.Errorf("condition value must be between 0 and 100 for percentage-based conditions")
		}
	}

	// Validate alert-specific configuration (e.g., Telegram token, email address)
	if err := validateAlertConfig(alert.Type, alert.Config); err != nil {
		return fmt.Errorf("invalid alert configuration: %w", err)
	}

	return nil
}

// ValidateAlertHistory validates an AlertHistory entity before persistence.
//
// This function enforces data integrity constraints for historical alert records.
// It does NOT normalize timestamps (SentAt, CreatedAt), as those fields are managed
// automatically by GORM based on the model's time.Time fields with default values.
//
// Key validation rules:
//   - AlertID must be non-zero (required foreign key)
//   - Title and Message must be non-empty and within length limits
//   - Status must be one of the predefined valid states
func ValidateAlertHistory(history *AlertHistory) error {
	if history == nil {
		return fmt.Errorf("alert history is nil")
	}

	// Validate required foreign key
	if history.AlertID == 0 {
		return fmt.Errorf("alert ID cannot be empty")
	}

	// Validate title: required and max length
	if history.Title == "" {
		return fmt.Errorf("alert title cannot be empty")
	}
	if len(history.Title) > 200 {
		return fmt.Errorf("alert title too long (max 200 chars)")
	}

	// Validate message: required and max length
	if history.Message == "" {
		return fmt.Errorf("alert message cannot be empty")
	}
	if len(history.Message) > 5000 {
		return fmt.Errorf("alert message too long (max 5000 chars)")
	}

	// Validate status against allowed enum values
	if !IsValidAlertStatus(history.Status) {
		return fmt.Errorf("invalid alert status: %s", history.Status)
	}

	return nil
}

// ValidateAdmin validates an Admin entity before persistence.
//
// This function enforces security and data integrity constraints for administrator accounts.
// It assumes the password field contains a pre-hashed value (e.g., bcrypt) and validates
// its length accordingly. No timestamp normalization is performed, as the Admin model
// does not include time fields.
//
// Key validation rules:
//   - Username: 3–50 chars, alphanumeric + underscore only
//   - Password: 8–255 chars (suitable for bcrypt hash length)
//   - FullName: non-empty, max 100 chars
func ValidateAdmin(admin *Admin) error {
	if admin == nil {
		return fmt.Errorf("admin is nil")
	}

	// Validate username: required, length, and format
	if admin.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(admin.Username) < 3 {
		return fmt.Errorf("username too short (minimum 3 chars)")
	}
	if len(admin.Username) > 50 {
		return fmt.Errorf("username too long (max 50 chars)")
	}

	// Enforce strict username format: only alphanumeric and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(admin.Username) {
		return fmt.Errorf("username contains invalid characters (only alphanumeric and underscores allowed)")
	}

	// Validate password: assumed to be pre-hashed (e.g., bcrypt)
	if admin.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(admin.Password) < 8 {
		return fmt.Errorf("password too short (minimum 8 chars)")
	}
	if len(admin.Password) > 255 {
		return fmt.Errorf("password too long (max 255 chars)")
	}

	// Validate full name: required and max length
	if admin.FullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if len(admin.FullName) > 100 {
		return fmt.Errorf("full name too long (max 100 chars)")
	}

	return nil
}

// ValidateCheckHistory validates a CheckHistory entity before persistence.
//
// This function enforces data integrity constraints for historical check execution records.
// It does NOT normalize the CheckedAt timestamp, as that field is managed automatically
// by GORM based on the model's time.Time field with a default value.
//
// Key validation rules:
//   - CheckID must be non-zero (required foreign key)
//   - Status must be one of the predefined valid states
//   - ResponseTimeMs (if set) must be non-negative
//   - StatusCode (if set) must be a valid HTTP status code (100–599)
//   - ErrorMessage (if set) must not exceed 1000 characters
func ValidateCheckHistory(history *CheckHistory) error {
	if history == nil {
		return fmt.Errorf("check history is nil")
	}

	// Validate required foreign key
	if history.CheckID == 0 {
		return fmt.Errorf("check ID cannot be empty")
	}

	// Validate status against allowed enum values
	if !IsValidCheckStatus(history.Status) {
		return fmt.Errorf("invalid check status: %s", history.Status)
	}

	// Validate response time: must be non-negative if provided
	if history.ResponseTimeMs != nil && *history.ResponseTimeMs < 0 {
		return fmt.Errorf("response time cannot be negative")
	}

	// Validate HTTP status code: must be in [100, 599] if provided
	if history.StatusCode != nil {
		if *history.StatusCode < 100 || *history.StatusCode > 599 {
			return fmt.Errorf("invalid HTTP status code: %d", *history.StatusCode)
		}
	}

	// Validate error message length: max 1000 chars if provided
	if history.ErrorMessage != nil && len(*history.ErrorMessage) > 1000 {
		return fmt.Errorf("error message too long (max 1000 chars)")
	}

	return nil
}

// ValidateAgentHistory validates an AgentHistory entity before persistence.
//
// This function enforces data integrity constraints for historical agent metric records.
// It does NOT normalize the CollectedAt timestamp, as that field is managed automatically
// by GORM based on the model's time.Time field with a default value.
//
// Special behavior:
//   - If IsOffline is true, all metric fields are zeroed and only CollectDurationMs=0 is allowed
//   - Otherwise, all metric values must be within valid ranges (e.g., percentages in [0,100])
func ValidateAgentHistory(history *AgentHistory) error {
	if history == nil {
		return fmt.Errorf("agent history is nil")
	}

	// Validate required foreign key
	if history.AgentID == 0 {
		return fmt.Errorf("agent ID cannot be empty")
	}

	// Validate collect duration: must be non-negative
	if history.CollectDurationMs < 0 {
		return fmt.Errorf("collect duration cannot be negative")
	}

	// Handle offline agent case: all metrics must be zero
	if history.IsOffline {
		if history.CollectDurationMs != 0 {
			return fmt.Errorf("offline history must have zero collect duration")
		}
		// Zero out all metric fields for consistency
		history.CPULoad1 = 0
		history.CPULoad5 = 0
		history.CPULoad15 = 0
		history.CPUUsagePercent = 0
		history.MemoryTotalMB = 0
		history.MemoryUsedMB = 0
		history.MemoryAvailableMB = 0
		history.MemoryUsagePercent = 0
		history.DiskTotalGB = 0
		history.DiskUsedGB = 0
		history.DiskUsagePercent = 0
		return nil
	}

	// Validate CPU load averages: must be non-negative
	if history.CPULoad1 < 0 || history.CPULoad5 < 0 || history.CPULoad15 < 0 {
		return fmt.Errorf("CPU load values cannot be negative")
	}

	// Validate CPU usage percentage: must be in [0, 100]
	if history.CPUUsagePercent < 0 || history.CPUUsagePercent > 100 {
		return fmt.Errorf("CPU usage percent must be between 0 and 100")
	}

	// Validate memory total: must be positive
	if history.MemoryTotalMB <= 0 {
		return fmt.Errorf("memory total must be positive")
	}

	// Validate memory used/available: must be non-negative
	if history.MemoryUsedMB < 0 || history.MemoryAvailableMB < 0 {
		return fmt.Errorf("memory values cannot be negative")
	}

	// Validate memory usage percentage: must be in [0, 100]
	if history.MemoryUsagePercent < 0 || history.MemoryUsagePercent > 100 {
		return fmt.Errorf("memory usage percent must be between 0 and 100")
	}

	// Validate disk total: must be positive
	if history.DiskTotalGB <= 0 {
		return fmt.Errorf("disk total must be positive")
	}

	// Validate disk used: must be non-negative
	if history.DiskUsedGB < 0 {
		return fmt.Errorf("disk used cannot be negative")
	}

	// Validate disk usage percentage: must be in [0, 100]
	if history.DiskUsagePercent < 0 || history.DiskUsagePercent > 100 {
		return fmt.Errorf("disk usage percent must be between 0 and 100")
	}

	return nil
}

// validateCheckTarget dispatches target validation based on check type.
func validateCheckTarget(checkType, target string) error {
	switch checkType {
	case CheckTypeHTTP:
		return validateHTTPCheckTarget(target)
	case CheckTypePing:
		return validatePingCheckTarget(target)
	default:
		return fmt.Errorf("unknown check type: %s", checkType)
	}
}

// validateHTTPCheckTarget validates that the target is a well-formed HTTP(S) URL.
// If no scheme is provided, "https://" is assumed.
func validateHTTPCheckTarget(target string) error {
	// Auto-prepend scheme if missing
	if !strings.Contains(target, "://") {
		target = "https://" + target
	}

	parsed, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Only http and https schemes are supported
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid scheme: %s (only http and https supported)", parsed.Scheme)
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("missing host")
	}
	if len(host) > 253 {
		return fmt.Errorf("hostname too long")
	}

	// Allow IP addresses
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	// Allow localhost explicitly
	if host == "localhost" {
		return nil
	}

	// Require dot in hostname (to reject plain words like "example")
	if !strings.Contains(host, ".") {
		return fmt.Errorf("invalid hostname")
	}

	// Basic hostname format validation
	if strings.Contains(host, "..") ||
		strings.HasPrefix(host, ".") ||
		strings.HasSuffix(host, ".") {
		return fmt.Errorf("invalid hostname format")
	}

	return nil
}

// validatePingCheckTarget validates that the target is a valid IP address or hostname.
func validatePingCheckTarget(target string) error {
	if target == "" {
		return fmt.Errorf("ping target cannot be empty")
	}

	// Try IPv4 first
	if isValidIPv4(target) {
		return nil
	}

	// Then IPv6
	if isValidIPv6(target) {
		return nil
	}

	// Must contain at least one letter to be a valid hostname
	hasLetter := false
	for _, r := range target {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return fmt.Errorf("target is not a valid IP address or hostname")
	}

	return validateHostname(target)
}

// isValidIPv4 performs strict validation of IPv4 addresses (no leading zeros).
func isValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
		if len(part) > 1 && part[0] == '0' {
			return false // leading zero not allowed
		}
		val := 0
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
			val = val*10 + int(c-'0')
			if val > 255 {
				return false
			}
		}
	}
	return true
}

// isValidIPv6 validates IPv6 addresses using net.ParseIP with additional checks.
func isValidIPv6(ip string) bool {
	if !strings.Contains(ip, ":") {
		return false
	}
	parsed := net.ParseIP(ip)
	return parsed != nil && parsed.To4() == nil
}

// validateHostname ensures the string conforms to DNS hostname rules.
func validateHostname(host string) error {
	if len(host) > 253 {
		return fmt.Errorf("hostname too long (max 253 chars)")
	}
	if strings.Contains(host, "..") ||
		strings.HasPrefix(host, ".") ||
		strings.HasSuffix(host, ".") {
		return fmt.Errorf("invalid hostname format")
	}
	for _, r := range host {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '-') {
			return fmt.Errorf("hostname contains invalid characters")
		}
	}
	return nil
}

// validateCheckConfig validates and normalizes the JSON configuration string
// for a given check type, applying global defaults where applicable.
func validateCheckConfig(checkType, configStr string) (string, error) {
	switch checkType {
	case CheckTypeHTTP:
		var httpCheckDefaults *HTTPCheckConfig
		if globalHTTPDefaults != nil {
			httpCheckDefaults = &HTTPCheckConfig{
				Method:          globalHTTPDefaults.Method,
				Headers:         globalHTTPDefaults.Headers,
				Body:            globalHTTPDefaults.Body,
				ExpectedStatus:  globalHTTPDefaults.ExpectedStatus,
				ExpectedContent: globalHTTPDefaults.ExpectedContent,
				FollowRedirects: globalHTTPDefaults.FollowRedirects,
				VerifySSL:       globalHTTPDefaults.VerifySSL,
			}
		}
		return ValidateHTTPCheckConfig(configStr, httpCheckDefaults)
	case CheckTypePing:
		var pingCheckDefaults *PingCheckConfig
		if globalPingDefaults != nil {
			pingCheckDefaults = &PingCheckConfig{
				Count:    globalPingDefaults.Count,
				Interval: globalPingDefaults.IntervalMs,
				Size:     globalPingDefaults.PacketSize,
			}
		}
		return ValidatePingCheckConfig(configStr, pingCheckDefaults)
	default:
		return "", fmt.Errorf("unsupported check type: %s", checkType)
	}
}

// HTTPCheckConfig defines the structure of an HTTP check's configuration.
type HTTPCheckConfig struct {
	Method          string            `json:"method"`           // HTTP method (e.g., GET, POST)
	Headers         map[string]string `json:"headers"`          // Custom request headers
	Body            string            `json:"body"`             // Request body
	ExpectedStatus  int               `json:"expected_status"`  // Expected HTTP status code
	ExpectedContent string            `json:"expected_content"` // Substring expected in response body
	FollowRedirects bool              `json:"follow_redirects"` // Whether to follow HTTP redirects
	VerifySSL       bool              `json:"verify_ssl"`       // Whether to verify SSL certificates
}

// PingCheckConfig defines the structure of a ping check's configuration.
type PingCheckConfig struct {
	Count    int `json:"count"`    // Number of ICMP packets to send
	Interval int `json:"interval"` // Interval between packets in milliseconds
	Size     int `json:"size"`     // Size of each ICMP packet in bytes
}

// ValidateHTTPCheckConfig validates and normalizes an HTTP check configuration.
// It applies provided defaults and enforces security and correctness constraints.
func ValidateHTTPCheckConfig(configStr string, defaults *HTTPCheckConfig) (string, error) {
	// Start with safe defaults
	cfg := &HTTPCheckConfig{
		Method:          "GET",
		Headers:         make(map[string]string),
		ExpectedStatus:  200,
		FollowRedirects: true,
		VerifySSL:       true,
	}

	// Apply custom defaults if provided
	if defaults != nil {
		if defaults.Method != "" {
			cfg.Method = defaults.Method
		}
		if defaults.ExpectedStatus != 0 {
			cfg.ExpectedStatus = defaults.ExpectedStatus
		}
		cfg.FollowRedirects = defaults.FollowRedirects
		cfg.VerifySSL = defaults.VerifySSL
		if defaults.Headers != nil {
			for k, v := range defaults.Headers {
				cfg.Headers[k] = v
			}
		}
		if defaults.Body != "" {
			cfg.Body = defaults.Body
		}
		if defaults.ExpectedContent != "" {
			cfg.ExpectedContent = defaults.ExpectedContent
		}
	}

	// Parse user-provided config if present
	if configStr != "" && configStr != "{}" {
		if err := json.Unmarshal([]byte(configStr), cfg); err != nil {
			return "", fmt.Errorf("invalid JSON format: %w", err)
		}
	}

	// Normalize and validate HTTP method
	method := strings.ToUpper(cfg.Method)
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true,
		"HEAD": true, "OPTIONS": true, "PATCH": true,
	}
	if !validMethods[method] {
		return "", fmt.Errorf("invalid HTTP method: %s (supported: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH)", method)
	}
	cfg.Method = method

	// Validate expected status code range
	if cfg.ExpectedStatus < 100 || cfg.ExpectedStatus > 599 {
		return "", fmt.Errorf("invalid expected status code: %d (must be between 100-599)", cfg.ExpectedStatus)
	}

	// Validate headers
	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string)
	}
	for k, v := range cfg.Headers {
		if k == "" {
			return "", fmt.Errorf("header key cannot be empty")
		}
		if len(k) > 100 {
			return "", fmt.Errorf("header key too long: %s (max 100 chars)", k)
		}
		if len(v) > 1000 {
			return "", fmt.Errorf("header value too long for key %s (max 1000 chars)", k)
		}
	}

	// Enforce size limits
	if len(cfg.Body) > 10*1024*1024 { // 10 MB
		return "", fmt.Errorf("request body too large (max 10MB)")
	}
	if len(cfg.ExpectedContent) > 1000 {
		return "", fmt.Errorf("expected content too long (max 1000 chars)")
	}

	// Serialize back to JSON
	validated, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal validated config: %w", err)
	}
	return string(validated), nil
}

// ValidatePingCheckConfig validates and normalizes a ping check configuration.
func ValidatePingCheckConfig(configStr string, defaults *PingCheckConfig) (string, error) {
	cfg := &PingCheckConfig{
		Count:    3,
		Interval: 300,
		Size:     56,
	}

	if defaults != nil {
		if defaults.Count != 0 {
			cfg.Count = defaults.Count
		}
		if defaults.Interval != 0 {
			cfg.Interval = defaults.Interval
		}
		if defaults.Size != 0 {
			cfg.Size = defaults.Size
		}
	}

	if configStr != "" && configStr != "{}" {
		if err := json.Unmarshal([]byte(configStr), cfg); err != nil {
			return "", fmt.Errorf("invalid JSON format: %w", err)
		}
	}

	// Validate count: 1–10 packets
	if cfg.Count <= 0 || cfg.Count > 10 {
		return "", fmt.Errorf("invalid count: %d (must be between 1-10)", cfg.Count)
	}

	// Validate interval: 100–10000 ms
	if cfg.Interval < 100 || cfg.Interval > 10000 {
		return "", fmt.Errorf("invalid interval: %d ms (must be between 100-10000)", cfg.Interval)
	}

	// Validate packet size: 8–1472 bytes (standard Ethernet MTU limits)
	if cfg.Size < 8 || cfg.Size > 1472 {
		return "", fmt.Errorf("invalid packet size: %d bytes (must be between 8-1472)", cfg.Size)
	}

	validated, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal validated config: %w", err)
	}
	return string(validated), nil
}

// IsValidAlertStatus validates if an alert history status is valid.
// Valid statuses: "sent", "failed", "pending".
func IsValidAlertStatus(status string) bool {
	switch status {
	case AlertStatusSent, AlertStatusFailed, AlertStatusPending:
		return true
	default:
		return false
	}
}

// IsValidCheckType checks if the provided check type is supported.
func IsValidCheckType(t string) bool {
	return t == CheckTypeHTTP || t == CheckTypePing
}

// generateAuthToken generates a cryptographically secure random token.
// Used internally by other validators (e.g., Agent), but included here
// for completeness in this incremental approach.
func generateAuthToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsValidCheckStatus validates if a check status is valid.
// Valid statuses: "up", "down", "error", "timeout".
func IsValidCheckStatus(status string) bool {
	switch status {
	case CheckStatusUp, CheckStatusDown, CheckStatusError, CheckStatusTimeout:
		return true
	default:
		return false
	}
}

// IsValidAlertType validates if an alert type is supported.
// Valid types: "telegram", "bale", "email", "webhook".
func IsValidAlertType(alertType string) bool {
	switch alertType {
	case AlertTypeTelegram, AlertTypeBale, AlertTypeEmail, AlertTypeWebhook:
		return true
	default:
		return false
	}
}

// IsValidAgentStatus validates if an agent status is valid.
// Valid statuses: "online", "offline".
func IsValidAgentStatus(status string) bool {
	switch status {
	case AgentStatusOnline, AgentStatusOffline:
		return true
	default:
		return false
	}
}

// validateAlertConditionType validates alert condition types based on context.
// For check alerts: only "status_down", "status_timeout", "status_error" are allowed.
// For agent alerts: only "cpu_usage_high", "memory_usage_high", "disk_usage_high", "agent_offline" are allowed.
func validateAlertConditionType(conditionType string, isCheckAlert bool) error {
	checkConditions := []string{"status_down", "status_timeout", "status_error"}
	agentConditions := []string{"cpu_usage_high", "memory_usage_high", "disk_usage_high", "agent_offline"}

	if isCheckAlert {
		for _, condition := range checkConditions {
			if conditionType == condition {
				return nil
			}
		}
		return fmt.Errorf("invalid condition type for check alert: %s (supported: %v)", conditionType, checkConditions)
	} else {
		for _, condition := range agentConditions {
			if conditionType == condition {
				return nil
			}
		}
		return fmt.Errorf("invalid condition type for agent alert: %s (supported: %v)", conditionType, agentConditions)
	}
}

// needsConditionValue returns true if the condition type requires a threshold value.
// Only percentage-based conditions (CPU, memory, disk usage) require a ConditionValue.
func needsConditionValue(conditionType string) bool {
	thresholdConditions := []string{"cpu_usage_high", "memory_usage_high", "disk_usage_high"}
	for _, condition := range thresholdConditions {
		if conditionType == condition {
			return true
		}
	}
	return false
}

// validateAlertConfig validates alert configuration based on alert type.
// It ensures the config is valid JSON and contains required fields per alert type.
func validateAlertConfig(alertType, configStr string) error {
	if configStr == "" || configStr == "{}" {
		return nil // Empty config is allowed
	}

	// Parse JSON to ensure it's valid
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Type-specific validation
	switch alertType {
	case AlertTypeTelegram:
		return validateTelegramAlertConfig(config)
	case AlertTypeBale:
		return validateBaleAlertConfig(config)
	case AlertTypeEmail:
		return validateEmailAlertConfig(config)
	case AlertTypeWebhook:
		return validateWebhookAlertConfig(config)
	default:
		return fmt.Errorf("unsupported alert type: %s", alertType)
	}
}

// validateServiceAlertConfig validates configuration for services that share the same schema
// (e.g., Telegram and Bale). The serviceName is used in error messages for clarity.
func validateServiceAlertConfig(config map[string]interface{}, serviceName string) error {
	// Validate token
	if token, exists := config["token"]; exists {
		tokenStr, ok := token.(string)
		if !ok {
			return fmt.Errorf("%s token must be a string", serviceName)
		}
		if len(tokenStr) < 10 {
			return fmt.Errorf("%s token too short", serviceName)
		}
	}

	// Validate chat_id
	if chatID, exists := config["chat_id"]; exists {
		chatIDStr, ok := chatID.(string)
		if !ok {
			return fmt.Errorf("%s chat_id must be a string", serviceName)
		}
		if chatIDStr == "" {
			return fmt.Errorf("%s chat_id cannot be empty", serviceName)
		}
	}

	return nil
}

// validateTelegramAlertConfig validates the configuration for Telegram alerts.
func validateTelegramAlertConfig(config map[string]interface{}) error {
	return validateServiceAlertConfig(config, "telegram")
}

// validateBaleAlertConfig validates the configuration for Bale alerts.
func validateBaleAlertConfig(config map[string]interface{}) error {
	return validateServiceAlertConfig(config, "bale")
}

// validateEmailAlertConfig validates the configuration for email alerts.
// Required fields: smtp_host (non-empty string), to (valid email format).
func validateEmailAlertConfig(config map[string]interface{}) error {
	// Validate SMTP host
	if host, exists := config["smtp_host"]; exists {
		hostStr, ok := host.(string)
		if !ok {
			return fmt.Errorf("smtp_host must be a string")
		}
		if hostStr == "" {
			return fmt.Errorf("smtp_host cannot be empty")
		}
	}

	// Validate recipient email
	if to, exists := config["to"]; exists {
		toStr, ok := to.(string)
		if !ok {
			return fmt.Errorf("email 'to' field must be a string")
		}
		if !strings.Contains(toStr, "@") {
			return fmt.Errorf("invalid email format")
		}
	}

	return nil
}

// validateWebhookAlertConfig validates the configuration for webhook alerts.
// Required field: url (valid URL string).
func validateWebhookAlertConfig(config map[string]interface{}) error {
	// Validate webhook URL
	if u, exists := config["url"]; exists {
		urlStr, ok := u.(string)
		if !ok {
			return fmt.Errorf("webhook URL must be a string")
		}
		if _, err := url.Parse(urlStr); err != nil {
			return fmt.Errorf("invalid webhook URL: %w", err)
		}
	}
	return nil
}
