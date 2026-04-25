// Package pivot provides functionality to transpose a snapshot's
// key-value pairs into a new snapshot keyed by value.
package pivot

import "fmt"

// Snapshot is the minimal interface required by the pivot runner.
type Snapshot interface {
	Vars() map[string]string
}

// Manager can load and save snapshots by name.
type Manager interface {
	Load(name string) (Snapshot, error)
	Save(name string, vars map[string]string) error
}

// Result describes the outcome of a pivot operation.
type Result struct {
	Src  string
	Dest string
	Keys int
}

// Run loads the snapshot named src, inverts its key→value mapping so that
// values become keys and keys become values, then saves the result as dest.
// If two keys share the same value the last one (in iteration order) wins;
// callers that need determinism should sort before calling.
func Run(m Manager, src, dest string) (Result, error) {
	snap, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("pivot: load %q: %w", src, err)
	}

	vars := snap.Vars()
	inverted := make(map[string]string, len(vars))
	for k, v := range vars {
		if v == "" {
			continue // skip blank values — they cannot become keys
		}
		inverted[v] = k
	}

	if err := m.Save(dest, inverted); err != nil {
		return Result{}, fmt.Errorf("pivot: save %q: %w", dest, err)
	}

	return Result{Src: src, Dest: dest, Keys: len(inverted)}, nil
}
