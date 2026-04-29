package unique_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/unique"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	store map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{store: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.store[name]
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
	m.store[name] = vars
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.store[name] = vars
}

func TestRunNoDuplicates(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1", "B": "2", "C": "3"})

	r, err := unique.Run(m, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Dropped) != 0 {
		t.Errorf("expected no drops, got %v", r.Dropped)
	}
	if len(r.Kept) != 3 {
		t.Errorf("expected 3 kept, got %d", len(r.Kept))
	}
}

func TestRunRemovesDuplicateValues(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "same", "B": "same", "C": "other"})

	r, err := unique.Run(m, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Dropped) != 1 {
		t.Errorf("expected 1 dropped, got %d", len(r.Dropped))
	}
	// B should be dropped because A comes first alphabetically.
	if r.Dropped[0] != "B" {
		t.Errorf("expected B to be dropped, got %s", r.Dropped[0])
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := unique.Run(m, "missing")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := unique.Result{Kept: []string{"A"}, Dropped: []string{"B", "C"}}
	out := unique.Format(r)
	if out == "" {
		t.Fatal("expected non-empty format output")
	}
	for _, k := range r.Dropped {
		if !contains(out, k) {
			t.Errorf("expected %q in format output", k)
		}
	}
}

func TestFormatNoDuplicates(t *testing.T) {
	r := unique.Result{Kept: []string{"A", "B"}}
	out := unique.Format(r)
	if out != "no duplicate values found" {
		t.Errorf("unexpected message: %q", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
