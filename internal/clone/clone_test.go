package clone_test

import (
	"fmt"
	"testing"
	"time"

	"envport/internal/clone"
	"envport/internal/snapshot"
)

type memManager struct {
	data map[string]*snapshot.Snapshot
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]*snapshot.Snapshot)}
}

func (m *memManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.data[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return s, nil
}

func (m *memManager) Save(name string, s *snapshot.Snapshot) error {
	m.data[name] = s
	return nil
}

func (m *memManager) List() ([]string, error) {
	out := make([]string, 0, len(m.data))
	for k := range m.data {
		out = append(out, k)
	}
	return out, nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = &snapshot.Snapshot{Name: name, Vars: vars, CreatedAt: time.Now()}
}

func TestRunClonesSnapshot(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"A": "1", "B": "2"})

	dest, err := clone.Run(m, "prod", "staging", clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "staging" {
		t.Fatalf("expected staging, got %s", dest)
	}
	s := m.data["staging"]
	if s.Vars["A"] != "1" || s.Vars["B"] != "2" {
		t.Fatal("vars not copied correctly")
	}
}

func TestRunAutoSuffix(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"X": "y"})
	seed(m, "prod-copy", map[string]string{})

	dest, err := clone.Run(m, "prod", "prod", clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "prod-copy-copy" {
		t.Fatalf("expected prod-copy-copy, got %s", dest)
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := clone.Run(m, "missing", "dest", clone.Options{})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRunEmptySourceName(t *testing.T) {
	m := newMemManager()
	_, err := clone.Run(m, "", "dest", clone.Options{})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}
