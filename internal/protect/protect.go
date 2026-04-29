// Package protect provides read-only protection for named snapshots.
// A protected snapshot cannot be overwritten or deleted until the
// protection is explicitly removed.
package protect

import (
	"errors"
	"fmt"
)

// ErrProtected is returned when an operation is attempted on a protected snapshot.
var ErrProtected = errors.New("snapshot is protected")

// Store persists protection records.
type Store interface {
	Set(name string) error
	Delete(name string) error
	Exists(name string) (bool, error)
	List() ([]string, error)
}

// Manager manages snapshot protection.
type Manager struct {
	store Store
}

// New creates a new Manager backed by the given Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Protect marks name as protected.
func (m *Manager) Protect(name string) error {
	ok, err := m.store.Exists(name)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%w: %s", ErrProtected, name)
	}
	return m.store.Set(name)
}

// Unprotect removes protection from name.
func (m *Manager) Unprotect(name string) error {
	ok, err := m.store.Exists(name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("snapshot %q is not protected", name)
	}
	return m.store.Delete(name)
}

// IsProtected reports whether name is protected.
func (m *Manager) IsProtected(name string) (bool, error) {
	return m.store.Exists(name)
}

// List returns all protected snapshot names.
func (m *Manager) List() ([]string, error) {
	return m.store.List()
}

// Guard returns ErrProtected if name is protected, nil otherwise.
func (m *Manager) Guard(name string) error {
	ok, err := m.store.Exists(name)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%w: %s", ErrProtected, name)
	}
	return nil
}
