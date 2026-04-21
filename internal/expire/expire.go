// Package expire provides snapshot expiration management.
package expire

import (
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned when no expiry is set for a snapshot.
var ErrNotFound = errors.New("no expiry set for snapshot")

// Manager handles expiry metadata for snapshots.
type Manager struct {
	store Store
}

// Store is the persistence interface for expiry records.
type Store interface {
	Set(name string, expiresAt time.Time) error
	Get(name string) (time.Time, error)
	Delete(name string) error
	List() (map[string]time.Time, error)
}

// New creates a new Manager backed by the given store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Set assigns an expiry duration to a named snapshot.
func (m *Manager) Set(name string, ttl time.Duration) error {
	if name == "" {
		return errors.New("snapshot name must not be empty")
	}
	if ttl <= 0 {
		return errors.New("ttl must be positive")
	}
	return m.store.Set(name, time.Now().Add(ttl))
}

// Get returns the expiry time for a snapshot.
func (m *Manager) Get(name string) (time.Time, error) {
	return m.store.Get(name)
}

// Remove clears the expiry for a snapshot.
func (m *Manager) Remove(name string) error {
	return m.store.Delete(name)
}

// Expired returns all snapshot names whose expiry has passed.
func (m *Manager) Expired() ([]string, error) {
	all, err := m.store.List()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var out []string
	for name, t := range all {
		if now.After(t) {
			out = append(out, name)
		}
	}
	return out, nil
}

// Format returns a human-readable summary of all expiry records.
func Format(records map[string]time.Time) string {
	if len(records) == 0 {
		return "no expiry records found\n"
	}
	out := ""
	for name, t := range records {
		status := "active"
		if time.Now().After(t) {
			status = "expired"
		}
		out += fmt.Sprintf("%-24s  %s  [%s]\n", name, t.Format(time.RFC3339), status)
	}
	return out
}
