package prefix_test

import (
	"errors"
	"testing"

	"envport/internal/prefix"
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

func TestRunAddAllKeys(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "snap", map[string]string{"FOO": "1", "BAR": "2"})

	res, err := prefix.RunAdd(mgr, "snap", "snap", "APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(res.Added))
	}
	vars, _ := mgr.Load("snap")
	if vars["APP_FOO"] != "1" || vars["APP_BAR"] != "2" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestRunAddEmptyPrefix(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "snap", map[string]string{"FOO": "1"})
	_, err := prefix.RunAdd(mgr, "snap", "snap", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestRunRemoveStripsPrefix(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "snap", map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"})

	res, err := prefix.RunRemove(mgr, "snap", "out", "APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	vars, _ := mgr.Load("out")
	if _, ok := vars["FOO"]; !ok {
		t.Errorf("FOO should be present after strip")
	}
	if _, ok := vars["OTHER"]; !ok {
		t.Errorf("OTHER should be preserved")
	}
}

func TestRunRemoveMissingSnapshot(t *testing.T) {
	mgr := newMemManager()
	_, err := prefix.RunRemove(mgr, "ghost", "out", "APP_")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := prefix.Result{Added: []string{"A", "B"}, Skipped: []string{"C"}}
	out := prefix.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
