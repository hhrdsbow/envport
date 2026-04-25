// Package group provides grouping of snapshots under named collections.
package group

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrNotFound is returned when a group does not exist.
var ErrNotFound = errors.New("group not found")

// ErrAlreadyExists is returned when a group already exists.
var ErrAlreadyExists = errors.New("group already exists")

// Store is the persistence interface for groups.
type Store interface {
	Get(name string) ([]string, error)
	Set(name string, members []string) error
	Delete(name string) error
	List() ([]string, error)
}

// Manager manages snapshot groups.
type Manager struct {
	store Store
}

// New returns a Manager backed by store.
func New(store Store) *Manager {
	return &Manager{store: store}
}

// Create initialises a new empty group.
func (m *Manager) Create(name string) error {
	if _, err := m.store.Get(name); err == nil {
		return fmt.Errorf("%w: %s", ErrAlreadyExists, name)
	}
	return m.store.Set(name, []string{})
}

// Add appends a snapshot name to the group, creating the group if absent.
func (m *Manager) Add(group, snapshot string) error {
	members, err := m.store.Get(group)
	if err != nil {
		members = []string{}
	}
	for _, m := range members {
		if m == snapshot {
			return nil // idempotent
		}
	}
	members = append(members, snapshot)
	return m.store.Set(group, members)
}

// Remove removes a snapshot name from the group.
func (m *Manager) Remove(group, snapshot string) error {
	members, err := m.store.Get(group)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, group)
	}
	next := members[:0]
	for _, mem := range members {
		if mem != snapshot {
			next = append(next, mem)
		}
	}
	return m.store.Set(group, next)
}

// Members returns the snapshot names belonging to group.
func (m *Manager) Members(group string) ([]string, error) {
	members, err := m.store.Get(group)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, group)
	}
	out := make([]string, len(members))
	copy(out, members)
	return out, nil
}

// Delete removes a group entirely.
func (m *Manager) Delete(name string) error {
	if _, err := m.store.Get(name); err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	return m.store.Delete(name)
}

// List returns all group names.
func (m *Manager) List() ([]string, error) {
	return m.store.List()
}

// MarshalMembers serialises members to JSON for storage.
func MarshalMembers(members []string) ([]byte, error) {
	return json.Marshal(members)
}

// UnmarshalMembers deserialises members from JSON.
func UnmarshalMembers(data []byte) ([]string, error) {
	var out []string
	err := json.Unmarshal(data, &out)
	return out, err
}
