// Package scope provides scoped environment variable views,
// allowing a named subset of keys to be isolated and managed together.
package scope

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNotFound is returned when a scope does not exist.
var ErrNotFound = errors.New("scope not found")

// Scope holds a named set of environment variable keys.
type Scope struct {
	Name string
	Keys []string
}

// Store is the persistence interface for scopes.
type Store interface {
	Get(name string) ([]string, error)
	Set(name string, keys []string) error
	Delete(name string) error
	List() ([]string, error)
}

// Manager manages scopes backed by a Store.
type Manager struct {
	store Store
}

// New returns a new Manager using the provided Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Add creates or replaces a scope with the given keys.
func (m *Manager) Add(name string, keys []string) error {
	if name == "" {
		return fmt.Errorf("scope name must not be empty")
	}
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	return m.store.Set(name, sorted)
}

// Get retrieves the keys belonging to a scope.
func (m *Manager) Get(name string) (*Scope, error) {
	keys, err := m.store.Get(name)
	if err != nil {
		return nil, err
	}
	return &Scope{Name: name, Keys: keys}, nil
}

// Delete removes a scope by name.
func (m *Manager) Delete(name string) error {
	return m.store.Delete(name)
}

// List returns all known scope names.
func (m *Manager) List() ([]string, error) {
	return m.store.List()
}

// Apply filters the provided vars map to only those keys in the scope.
func (m *Manager) Apply(name string, vars map[string]string) (map[string]string, error) {
	sc, err := m.Get(name)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(sc.Keys))
	for _, k := range sc.Keys {
		if v, ok := vars[k]; ok {
			out[k] = v
		}
	}
	return out, nil
}
