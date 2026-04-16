package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultStoreFile = ".envport.json"

// Store manages named snapshot references on disk.
type Store struct {
	path string
	data map[string]string // name -> snapshot file path
}

// Open loads or creates a store at the given directory.
func Open(dir string) (*Store, error) {
	p := filepath.Join(dir, defaultStoreFile)
	s := &Store{path: p, data: make(map[string]string)}

	bytes, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &s.data); err != nil {
		return nil, err
	}
	return s, nil
}

// Set registers a named snapshot path.
func (s *Store) Set(name, snapshotPath string) error {
	s.data[name] = snapshotPath
	return s.save()
}

// Get returns the snapshot path for the given name.
func (s *Store) Get(name string) (string, bool) {
	v, ok := s.data[name]
	return v, ok
}

// Delete removes a named snapshot entry.
func (s *Store) Delete(name string) error {
	delete(s.data, name)
	return s.save()
}

// List returns all registered snapshot names.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.data))
	for k := range s.data {
		names = append(names, k)
	}
	return names
}

func (s *Store) save() error {
	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, bytes, 0644)
}
