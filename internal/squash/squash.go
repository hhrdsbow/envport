// Package squash merges multiple snapshots into a single snapshot,
// keeping the last-write-wins value for each key.
package squash

import "fmt"

// Snapshot is the minimal interface squash needs from a profile manager.
type Snapshot interface {
	Vars() map[string]string
	Name() string
}

// Manager can load and save snapshots.
type Manager interface {
	Load(name string) (Snapshot, error)
	Save(name string, vars map[string]string) error
}

// Result describes what squash produced.
type Result struct {
	Dest   string
	Sources []string
	KeyCount int
}

// Run loads each source snapshot in order, merges their variables
// (later sources win on conflict), and saves the result as dest.
func Run(m Manager, sources []string, dest string) (Result, error) {
	if len(sources) < 2 {
		return Result{}, fmt.Errorf("squash requires at least two source snapshots")
	}

	merged := make(map[string]string)

	for _, name := range sources {
		snap, err := m.Load(name)
		if err != nil {
			return Result{}, fmt.Errorf("load %q: %w", name, err)
		}
		for k, v := range snap.Vars() {
			merged[k] = v
		}
	}

	if err := m.Save(dest, merged); err != nil {
		return Result{}, fmt.Errorf("save %q: %w", dest, err)
	}

	return Result{
		Dest:    dest,
		Sources: sources,
		KeyCount: len(merged),
	}, nil
}

// Format returns a human-readable summary of a squash Result.
func Format(r Result) string {
	return fmt.Sprintf("squashed %d source(s) into %q (%d keys)",
		len(r.Sources), r.Dest, r.KeyCount)
}
