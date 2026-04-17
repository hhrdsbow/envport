package tag

import (
	"fmt"
	"sort"
	"strings"
)

// Manager handles tagging of snapshots.
type Manager struct {
	store tagStore
}

type tagStore interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	Delete(key string) error
	List() ([]string, error)
}

const prefix = "tag:"

// New returns a new Manager backed by the given store.
func New(s tagStore) *Manager {
	return &Manager{store: s}
}

// Add associates a tag with a snapshot name.
func (m *Manager) Add(tag, snapshot string) error {
	key := prefix + tag
	existing, _ := m.Get(tag)
	for _, s := range existing {
		if s == snapshot {
			return nil
		}
	}
	existing = append(existing, snapshot)
	sort.Strings(existing)
	return m.store.Set(key, []byte(strings.Join(existing, "\n")))
}

// Remove disassociates a tag from a snapshot name.
func (m *Manager) Remove(tag, snapshot string) error {
	existing, err := m.Get(tag)
	if err != nil {
		return err
	}
	updated := existing[:0]
	for _, s := range existing {
		if s != snapshot {
			updated = append(updated, s)
		}
	}
	if len(updated) == 0 {
		return m.store.Delete(prefix + tag)
	}
	return m.store.Set(prefix+tag, []byte(strings.Join(updated, "\n")))
}

// Get returns all snapshot names associated with a tag.
func (m *Manager) Get(tag string) ([]string, error) {
	data, err := m.store.Get(prefix + tag)
	if err != nil {
		return nil, fmt.Errorf("tag %q not found", tag)
	}
	if len(data) == 0 {
		return nil, nil
	}
	return strings.Split(string(data), "\n"), nil
}

// List returns all known tags.
func (m *Manager) List() ([]string, error) {
	keys, err := m.store.List()
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			tags = append(tags, strings.TrimPrefix(k, prefix))
		}
	}
	sort.Strings(tags)
	return tags, nil
}
