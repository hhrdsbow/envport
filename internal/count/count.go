// Package count provides functionality to count environment variables
// in a snapshot, optionally filtered by key pattern.
package count

import (
	"fmt"
	"strings"
)

// Manager is the interface required to load snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Result holds the output of a Count operation.
type Result struct {
	Name    string
	Total   int
	Matched int
	Pattern string
}

// Run counts the environment variables in the named snapshot.
// If pattern is non-empty, only keys containing the pattern (case-insensitive)
// are counted and reported as Matched.
func Run(m Manager, name, pattern string) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("count: load %q: %w", name, err)
	}

	res := Result{
		Name:    name,
		Total:   len(vars),
		Pattern: pattern,
	}

	if pattern == "" {
		res.Matched = len(vars)
		return res, nil
	}

	lower := strings.ToLower(pattern)
	for k := range vars {
		if strings.Contains(strings.ToLower(k), lower) {
			res.Matched++
		}
	}
	return res, nil
}

// Format returns a human-readable summary of the Result.
func Format(r Result) string {
	if r.Pattern == "" {
		return fmt.Sprintf("%s: %d variable(s)", r.Name, r.Total)
	}
	return fmt.Sprintf("%s: %d/%d variable(s) match %q", r.Name, r.Matched, r.Total, r.Pattern)
}
