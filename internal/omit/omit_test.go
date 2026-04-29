package omit_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/omit"
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

	r, err := omit.Run(m, "dev", []string{"A", "C"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(r.Removed))
	}
	if r.Remaining != 1 {
		t.Fatalf("expected 1 remaining, got %d", r.Remaining)
	}
	vars, _ := m.Load("dev")
	if _, ok := vars["B"]; !ok {
		t.Error("expected key B to remain")
	}
}

func TestRunIgnoresMissingKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"X": "10"})

	r, err := omit.Run(m, "dev", []string{"NOPE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Removed) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(r.Removed))
	}
	if r.Remaining != 1 {
		t.Fatalf("expected 1 remaining, got %d", r.Remaining)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := omit.Run(m, "ghost", []string{"A"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunEmptyName(t *testing.T) {
	m := newMemManager()
	_, err := omit.Run(m, "", []string{"A"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestFormat(t *testing.T) {
	r := omit.Result{Removed: []string{"A", "B"}, Remaining: 3}
	out := omit.Format(r)
	if out == "" {
		t.Fatal("expected non-empty format output")
	}
}

func TestFormatNoRemovals(t *testing.T) {
	r := omit.Result{Removed: nil, Remaining: 5}
	out := omit.Format(r)
	if out == "" {
		t.Fatal("expected non-empty format output")
	}
}
