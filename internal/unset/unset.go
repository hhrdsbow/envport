// Package unset removes specified keys from a snapshot and saves the result.
package unset

import (
	"errors"
	"fmt"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of an unset operation.
type Result struct {
	Removed []string
	NotFound []string
}

// Run removes keys from the named snapshot. If strict is true, referencing a
// key that does not exist in the snapshot is treated as an error.
func Run(m Manager, name string, keys []string, strict bool) (*Result, error) {
	if name == "" {
		return nil, errors.New("snapshot name must not be empty")
	}
	if len(keys) == 0 {
		return nil, errors.New("at least one key must be specified")
	}

	vars, err := m.Load(name)
	if err != nil {
		return nil, fmt.Errorf("load snapshot %q: %w", name, err)
	}

	res := &Result{}
	for _, k := range keys {
		if _, ok := vars[k]; ok {
			delete(vars, k)
			res.Removed = append(res.Removed, k)
		} else {
			if strict {
				return nil, fmt.Errorf("key %q not found in snapshot %q", k, name)
			}
			res.NotFound = append(res.NotFound, k)
		}
	}

	if err := m.Save(name, vars); err != nil {
		return nil, fmt.Errorf("save snapshot %q: %w", name, err)
	}
	return res, nil
}

// Format returns a human-readable summary of the result.
func Format(r *Result) string {
	if len(r.Removed) == 0 {
		return "no keys removed"
	}
	out := fmt.Sprintf("removed %d key(s):", len(r.Removed))
	for _, k := range r.Removed {
		out += "\n  - " + k
	}
	if len(r.NotFound) > 0 {
		out += fmt.Sprintf("\n%d key(s) not found (skipped):", len(r.NotFound))
		for _, k := range r.NotFound {
			out += "\n  ? " + k
		}
	}
	return out
}
