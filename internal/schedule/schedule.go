// Package schedule provides functionality for managing named snapshot schedules.
package schedule

import (
	"encoding/json"
	"fmt"
	"time"
)

// Entry represents a scheduled snapshot configuration.
type Entry struct {
	Name     string        `json:"name"`
	Profile  string        `json:"profile"`
	Interval time.Duration `json:"interval"`
	LastRun  time.Time     `json:"last_run"`
	Enabled  bool          `json:"enabled"`
}

// Store is the persistence interface for schedule entries.
type Store interface {
	Set(name string, data []byte) error
	Get(name string) ([]byte, error)
	Delete(name string) error
	List() ([]string, error)
}

// Manager handles schedule CRUD operations.
type Manager struct {
	store Store
}

// New returns a new Manager backed by the given Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Add creates or replaces a schedule entry.
func (m *Manager) Add(e Entry) error {
	if e.Name == "" {
		return fmt.Errorf("schedule name must not be empty")
	}
	if e.Profile == "" {
		return fmt.Errorf("profile must not be empty")
	}
	if e.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	e.Enabled = true
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return m.store.Set(e.Name, data)
}

// Get retrieves a schedule entry by name.
func (m *Manager) Get(name string) (Entry, error) {
	data, err := m.store.Get(name)
	if err != nil {
		return Entry{}, fmt.Errorf("schedule %q not found", name)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, err
	}
	return e, nil
}

// Remove deletes a schedule entry by name.
func (m *Manager) Remove(name string) error {
	return m.store.Delete(name)
}

// List returns all schedule entry names.
func (m *Manager) List() ([]Entry, error) {
	names, err := m.store.List()
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, n := range names {
		e, err := m.Get(n)
		if err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// UpdateLastRun records the time a schedule was last executed.
func (m *Manager) UpdateLastRun(name string, t time.Time) error {
	e, err := m.Get(name)
	if err != nil {
		return err
	}
	e.LastRun = t
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return m.store.Set(name, data)
}
