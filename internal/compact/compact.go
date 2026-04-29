// Package compact removes duplicate values from a snapshot, keeping only
// the last occurrence of each value across all keys.
package compact

import (
	"fmt"
	"sort"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a compact operation.
type Result struct {
	Removed []string // keys whose values were duplicates
	Kept    int      // number of unique keys retained
}

// Run compacts the snapshot identified by name: for every value that appears
// more than once, only the lexicographically last key is kept.
func Run(m Manager, name string) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("compact: load %q: %w", name, err)
	}

	// Build reverse map: value -> list of keys
	reverse := make(map[string][]string)
	for k, v := range vars {
		reverse[v] = append(reverse[v], k)
	}

	removed := []string{}
	compacted := make(map[string]string, len(vars))

	for v, keys := range reverse {
		if len(keys) == 1 {
			compacted[keys[0]] = v
			continue
		}
		// Keep the last key alphabetically; drop the rest.
		sort.Strings(keys)
		keeper := keys[len(keys)-1]
		compacted[keeper] = v
		removed = append(removed, keys[:len(keys)-1]...)
	}

	sort.Strings(removed)

	if err := m.Save(name, compacted); err != nil {
		return Result{}, fmt.Errorf("compact: save %q: %w", name, err)
	}

	return Result{Removed: removed, Kept: len(compacted)}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Removed) == 0 {
		return fmt.Sprintf("compact: no duplicate values found (%d keys kept)", r.Kept)
	}
	out := fmt.Sprintf("compact: removed %d duplicate key(s), %d kept\n", len(r.Removed), r.Kept)
	for _, k := range r.Removed {
		out += fmt.Sprintf("  - %s\n", k)
	}
	return out
}
