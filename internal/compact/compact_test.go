package compact_test

import (
	"errors"
	"testing"

	"envport/internal/compact"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	data map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = vars
}

func TestRunNoDuplicates(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1", "B": "2", "C": "3"})

	r, err := compact.Run(m, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removed keys, got %v", r.Removed)
	}
	if r.Kept != 3 {
		t.Errorf("expected 3 kept, got %d", r.Kept)
	}
}

func TestRunRemovesDuplicateValues(t *testing.T) {
	m := newMemManager()
	// B and A share value "same"; C is unique.
	seed(m, "prod", map[string]string{"A": "same", "B": "same", "C": "unique"})

	r, err := compact.Run(m, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %v", r.Removed)
	}
	// "A" is removed; "B" (last alphabetically) is kept.
	if r.Removed[0] != "A" {
		t.Errorf("expected A removed, got %s", r.Removed[0])
	}
	if r.Kept != 2 {
		t.Errorf("expected 2 kept, got %d", r.Kept)
	}

	vars, _ := m.Load("prod")
	if _, ok := vars["A"]; ok {
		t.Error("A should have been removed")
	}
	if v, ok := vars["B"]; !ok || v != "same" {
		t.Error("B should be kept with value 'same'")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := compact.Run(m, "missing")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := compact.Result{Removed: []string{"OLD_KEY"}, Kept: 5}
	out := compact.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormatNoRemovals(t *testing.T) {
	r := compact.Result{Removed: nil, Kept: 4}
	out := compact.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
