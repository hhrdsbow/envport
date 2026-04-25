package flatten_test

import (
	"errors"
	"testing"

	"envport/internal/flatten"
)

// --- test doubles ---

type memSnap struct {
	name string
	vars map[string]string
}

func (s *memSnap) Name() string            { return s.name }
func (s *memSnap) Vars() map[string]string { return s.vars }

type memManager struct {
	snaps map[string]*memSnap
}

func (m *memManager) Load(name string) (flatten.Snapshot, error) {
	s, ok := m.snaps[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return s, nil
}

func newManager(snaps ...*memSnap) *memManager {
	m := &memManager{snaps: make(map[string]*memSnap)}
	for _, s := range snaps {
		m.snaps[s.name] = s
	}
	return m
}

// --- tests ---

func TestRunNoNames(t *testing.T) {
	mgr := newManager()
	_, err := flatten.Run(mgr, nil)
	if err == nil {
		t.Fatal("expected error for empty names")
	}
}

func TestRunSingleSnapshot(t *testing.T) {
	mgr := newManager(&memSnap{name: "base", vars: map[string]string{"A": "1", "B": "2"}})
	res, err := flatten.Run(mgr, []string{"base"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "1" || res.Vars["B"] != "2" {
		t.Fatalf("unexpected vars: %v", res.Vars)
	}
	if res.Sources["A"] != "base" {
		t.Fatalf("expected source 'base', got %q", res.Sources["A"])
	}
}

func TestRunLaterOverrides(t *testing.T) {
	mgr := newManager(
		&memSnap{name: "base", vars: map[string]string{"A": "base_val", "B": "shared"}},
		&memSnap{name: "override", vars: map[string]string{"A": "new_val", "C": "extra"}},
	)
	res, err := flatten.Run(mgr, []string{"base", "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "new_val" {
		t.Errorf("expected A=new_val, got %q", res.Vars["A"])
	}
	if res.Vars["B"] != "shared" {
		t.Errorf("expected B=shared, got %q", res.Vars["B"])
	}
	if res.Sources["A"] != "override" {
		t.Errorf("expected source 'override' for A, got %q", res.Sources["A"])
	}
	if res.Sources["B"] != "base" {
		t.Errorf("expected source 'base' for B, got %q", res.Sources["B"])
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	mgr := newManager()
	_, err := flatten.Run(mgr, []string{"ghost"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
