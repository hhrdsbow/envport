package validate

import (
	"fmt"
	"strings"
)

// Result holds the outcome of a validation check.
type Result struct {
	Key     string
	Message string
}

// Manager defines what a snapshot manager must expose for validation.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Run checks a named snapshot against a set of required keys.
// It returns missing keys and keys whose values are empty.
func Run(mgr Manager, name string, required []string) ([]Result, error) {
	vars, err := mgr.Load(name)
	if err != nil {
		return nil, fmt.Errorf("validate: load %q: %w", name, err)
	}

	var results []Result
	for _, key := range required {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		val, ok := vars[key]
		if !ok {
			results = append(results, Result{Key: key, Message: "missing"})
		} else if strings.TrimSpace(val) == "" {
			results = append(results, Result{Key: key, Message: "empty"})
		}
	}
	return results, nil
}

// Format renders validation results as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "all required keys present and non-empty"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "  [%s] %s\n", r.Message, r.Key)
	}
	return sb.String()
}
