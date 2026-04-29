// Package omit removes a specified set of keys from a snapshot.
package omit

import (
	"fmt"
	"sort"
)

// Manager describes the profile store operations required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of an omit operation.
type Result struct {
	Removed []string
	Remaining int
}

// Run removes the given keys from the named snapshot and saves it back.
// Keys that do not exist in the snapshot are silently ignored.
func Run(m Manager, name string, keys []string) (Result, error) {
	if name == "" {
		return Result{}, fmt.Errorf("snapshot name must not be empty")
	}
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("load %q: %w", name, err)
	}

	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	var removed []string
	for k := range keySet {
		if _, ok := vars[k]; ok {
			delete(vars, k)
			removed = append(removed, k)
		}
	}
	sort.Strings(removed)

	if err := m.Save(name, vars); err != nil {
		return Result{}, fmt.Errorf("save %q: %w", name, err)
	}
	return Result{Removed: removed, Remaining: len(vars)}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Removed) == 0 {
		return fmt.Sprintf("no keys removed; %d key(s) remaining", r.Remaining)
	}
	out := fmt.Sprintf("removed %d key(s): ", len(r.Removed))
	for i, k := range r.Removed {
		if i > 0 {
			out += ", "
		}
		out += k
	}
	out += fmt.Sprintf("\n%d key(s) remaining", r.Remaining)
	return out
}
