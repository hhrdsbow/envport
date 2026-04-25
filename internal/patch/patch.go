// Package patch applies partial updates to an existing snapshot,
// overwriting only the keys provided while leaving others intact.
package patch

import (
	"errors"
	"fmt"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds a summary of the patch operation.
type Result struct {
	Name    string
	Added   []string
	Updated []string
}

// Run merges updates into the snapshot identified by name.
// Keys present in updates but absent from the snapshot are added;
// keys present in both are overwritten. Existing keys not in updates
// are left unchanged.
func Run(m Manager, name string, updates map[string]string) (Result, error) {
	if name == "" {
		return Result{}, errors.New("patch: name must not be empty")
	}
	if len(updates) == 0 {
		return Result{}, errors.New("patch: updates must not be empty")
	}

	current, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("patch: load %q: %w", name, err)
	}

	result := Result{Name: name}
	merged := make(map[string]string, len(current))
	for k, v := range current {
		merged[k] = v
	}

	for k, v := range updates {
		if _, exists := merged[k]; exists {
			result.Updated = append(result.Updated, k)
		} else {
			result.Added = append(result.Added, k)
		}
		merged[k] = v
	}

	if err := m.Save(name, merged); err != nil {
		return Result{}, fmt.Errorf("patch: save %q: %w", name, err)
	}

	return result, nil
}

// Format returns a human-readable summary of the patch result.
func Format(r Result) string {
	return fmt.Sprintf("patched %q: %d added, %d updated",
		r.Name, len(r.Added), len(r.Updated))
}
