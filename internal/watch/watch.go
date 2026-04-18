// Package watch detects changes to the current environment compared to a saved snapshot.
package watch

import (
	"fmt"

	"envport/internal/diff"
	"envport/internal/env"
	"envport/internal/snapshot"
)

// Manager is the interface required to load snapshots.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
}

// Result holds the outcome of a watch check.
type Result struct {
	Name    string
	Changed bool
	Report  string
}

// Check compares the current environment against the named snapshot.
// It returns a Result describing whether any variables have drifted.
func Check(m Manager, name string, keys []string) (*Result, error) {
	snap, err := m.Load(name)
	if err != nil {
		return nil, fmt.Errorf("watch: load snapshot %q: %w", name, err)
	}

	current := env.Capture()
	if len(keys) > 0 {
		current = env.FilterKeys(current, keys)
	}

	base := snap.Vars
	if len(keys) > 0 {
		base = env.FilterKeys(base, keys)
	}

	changes := diff.Compute(base, current)
	report := diff.Format(changes)

	return &Result{
		Name:    name,
		Changed: len(changes) > 0,
		Report:  report,
	}, nil
}
