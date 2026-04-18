// Package audit records snapshot operations for accountability.
package audit

import (
	"encoding/json"
	"fmt"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Profile   string    `json:"profile"`
	Detail    string    `json:"detail,omitempty"`
}

// Store is the persistence interface for audit entries.
type Store interface {
	Append(entry Entry) error
	List() ([]Entry, error)
}

// Manager handles audit logging.
type Manager struct {
	store Store
}

// New creates a new Manager backed by store.
func New(store Store) *Manager {
	return &Manager{store: store}
}

// Record writes an audit entry for the given operation and profile.
func (m *Manager) Record(operation, profile, detail string) error {
	return m.store.Append(Entry{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		Profile:   profile,
		Detail:    detail,
	})
}

// List returns all recorded audit entries.
func (m *Manager) List() ([]Entry, error) {
	return m.store.List()
}

// Format renders an entry as a human-readable string.
func Format(e Entry) string {
	ts := e.Timestamp.Format(time.RFC3339)
	if e.Detail != "" {
		return fmt.Sprintf("%s  %-10s  %-20s  %s", ts, e.Operation, e.Profile, e.Detail)
	}
	return fmt.Sprintf("%s  %-10s  %s", ts, e.Operation, e.Profile)
}

// MarshalEntry serialises an entry to JSON.
func MarshalEntry(e Entry) ([]byte, error) {
	return json.Marshal(e)
}
