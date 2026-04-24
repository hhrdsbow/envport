// Package redact provides utilities for masking sensitive environment
// variable values before display or export.
package redact

import (
	"strings"
)

// DefaultPatterns is the list of key substrings that trigger redaction.
var DefaultPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const mask = "***"

// IsSensitive reports whether the given key name matches any of the
// provided patterns (case-insensitive substring match).
func IsSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Apply returns a copy of vars where values whose keys are considered
// sensitive are replaced with the mask string.
func Apply(vars map[string]string, patterns []string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if IsSensitive(k, patterns) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

// ApplyDefault is a convenience wrapper that uses DefaultPatterns.
func ApplyDefault(vars map[string]string) map[string]string {
	return Apply(vars, DefaultPatterns)
}
