// Package suffix provides functionality to add or remove a suffix
// from environment variable values within a snapshot.
package suffix

import (
	"fmt"
	"strings"
)

// Manager is the interface required to load and save snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a suffix operation.
type Result struct {
	Modified []string
	Skipped  []string
}

// RunAdd appends suffix to the values of the specified keys (or all keys if
// keys is empty) in the named snapshot and saves it back.
func RunAdd(m Manager, name, sfx string, keys []string) (Result, error) {
	if sfx == "" {
		return Result{}, fmt.Errorf("suffix must not be empty")
	}
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, err
	}
	return applyOp(m, name, vars, sfx, keys, func(v, s string) string {
		return v + s
	})
}

// RunRemove strips suffix from the values of the specified keys (or all keys
// if keys is empty) in the named snapshot and saves it back.
func RunRemove(m Manager, name, sfx string, keys []string) (Result, error) {
	if sfx == "" {
		return Result{}, fmt.Errorf("suffix must not be empty")
	}
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, err
	}
	return applyOp(m, name, vars, sfx, keys, func(v, s string) string {
		return strings.TrimSuffix(v, s)
	})
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "modified: %d, skipped: %d\n", len(r.Modified), len(r.Skipped))
	for _, k := range r.Modified {
		fmt.Fprintf(&sb, "  ~ %s\n", k)
	}
	return sb.String()
}

func applyOp(m Manager, name string, vars map[string]string, sfx string, keys []string, fn func(string, string) string) (Result, error) {
	target := toSet(keys)
	var res Result
	for k, v := range vars {
		if len(target) > 0 && !target[k] {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		newVal := fn(v, sfx)
		if newVal != v {
			vars[k] = newVal
			res.Modified = append(res.Modified, k)
		} else {
			res.Skipped = append(res.Skipped, k)
		}
	}
	if err := m.Save(name, vars); err != nil {
		return Result{}, err
	}
	return res, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
