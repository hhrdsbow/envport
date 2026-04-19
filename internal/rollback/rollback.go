// Package rollback restores a previous snapshot from history.
package rollback

import (
	"fmt"

	"envport/internal/history"
	"envport/internal/snapshot"
)

// Manager defines the profile operations needed for rollback.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
	Save(name string, snap *snapshot.Snapshot) error
}

// HistoryReader lists history entries for a profile.
type HistoryReader interface {
	List(profile string) ([]history.Entry, error)
}

// Result describes what rollback did.
type Result struct {
	Profile  string
	RolledTo string
}

// Run rolls back the named profile to the nth previous snapshot recorded in history.
// offset 1 means the most recent history entry before the current state.
func Run(mgr Manager, hr HistoryReader, profile string, offset int) (*Result, error) {
	if offset < 1 {
		return nil, fmt.Errorf("offset must be >= 1")
	}

	entries, err := hr.List(profile)
	if err != nil {
		return nil, fmt.Errorf("listing history: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no history found for profile %q", profile)
	}
	if offset > len(entries) {
		return nil, fmt.Errorf("offset %d exceeds history length %d", offset, len(entries))
	}

	// entries are newest-first
	target := entries[offset-1]

	snap, err := mgr.Load(target.SnapshotRef)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot %q: %w", target.SnapshotRef, err)
	}

	if err := mgr.Save(profile, snap); err != nil {
		return nil, fmt.Errorf("saving rollback: %w", err)
	}

	return &Result{Profile: profile, RolledTo: target.SnapshotRef}, nil
}
