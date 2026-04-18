// Package prune removes snapshots older than a given age or beyond a count limit.
package prune

import (
	"fmt"
	"time"
)

// Snapshot is the minimal interface prune needs.
type Snapshot interface {
	Name() string
	CreatedAt() time.Time
}

// Manager can list and delete snapshots.
type Manager interface {
	List() ([]Snapshot, error)
	Delete(name string) error
}

// Options controls pruning behaviour.
type Options struct {
	// OlderThan removes snapshots created before this time. Zero means no age limit.
	OlderThan time.Time
	// KeepLast keeps at most this many snapshots (newest first). 0 means no limit.
	KeepLast int
	// DryRun lists what would be deleted without deleting.
	DryRun bool
}

// Result holds the names of pruned (or would-be-pruned) snapshots.
type Result struct {
	Pruned []string
}

// Run executes the prune operation and returns the result.
func Run(m Manager, opts Options) (Result, error) {
	snaps, err := m.List()
	if err != nil {
		return Result{}, fmt.Errorf("prune: list: %w", err)
	}

	// Sort newest-first (assume List returns arbitrary order).
	sorted := make([]Snapshot, len(snaps))
	copy(sorted, snaps)
	sortByAge(sorted)

	toDelete := map[string]bool{}

	if !opts.OlderThan.IsZero() {
		for _, s := range sorted {
			if s.CreatedAt().Before(opts.OlderThan) {
				toDelete[s.Name()] = true
			}
		}
	}

	if opts.KeepLast > 0 {
		for i, s := range sorted {
			if i >= opts.KeepLast {
				toDelete[s.Name()] = true
			}
		}
	}

	var result Result
	for _, s := range sorted {
		if !toDelete[s.Name()] {
			continue
		}
		result.Pruned = append(result.Pruned, s.Name())
		if !opts.DryRun {
			if err := m.Delete(s.Name()); err != nil {
				return result, fmt.Errorf("prune: delete %q: %w", s.Name(), err)
			}
		}
	}
	return result, nil
}

func sortByAge(snaps []Snapshot) {
	// Simple insertion sort (lists are small).
	for i := 1; i < len(snaps); i++ {
		for j := i; j > 0 && snaps[j].CreatedAt().After(snaps[j-1].CreatedAt()); j-- {
			snaps[j], snaps[j-1] = snaps[j-1], snaps[j]
		}
	}
}
