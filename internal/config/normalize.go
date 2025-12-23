package config

import "strings"

// normalizeConfig normalizes configuration values.
func normalizeConfig(c *Config) {
	// Normalize log level to lowercase
	c.Log.Level = strings.ToLower(c.Log.Level)
}
