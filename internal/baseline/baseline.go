// Package baseline provides functionality to mark a snapshot as the
// baseline for a profile, enabling drift detection relative to that point.
package baseline

import (
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned when no baseline exists for a profile.
var ErrNotFound = errors.New("baseline: no baseline set")

// Entry records the snapshot name pinned as the baseline and when it was set.
type Entry struct {
	Profile     string    `json:"profile"`
	Snapshot    string    `json:"snapshot"`
	SetAt       time.Time `json:"set_at"`
}

// Store is the persistence interface required by Manager.
type Store interface {
	Set(profile, snapshot string, at time.Time) error
	Get(profile string) (Entry, error)
	Delete(profile string) error
	List() ([]Entry, error)
}

// Manager manages baseline entries.
type Manager struct {
	store Store
}

// New creates a new Manager backed by the provided Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Set marks snapshot as the baseline for profile.
func (m *Manager) Set(profile, snapshot string) error {
	if profile == "" {
		return fmt.Errorf("baseline: profile name required")
	}
	if snapshot == "" {
		return fmt.Errorf("baseline: snapshot name required")
	}
	return m.store.Set(profile, snapshot, time.Now().UTC())
}

// Get returns the current baseline entry for profile.
func (m *Manager) Get(profile string) (Entry, error) {
	return m.store.Get(profile)
}

// Clear removes the baseline for profile.
func (m *Manager) Clear(profile string) error {
	return m.store.Delete(profile)
}

// List returns all baseline entries.
func (m *Manager) List() ([]Entry, error) {
	return m.store.List()
}
