package dedup_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/dedup"
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
	seed(m, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seed(m, "b", map[string]string{"BAZ": "3"})

	r, err := dedup.Run(m, []string{"a", "b"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Dropped) != 0 {
		t.Errorf("expected no dropped keys, got %v", r.Dropped)
	}
	if len(r.Kept) != 3 {
		t.Errorf("expected 3 kept keys, got %d", len(r.Kept))
	}
}

func TestRunDropsDuplicates(t *testing.T) {
	m := newMemManager()
	seed(m, "high", map[string]string{"FOO": "high", "ONLY_HIGH": "yes"})
	seed(m, "low", map[string]string{"FOO": "low", "ONLY_LOW": "yes"})

	r, err := dedup.Run(m, []string{"high", "low"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Kept["FOO"] != "high" {
		t.Errorf("expected FOO=high, got %q", r.Kept["FOO"])
	}
	if len(r.Dropped) != 1 || r.Dropped[0] != "FOO" {
		t.Errorf("expected [FOO] dropped, got %v", r.Dropped)
	}
}

func TestRunSavesDst(t *testing.T) {
	m := newMemManager()
	seed(m, "x", map[string]string{"A": "1"})
	seed(m, "y", map[string]string{"B": "2"})

	_, err := dedup.Run(m, []string{"x", "y"}, "merged")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := m.Load("merged")
	if err != nil {
		t.Fatalf("merged snapshot not saved: %v", err)
	}
	if v["A"] != "1" || v["B"] != "2" {
		t.Errorf("unexpected merged content: %v", v)
	}
}

func TestRunTooFewNames(t *testing.T) {
	m := newMemManager()
	_, err := dedup.Run(m, []string{"only"}, "")
	if err == nil {
		t.Fatal("expected error for fewer than two names")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	seed(m, "exists", map[string]string{"K": "v"})
	_, err := dedup.Run(m, []string{"exists", "missing"}, "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := dedup.Result{Kept: map[string]string{"A": "1"}, Dropped: []string{}}
	out := dedup.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
