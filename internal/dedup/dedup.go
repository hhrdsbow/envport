// Package dedup removes duplicate keys across multiple snapshots,
// keeping the value from the highest-priority (first) snapshot.
package dedup

import (
	"fmt"
	"sort"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a deduplication run.
type Result struct {
	Kept    map[string]string // final merged map
	Dropped []string          // keys that appeared in lower-priority snapshots and were dropped
}

// Run merges the named snapshots in priority order (first = highest priority).
// Duplicate keys from lower-priority snapshots are dropped.
// When dst is non-empty the merged result is saved under that name.
func Run(m Manager, names []string, dst string) (Result, error) {
	if len(names) < 2 {
		return Result{}, fmt.Errorf("dedup requires at least two snapshot names")
	}

	seen := make(map[string]bool)
	merged := make(map[string]string)
	droppedSet := make(map[string]bool)

	for _, name := range names {
		vars, err := m.Load(name)
		if err != nil {
			return Result{}, fmt.Errorf("load %q: %w", name, err)
		}
		for k, v := range vars {
			if seen[k] {
				droppedSet[k] = true
				continue
			}
			seen[k] = true
			merged[k] = v
		}
	}

	dropped := make([]string, 0, len(droppedSet))
	for k := range droppedSet {
		dropped = append(dropped, k)
	}
	sort.Strings(dropped)

	if dst != "" {
		if err := m.Save(dst, merged); err != nil {
			return Result{}, fmt.Errorf("save %q: %w", dst, err)
		}
	}

	return Result{Kept: merged, Dropped: dropped}, nil
}

// Format returns a human-readable summary of the dedup result.
func Format(r Result) string {
	if len(r.Dropped) == 0 {
		return fmt.Sprintf("dedup complete: %d keys kept, no duplicates found", len(r.Kept))
	}
	return fmt.Sprintf("dedup complete: %d keys kept, %d duplicate keys dropped: %v",
		len(r.Kept), len(r.Dropped), r.Dropped)
}
