// Package flatten merges multiple snapshots into a single flat variable map,
// applying each snapshot in order so later entries override earlier ones.
package flatten

import (
	"errors"
	"fmt"
)

// Snapshot is the minimal interface required by the flatten package.
type Snapshot interface {
	Name() string
	Vars() map[string]string
}

// Manager can load snapshots by name.
type Manager interface {
	Load(name string) (Snapshot, error)
}

// Result holds the merged variable map and a record of which snapshot
// each key was last contributed by.
type Result struct {
	Vars    map[string]string
	Sources map[string]string // key -> snapshot name
}

// Run loads each named snapshot in order and merges their variables.
// Later snapshots override keys from earlier ones.
func Run(mgr Manager, names []string) (*Result, error) {
	if len(names) == 0 {
		return nil, errors.New("flatten: at least one snapshot name is required")
	}

	out := &Result{
		Vars:    make(map[string]string),
		Sources: make(map[string]string),
	}

	for _, name := range names {
		snap, err := mgr.Load(name)
		if err != nil {
			return nil, fmt.Errorf("flatten: loading %q: %w", name, err)
		}
		for k, v := range snap.Vars() {
			out.Vars[k] = v
			out.Sources[k] = snap.Name()
		}
	}

	return out, nil
}
