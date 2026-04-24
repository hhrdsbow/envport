// Package annotate provides functionality to attach and retrieve
// human-readable notes/annotations on named snapshots.
package annotate

import "fmt"

// Store is the persistence interface required by the annotation manager.
type Store interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	List() (map[string]string, error)
}

// Manager handles annotation operations.
type Manager struct {
	store Store
}

// New returns a new Manager backed by the given Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

func annotationKey(name string) string {
	return "annotation:" + name
}

// Set attaches a note to the named snapshot.
func (m *Manager) Set(name, note string) error {
	if name == "" {
		return fmt.Errorf("annotate: name must not be empty")
	}
	return m.store.Set(annotationKey(name), note)
}

// Get retrieves the note attached to the named snapshot.
// Returns an empty string and no error when no annotation exists.
func (m *Manager) Get(name string) (string, error) {
	val, err := m.store.Get(annotationKey(name))
	if err != nil {
		// Treat a missing key as an empty annotation.
		return "", nil
	}
	return val, nil
}

// Remove deletes the annotation for the named snapshot.
func (m *Manager) Remove(name string) error {
	if name == "" {
		return fmt.Errorf("annotate: name must not be empty")
	}
	return m.store.Delete(annotationKey(name))
}

// List returns a map of snapshot name → annotation for all annotated snapshots.
func (m *Manager) List() (map[string]string, error) {
	all, err := m.store.List()
	if err != nil {
		return nil, err
	}
	prefix := "annotation:"
	out := make(map[string]string)
	for k, v := range all {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out[k[len(prefix):]] = v
		}
	}
	return out, nil
}
