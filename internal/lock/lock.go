// Package lock provides snapshot locking to prevent concurrent modification.
package lock

import (
	"errors"
	"time"
)

// ErrAlreadyLocked is returned when a snapshot is already locked.
var ErrAlreadyLocked = errors.New("snapshot is locked")

// ErrNotLocked is returned when trying to unlock a snapshot that isn't locked.
var ErrNotLocked = errors.New("snapshot is not locked")

// Entry represents a lock record.
type Entry struct {
	Name      string
	LockedAt  time.Time
	Reason    string
}

// Store is the persistence interface for locks.
type Store interface {
	Set(name string, entry Entry) error
	Get(name string) (Entry, bool, error)
	Delete(name string) error
	List() ([]Entry, error)
}

// Manager handles lock operations.
type Manager struct {
	store Store
}

// New creates a new Manager.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Lock locks a snapshot with an optional reason.
func (m *Manager) Lock(name, reason string) error {
	_, exists, err := m.store.Get(name)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyLocked
	}
	return m.store.Set(name, Entry{
		Name:     name,
		LockedAt: time.Now().UTC(),
		Reason:   reason,
	})
}

// Unlock removes the lock from a snapshot.
func (m *Manager) Unlock(name string) error {
	_, exists, err := m.store.Get(name)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotLocked
	}
	return m.store.Delete(name)
}

// IsLocked reports whether a snapshot is locked.
func (m *Manager) IsLocked(name string) (bool, error) {
	_, exists, err := m.store.Get(name)
	return exists, err
}

// List returns all current locks.
func (m *Manager) List() ([]Entry, error) {
	return m.store.List()
}
