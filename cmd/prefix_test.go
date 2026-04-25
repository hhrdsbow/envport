package cmd_test

import (
	"bytes"
	"errors"
	"testing"

	"envport/internal/prefix"
)

// memPrefixManager satisfies prefix.Manager for command-level tests.
type memPrefixManager struct {
	data map[string]map[string]string
}

func newMemPrefixManager() *memPrefixManager {
	return &memPrefixManager{data: make(map[string]map[string]string)}
}

func (m *memPrefixManager) Load(name string) (map[string]string, error) {
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

func (m *memPrefixManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func TestPrefixAddCommand(t *testing.T) {
	mgr := newMemPrefixManager()
	mgr.data["snap"] = map[string]string{"FOO": "bar", "BAZ": "qux"}

	res, err := prefix.RunAdd(mgr, "snap", "snap", "CI_")
	if err != nil {
		t.Fatalf("RunAdd: %v", err)
	}
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added keys, got %d", len(res.Added))
	}
	if _, ok := mgr.data["snap"]["CI_FOO"]; !ok {
		t.Error("CI_FOO should exist after add")
	}
}

func TestPrefixRemoveCommand(t *testing.T) {
	mgr := newMemPrefixManager()
	mgr.data["snap"] = map[string]string{"CI_FOO": "bar", "CI_BAZ": "qux", "KEEP": "yes"}

	res, err := prefix.RunRemove(mgr, "snap", "clean", "CI_")
	if err != nil {
		t.Fatalf("RunRemove: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPrefixFormatOutput(t *testing.T) {
	r := prefix.Result{Added: []string{"A"}, Skipped: []string{"B", "C"}}
	out := prefix.Format(r)
	buf := bytes.NewBufferString(out)
	if buf.Len() == 0 {
		t.Error("format output must not be empty")
	}
}
