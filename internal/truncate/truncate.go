// Package truncate provides functionality to shorten snapshot variable values
// to a maximum byte length, useful for display or storage constraints.
package truncate

import "fmt"

// Manager describes the interface required to load and save snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Options controls how truncation is applied.
type Options struct {
	// MaxLen is the maximum number of characters allowed per value.
	// Values longer than MaxLen are cut and suffixed with Ellipsis.
	MaxLen int
	// Ellipsis is appended to truncated values. Defaults to "...".
	Ellipsis string
	// Keys restricts truncation to the listed keys. Empty means all keys.
	Keys []string
}

// Result holds per-key information about what was truncated.
type Result struct {
	Key      string
	Original string
	Truncated string
}

// Run loads the named snapshot, truncates values according to opts, saves the
// result back, and returns a slice describing every change made.
func Run(m Manager, name string, opts Options) ([]Result, error) {
	if opts.MaxLen <= 0 {
		return nil, fmt.Errorf("truncate: MaxLen must be greater than zero")
	}
	ellipsis := opts.Ellipsis
	if ellipsis == "" {
		ellipsis = "..."
	}

	vars, err := m.Load(name)
	if err != nil {
		return nil, fmt.Errorf("truncate: load %q: %w", name, err)
	}

	filter := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		filter[k] = true
	}

	var results []Result
	for k, v := range vars {
		if len(filter) > 0 && !filter[k] {
			continue
		}
		if len(v) > opts.MaxLen {
			short := v[:opts.MaxLen] + ellipsis
			results = append(results, Result{Key: k, Original: v, Truncated: short})
			vars[k] = short
		}
	}

	if len(results) > 0 {
		if err := m.Save(name, vars); err != nil {
			return nil, fmt.Errorf("truncate: save %q: %w", name, err)
		}
	}
	return results, nil
}
