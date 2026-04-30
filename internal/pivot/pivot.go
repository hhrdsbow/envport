// Package pivot swaps the keys and values of a snapshot so that every
// value becomes a key and every key becomes its value.
// Keys that would collide after pivoting are deduplicated by keeping the
// lexicographically first original key.
package pivot

import (
	"errors"
	"fmt"
	"sort"
)

// Manager is the minimal interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Run loads src, inverts key↔value, and saves the result to dst.
// If dst is empty the result is saved back to src.
func Run(m Manager, src, dst string) (map[string]string, error) {
	if src == "" {
		return nil, errors.New("source name must not be empty")
	}

	vars, err := m.Load(src)
	if err != nil {
		return nil, fmt.Errorf("load %q: %w", src, err)
	}

	pivoted, err := invert(vars)
	if err != nil {
		return nil, err
	}

	target := src
	if dst != "" {
		target = dst
	}

	if err := m.Save(target, pivoted); err != nil {
		return nil, fmt.Errorf("save %q: %w", target, err)
	}

	return pivoted, nil
}

// invert swaps keys and values, resolving collisions by preferring the
// lexicographically smallest original key.
func invert(vars map[string]string) (map[string]string, error) {
	// Collect keys in sorted order so collision resolution is deterministic.
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make(map[string]string, len(vars))
	for _, k := range keys {
		v := vars[k]
		if v == "" {
			continue // skip blank values — they cannot become valid keys
		}
		if _, exists := out[v]; !exists {
			out[v] = k
		}
	}
	return out, nil
}
