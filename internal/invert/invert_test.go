package invert_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/envport/internal/invert"
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

func TestRunInvertsKeyValues(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"FOO": "bar", "BAZ": "qux"})

	r, err := invert.Run(m, "src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count != 2 {
		t.Fatalf("expected 2 entries, got %d", r.Count)
	}

	got, _ := m.Load("dst")
	if got["bar"] != "FOO" {
		t.Errorf("expected bar→FOO, got %q", got["bar"])
	}
	if got["qux"] != "BAZ" {
		t.Errorf("expected qux→BAZ, got %q", got["qux"])
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := invert.Run(m, "missing", "dst")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRunConflictDetected(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"A": "same", "B": "same"})

	r, err := invert.Run(m, "src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Conflicts) == 0 {
		t.Fatal("expected at least one conflict")
	}
}

func TestFormat(t *testing.T) {
	r := invert.Result{Source: "s", Destination: "d", Count: 3}
	out := invert.Format(r)
	if !strings.Contains(out, "3 key(s)") {
		t.Errorf("unexpected format output: %q", out)
	}
}

func TestFormatWithConflicts(t *testing.T) {
	r := invert.Result{
		Source: "s", Destination: "d", Count: 2,
		Conflicts: []string{`A and B share value "dup"`},
	}
	out := invert.Format(r)
	if !strings.Contains(out, "conflict") {
		t.Errorf("expected conflict in output: %q", out)
	}
}
