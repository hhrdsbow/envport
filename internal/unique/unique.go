// Package unique removes duplicate values (not keys) from a snapshot,
// keeping only the first key encountered for each distinct value.
package unique

import (
	"fmt"
	"sort"
)

// Manager can load and save snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a unique operation.
type Result struct {
	Kept    []string
	Dropped []string
}

// Run loads the snapshot identified by name, removes keys whose values are
// duplicates of an already-seen value, and saves the result back.
// Keys are processed in sorted order so the behaviour is deterministic.
func Run(m Manager, name string) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("unique: load %q: %w", name, err)
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	seen := make(map[string]bool)
	result := make(map[string]string)
	var kept, dropped []string

	for _, k := range keys {
		v := vars[k]
		if seen[v] {
			dropped = append(dropped, k)
			continue
		}
		seen[v] = true
		result[k] = v
		kept = append(kept, k)
	}

	if err := m.Save(name, result); err != nil {
		return Result{}, fmt.Errorf("unique: save %q: %w", name, err)
	}

	return Result{Kept: kept, Dropped: dropped}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Dropped) == 0 {
		return "no duplicate values found"
	}
	msg := fmt.Sprintf("removed %d key(s) with duplicate values:", len(r.Dropped))
	for _, k := range r.Dropped {
		msg += "\n  - " + k
	}
	return msg
}
