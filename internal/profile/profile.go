package profile

import (
	"fmt"

	"github.com/user/envport/internal/env"
	"github.com/user/envport/internal/snapshot"
	"github.com/user/envport/internal/store"
)

// Manager handles named environment profiles backed by a store.
type Manager struct {
	st *store.Store
}

// New returns a Manager using the given store.
func New(st *store.Store) *Manager {
	return &Manager{st: st}
}

// Save captures current environment variables and saves them under name.
func (m *Manager) Save(name string) error {
	vars := env.Capture()
	snap := snapshot.New(name, vars)
	data, err := snap.Marshal()
	if err != nil {
		return fmt.Errorf("profile save: %w", err)
	}
	return m.st.Set(name, data)
}

// Load retrieves the named profile and returns its snapshot.
func (m *Manager) Load(name string) (*snapshot.Snapshot, error) {
	data, err := m.st.Get(name)
	if err != nil {
		return nil, fmt.Errorf("profile load: %w", err)
	}
	snap, err := snapshot.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("profile load: %w", err)
	}
	return snap, nil
}

// Apply loads the named profile and applies it to the current process environment.
func (m *Manager) Apply(name string) error {
	snap, err := m.Load(name)
	if err != nil {
		return err
	}
	return env.Apply(snap.Vars)
}

// Delete removes a named profile from the store.
func (m *Manager) Delete(name string) error {
	return m.st.Delete(name)
}

// List returns all saved profile names.
func (m *Manager) List() ([]string, error) {
	return m.st.List()
}

// Diff returns the difference between the named profile and the current environment.
func (m *Manager) Diff(name string) (env.DiffResult, error) {
	snap, err := m.Load(name)
	if err != nil {
		return env.DiffResult{}, err
	}
	current := env.Capture()
	return env.Diff(snap.Vars, current), nil
}
