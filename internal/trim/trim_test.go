package trim_test

import (
	"errors"
	"testing"

	"envport/internal/trim"
)

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

func TestRunRemovesKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1", "B": "2", "C": "3"})

	res, err := trim.Run(m, "dev", []string{"A", "C"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if res.Kept != 1 {
		t.Errorf("expected 1 kept, got %d", res.Kept)
	}
	vars, _ := m.Load("dev")
	if _, ok := vars["A"]; ok {
		t.Error("key A should have been removed")
	}
}

func TestRunDryRunDoesNotPersist(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"X": "10", "Y": "20"})

	res, err := trim.Run(m, "dev", []string{"X"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 in removed list, got %d", len(res.Removed))
	}
	vars, _ := m.Load("dev")
	if _, ok := vars["X"]; !ok {
		t.Error("key X should still exist after dry run")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := trim.Run(m, "ghost", []string{"A"}, false)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunNoKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})
	_, err := trim.Run(m, "dev", nil, false)
	if err == nil {
		t.Fatal("expected error when no keys provided")
	}
}

func TestRunIgnoresMissingKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})
	res, err := trim.Run(m, "dev", []string{"NONEXISTENT"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(res.Removed))
	}
}
