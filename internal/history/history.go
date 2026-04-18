package history

import (
	"encoding/json"
	"fmt"
	"time"
)

// Entry records a single operation performed on a snapshot.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Snapshot  string    `json:"snapshot"`
	Detail    string    `json:"detail,omitempty"`
}

// Store is the persistence interface required by Manager.
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
}

// Manager records and retrieves operation history.
type Manager struct {
	store Store
}

const storeKey = "__history__"

// New returns a new Manager backed by store.
func New(store Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) load() ([]Entry, error) {
	data, err := m.store.Get(storeKey)
	if err != nil || data == nil {
		return []Entry{}, nil
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("history: decode: %w", err)
	}
	return entries, nil
}

func (m *Manager) save(entries []Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("history: encode: %w", err)
	}
	return m.store.Set(storeKey, data)
}

// Record appends a new entry to the history log.
func (m *Manager) Record(op, snapshot, detail string) error {
	entries, err := m.load()
	if err != nil {
		return err
	}
	entries = append(entries, Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Snapshot:  snapshot,
		Detail:    detail,
	})
	return m.save(entries)
}

// List returns all recorded history entries.
func (m *Manager) List() ([]Entry, error) {
	return m.load()
}

// Clear removes all history entries.
func (m *Manager) Clear() error {
	return m.save([]Entry{})
}
