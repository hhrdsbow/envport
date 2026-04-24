// Package quota enforces limits on the number of snapshots stored per profile.
// It allows setting a maximum count and provides a check that returns which
// snapshots should be evicted (oldest first) to stay within the limit.
package quota

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// ErrNoQuota is returned when no quota has been set for a profile.
var ErrNoQuota = errors.New("no quota set")

// SnapshotMeta holds the minimal metadata needed to evaluate quota.
type SnapshotMeta struct {
	Name      string
	CreatedAt time.Time
}

// Store is the persistence interface required by Manager.
type Store interface {
	GetQuota(profile string) (int, error)
	SetQuota(profile string, max int) error
	DeleteQuota(profile string) error
}

// Manager manages per-profile snapshot quotas.
type Manager struct {
	store Store
}

// New creates a new Manager backed by the provided Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Set configures a maximum snapshot count for the given profile.
// max must be >= 1.
func (m *Manager) Set(profile string, max int) error {
	if max < 1 {
		return fmt.Errorf("quota max must be at least 1, got %d", max)
	}
	return m.store.SetQuota(profile, max)
}

// Get returns the current quota for a profile.
// Returns ErrNoQuota if none has been set.
func (m *Manager) Get(profile string) (int, error) {
	return m.store.GetQuota(profile)
}

// Remove deletes the quota for a profile, disabling enforcement.
func (m *Manager) Remove(profile string) error {
	return m.store.DeleteQuota(profile)
}

// Evictions returns the snapshots that should be deleted so that adding one
// more snapshot would not exceed the quota. Snapshots are ordered oldest-first;
// the caller is responsible for actually deleting them.
//
// If no quota is set for the profile, Evictions returns nil, nil.
// If the current count is within the limit, no evictions are needed.
func (m *Manager) Evictions(profile string, current []SnapshotMeta) ([]SnapshotMeta, error) {
	max, err := m.store.GetQuota(profile)
	if errors.Is(err, ErrNoQuota) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("quota: get quota for %q: %w", profile, err)
	}

	// We need room for one additional snapshot.
	excess := len(current) - max + 1
	if excess <= 0 {
		return nil, nil
	}

	// Sort oldest first so we evict the most stale entries.
	sorted := make([]SnapshotMeta, len(current))
	copy(sorted, current)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})

	return sorted[:excess], nil
}
