// Package freeze provides functionality to mark snapshots as immutable,
// preventing modification or deletion until explicitly unfrozen.
package freeze

import "fmt"

// Store is the persistence layer for frozen snapshot names.
type Store interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	List() (map[string]string, error)
}

const frozenValue = "frozen"

// Manager manages frozen snapshot state.
type Manager struct {
	store Store
}

// New returns a new Manager backed by the provided Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Freeze marks the named snapshot as frozen.
func (m *Manager) Freeze(name string) error {
	return m.store.Set(name, frozenValue)
}

// Unfreeze removes the frozen mark from the named snapshot.
func (m *Manager) Unfreeze(name string) error {
	v, err := m.store.Get(name)
	if err != nil || v == "" {
		return fmt.Errorf("snapshot %q is not frozen", name)
	}
	return m.store.Delete(name)
}

// IsFrozen reports whether the named snapshot is currently frozen.
func (m *Manager) IsFrozen(name string) (bool, error) {
	v, err := m.store.Get(name)
	if err != nil {
		return false, err
	}
	return v == frozenValue, nil
}

// List returns the names of all frozen snapshots.
func (m *Manager) List() ([]string, error) {
	all, err := m.store.List()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(all))
	for k, v := range all {
		if v == frozenValue {
			out = append(out, k)
		}
	}
	return out, nil
}
