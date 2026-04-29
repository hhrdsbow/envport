package flatten_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/flatten"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	store map[string]map[string]string
}

func newManager() *memManager {
	return &memManager{store: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.store[name]
	if !ok {
		return nil, errors.New("not found")
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.store[name] = vars
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.store[name] = vars
}

func TestRunNoNames(t *testing.T) {
	m := newManager()
	_, err := flatten.Run(m, nil, "out")
	if err == nil {
		t.Fatal("expected error for empty names")
	}
}

func TestRunMergesInOrder(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seed(m, "b", map[string]string{"BAR": "99", "BAZ": "3"})

	res, err := flatten.Run(m, []string{"a", "b"}, "merged")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "1" {
		t.Errorf("FOO: got %q want %q", res.Vars["FOO"], "1")
	}
	if res.Vars["BAR"] != "99" {
		t.Errorf("BAR: got %q want %q", res.Vars["BAR"], "99")
	}
	if res.Vars["BAZ"] != "3" {
		t.Errorf("BAZ: got %q want %q", res.Vars["BAZ"], "3")
	}
	if res.Conflicts != 1 {
		t.Errorf("conflicts: got %d want 1", res.Conflicts)
	}
}

func TestRunSavesPersisted(t *testing.T) {
	m := newManager()
	seed(m, "src", map[string]string{"X": "hello"})

	_, err := flatten.Run(m, []string{"src"}, "dest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := m.Load("dest")
	if err != nil {
		t.Fatalf("load dest: %v", err)
	}
	if loaded["X"] != "hello" {
		t.Errorf("X: got %q want %q", loaded["X"], "hello")
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newManager()
	_, err := flatten.Run(m, []string{"missing"}, "out")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}
