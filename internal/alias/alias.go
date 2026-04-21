// Package alias manages short aliases for snapshot names.
package alias

import "fmt"

// Store is the persistence interface for alias mappings.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	List() (map[string]string, error)
}

// Manager handles alias operations.
type Manager struct {
	store Store
}

// New returns a new Manager backed by store.
func New(store Store) *Manager {
	return &Manager{store: store}
}

// Add creates or overwrites an alias pointing to target.
func (m *Manager) Add(alias, target string) error {
	if alias == "" || target == "" {
		return fmt.Errorf("alias and target must not be empty")
	}
	return m.store.Set(alias, target)
}

// Remove deletes an alias.
func (m *Manager) Remove(alias string) error {
	return m.store.Delete(alias)
}

// Resolve returns the snapshot name an alias points to.
// If the alias does not exist it returns ErrNotFound.
func (m *Manager) Resolve(alias string) (string, error) {
	target, err := m.store.Get(alias)
	if err != nil {
		return "", fmt.Errorf("alias %q not found: %w", alias, err)
	}
	return target, nil
}

// List returns all alias → target pairs.
func (m *Manager) List() (map[string]string, error) {
	return m.store.List()
}
