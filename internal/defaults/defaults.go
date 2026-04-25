// Package defaults manages per-profile default key-value pairs that are
// automatically merged into a snapshot when it is loaded.
package defaults

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned when no defaults entry exists for a profile.
var ErrNotFound = errors.New("defaults: profile not found")

// Store is the persistence interface required by Manager.
type Store interface {
	Get(profile string) (map[string]string, error)
	Set(profile string, vals map[string]string) error
	Delete(profile string) error
	List() ([]string, error)
}

// Manager handles default key-value pairs for named profiles.
type Manager struct {
	store Store
}

// New returns a Manager backed by store.
func New(store Store) *Manager {
	return &Manager{store: store}
}

// Set stores defaults for profile, merging with any existing entries.
func (m *Manager) Set(profile string, pairs map[string]string) error {
	if profile == "" {
		return fmt.Errorf("defaults: profile name must not be empty")
	}
	existing, err := m.store.Get(profile)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	if existing == nil {
		existing = make(map[string]string)
	}
	for k, v := range pairs {
		existing[k] = v
	}
	return m.store.Set(profile, existing)
}

// Get returns the defaults map for profile.
func (m *Manager) Get(profile string) (map[string]string, error) {
	return m.store.Get(profile)
}

// Remove deletes specific keys from a profile's defaults.
func (m *Manager) Remove(profile string, keys []string) error {
	existing, err := m.store.Get(profile)
	if err != nil {
		return err
	}
	for _, k := range keys {
		delete(existing, k)
	}
	if len(existing) == 0 {
		return m.store.Delete(profile)
	}
	return m.store.Set(profile, existing)
}

// Clear removes all defaults for profile.
func (m *Manager) Clear(profile string) error {
	return m.store.Delete(profile)
}

// List returns all profile names that have defaults registered.
func (m *Manager) List() ([]string, error) {
	return m.store.List()
}

// Apply merges the stored defaults for profile into vars, returning a new map.
// Existing keys in vars are NOT overwritten.
func (m *Manager) Apply(profile string, vars map[string]string) (map[string]string, error) {
	defaults, err := m.store.Get(profile)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return vars, nil
		}
		return nil, err
	}
	out := make(map[string]string, len(vars)+len(defaults))
	for k, v := range defaults {
		out[k] = v
	}
	for k, v := range vars {
		out[k] = v
	}
	return out, nil
}
