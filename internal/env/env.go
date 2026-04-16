package env

import (
	"fmt"
	"os"
	"strings"
)

// Capture returns a map of all current environment variables.
func Capture() map[string]string {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		env[key] = value
	}
	return env
}

// Apply sets environment variables from the provided map.
func Apply(vars map[string]string) error {
	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("setting %q: %w", key, err)
		}
	}
	return nil
}

// Diff returns keys that differ between current environment and the given map.
// Returns added, removed, and changed key sets.
func Diff(current, target map[string]string) (added, removed, changed []string) {
	for k, tv := range target {
		cv, exists := current[k]
		if !exists {
			added = append(added, k)
		} else if cv != tv {
			changed = append(changed, k)
		}
	}
	for k := range current {
		if _, exists := target[k]; !exists {
			removed = append(removed, k)
		}
	}
	return
}

// FilterKeys returns a new map containing only the specified keys.
func FilterKeys(vars map[string]string, keys []string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := vars[k]; ok {
			result[k] = v
		}
	}
	return result
}
