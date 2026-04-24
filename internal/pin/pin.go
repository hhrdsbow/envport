// Package pin manages "pinned" snapshots — snapshots that are protected from
// automated pruning or expiry operations.
package pin

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNotPinned is returned when an unpin is attempted on a name that is not
// currently pinned.
var ErrNotPinned = errors.New("pin: snapshot is not pinned")

// Store is the persistence interface required by the pin manager.
type Store interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	List() (map[string]string, error)
}

const prefix = "pin:"

// Manager wraps a Store and provides pin / unpin / list operations.
type Manager struct {
	store Store
}

// New returns a Manager backed by the provided Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Add marks the named snapshot as pinned.
func (m *Manager) Add(name string) error {
	if err := m.store.Set(prefix+name, "1"); err != nil {
		return fmt.Errorf("pin add %q: %w", name, err)
	}
	return nil
}

// Remove unmarks the named snapshot.
func (m *Manager) Remove(name string) error {
	if _, err := m.store.Get(prefix + name); err != nil {
		return ErrNotPinned
	}
	if err := m.store.Delete(prefix + name); err != nil {
		return fmt.Errorf("pin remove %q: %w", name, err)
	}
	return nil
}

// IsPinned reports whether the named snapshot is currently pinned.
func (m *Manager) IsPinned(name string) bool {
	_, err := m.store.Get(prefix + name)
	return err == nil
}

// List returns the names of all pinned snapshots in sorted order.
func (m *Manager) List() ([]string, error) {
	all, err := m.store.List()
	if err != nil {
		return nil, fmt.Errorf("pin list: %w", err)
	}
	var names []string
	for k := range all {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			names = append(names, k[len(prefix):])
		}
	}
	sort.Strings(names)
	return names, nil
}
