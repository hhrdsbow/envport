package unset_test

import (
	"errors"
	"testing"

	"envport/internal/unset"
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

func TestRunRemovesKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1", "B": "2", "C": "3"})

	res, err := unset.Run(m, "dev", []string{"A", "C"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
	vars, _ := m.Load("dev")
	if _, ok := vars["A"]; ok {
		t.Error("key A should have been removed")
	}
	if _, ok := vars["B"]; !ok {
		t.Error("key B should still exist")
	}
}

func TestRunNotFoundStrict(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})

	_, err := unset.Run(m, "dev", []string{"MISSING"}, true)
	if err == nil {
		t.Fatal("expected error for missing key in strict mode")
	}
}

func TestRunNotFoundLenient(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})

	res, err := unset.Run(m, "dev", []string{"MISSING"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.NotFound) != 1 || res.NotFound[0] != "MISSING" {
		t.Errorf("expected MISSING in NotFound, got %v", res.NotFound)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := unset.Run(m, "ghost", []string{"A"}, false)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormatOutput(t *testing.T) {
	res := &unset.Result{
		Removed:  []string{"FOO", "BAR"},
		NotFound: []string{"BAZ"},
	}
	out := unset.Format(res)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
