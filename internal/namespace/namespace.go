// Package namespace provides grouping of snapshots under named namespaces.
package namespace

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrNotFound is returned when a namespace does not exist.
var ErrNotFound = errors.New("namespace not found")

// Store is the persistence interface required by Manager.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	List() ([]string, error)
}

// Manager handles namespace CRUD operations.
type Manager struct {
	store Store
}

const prefix = "ns:"

// New returns a Manager backed by the given Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Add creates a new namespace. Returns an error if it already exists.
func (m *Manager) Add(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("namespace name must not be empty")
	}
	_, err := m.store.Get(prefix + name)
	if err == nil {
		return fmt.Errorf("namespace %q already exists", name)
	}
	return m.store.Set(prefix+name, name)
}

// Remove deletes a namespace by name.
func (m *Manager) Remove(name string) error {
	_, err := m.store.Get(prefix + name)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	return m.store.Delete(prefix + name)
}

// List returns all namespace names in sorted order.
func (m *Manager) List() ([]string, error) {
	keys, err := m.store.List()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			names = append(names, strings.TrimPrefix(k, prefix))
		}
	}
	sort.Strings(names)
	return names, nil
}

// Exists reports whether the given namespace is registered.
func (m *Manager) Exists(name string) bool {
	_, err := m.store.Get(prefix + name)
	return err == nil
}
