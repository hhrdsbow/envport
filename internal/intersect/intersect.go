// Package intersect computes the common keys and their values across
// two or more named snapshots.
package intersect

import "fmt"

// Manager is the interface required to load snapshots by name.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Result holds the keys present in every snapshot along with each
// snapshot's value for that key.
type Result struct {
	// Keys lists the common keys in sorted order.
	Keys []string
	// Values maps key -> (snapshotName -> value).
	Values map[string]map[string]string
}

// Run loads each named snapshot and returns the keys that appear in
// all of them together with their per-snapshot values.
func Run(m Manager, names []string) (*Result, error) {
	if len(names) < 2 {
		return nil, fmt.Errorf("intersect: at least two snapshot names are required")
	}

	snaps := make([]map[string]string, 0, len(names))
	for _, n := range names {
		vars, err := m.Load(n)
		if err != nil {
			return nil, fmt.Errorf("intersect: load %q: %w", n, err)
		}
		snaps = append(snaps, vars)
	}

	// Seed with keys from the first snapshot, then narrow down.
	candidates := make(map[string]struct{})
	for k := range snaps[0] {
		candidates[k] = struct{}{}
	}
	for _, s := range snaps[1:] {
		for k := range candidates {
			if _, ok := s[k]; !ok {
				delete(candidates, k)
			}
		}
	}

	// Build sorted key list.
	keys := make([]string, 0, len(candidates))
	for k := range candidates {
		keys = append(keys, k)
	}
	sortStrings(keys)

	// Populate per-snapshot values.
	values := make(map[string]map[string]string, len(keys))
	for _, k := range keys {
		perSnap := make(map[string]string, len(names))
		for i, n := range names {
			perSnap[n] = snaps[i][k]
		}
		values[k] = perSnap
	}

	return &Result{Keys: keys, Values: values}, nil
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
