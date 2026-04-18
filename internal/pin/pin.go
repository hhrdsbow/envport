package pin

import (
	"errors"
	"sort"
)

// ErrNotFound is returned when a pinned snapshot is not found.
var ErrNotFound = errors.New("pin not found")

// Store is the persistence interface for pins.
type Store interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	List() (map[string]string, error)
}

// Manager manages pinned snapshots (alias -> snapshot name).
type Manager struct {
	store Store
}

// New creates a new pin Manager.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Pin associates an alias with a snapshot name.
func (m *Manager) Pin(alias, snapshot string) error {
	if alias == "" {
		return errors.New("alias must not be empty")
	}
	if snapshot == "" {
		return errors.New("snapshot must not be empty")
	}
	return m.store.Set(alias, snapshot)
}

// Resolve returns the snapshot name for the given alias.
func (m *Manager) Resolve(alias string) (string, error) {
	v, err := m.store.Get(alias)
	if err != nil {
		return "", ErrNotFound
	}
	return v, nil
}

// Unpin removes the alias.
func (m *Manager) Unpin(alias string) error {
	return m.store.Delete(alias)
}

// List returns all alias -> snapshot mappings sorted by alias.
func (m *Manager) List() ([][2]string, error) {
	all, err := m.store.List()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([][2]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, [2]string{k, all[k]})
	}
	return out, nil
}
