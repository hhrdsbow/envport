// Package trim provides functionality to remove specific keys from a snapshot.
package trim

import (
	"fmt"
)

// Manager describes the storage operations required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a trim operation.
type Result struct {
	Removed []string
	Kept    int
}

// Run removes the given keys from the named snapshot and saves it back.
// If dryRun is true the snapshot is not modified.
func Run(m Manager, name string, keys []string, dryRun bool) (Result, error) {
	if len(keys) == 0 {
		return Result{}, fmt.Errorf("trim: no keys specified")
	}

	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("trim: load %q: %w", name, err)
	}

	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	var removed []string
	for _, k := range keys {
		if _, exists := vars[k]; exists {
			removed = append(removed, k)
			if !dryRun {
				delete(vars, k)
			}
		}
	}

	if !dryRun && len(removed) > 0 {
		if err := m.Save(name, vars); err != nil {
			return Result{}, fmt.Errorf("trim: save %q: %w", name, err)
		}
	}

	return Result{Removed: removed, Kept: len(vars)}, nil
}
