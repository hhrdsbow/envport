package squash_test

import (
	"errors"
	"testing"

	"envport/internal/squash"
)

// --- in-memory test doubles ---

type memSnap struct {
	name string
	vars map[string]string
}

func (s *memSnap) Vars() map[string]string { return s.vars }
func (s *memSnap) Name() string            { return s.name }

type memManager struct {
	snaps map[string]*memSnap
}

func newManager() *memManager {
	return &memManager{snaps: make(map[string]*memSnap)}
}

func (m *memManager) Load(name string) (squash.Snapshot, error) {
	s, ok := m.snaps[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.snaps[name] = &memSnap{name: name, vars: vars}
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.snaps[name] = &memSnap{name: name, vars: vars}
}

// --- tests ---

func TestRunMergesNoConflicts(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seed(m, "b", map[string]string{"BAZ": "3"})

	r, err := squash.Run(m, []string{"a", "b"}, "out")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.KeyCount != 3 {
		t.Errorf("expected 3 keys, got %d", r.KeyCount)
	}
}

func TestRunLaterSourceWins(t *testing.T) {
	m := newManager()
	seed(m, "base", map[string]string{"KEY": "old", "X": "1"})
	seed(m, "override", map[string]string{"KEY": "new"})

	_, err := squash.Run(m, []string{"base", "override"}, "merged")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, _ := m.Load("merged")
	if out.Vars()["KEY"] != "new" {
		t.Errorf("expected 'new', got %q", out.Vars()["KEY"])
	}
	if out.Vars()["X"] != "1" {
		t.Errorf("expected X=1, got %q", out.Vars()["X"])
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"FOO": "1"})

	_, err := squash.Run(m, []string{"a", "missing"}, "out")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRunTooFewSources(t *testing.T) {
	m := newManager()
	_, err := squash.Run(m, []string{"only"}, "out")
	if err == nil {
		t.Fatal("expected error for fewer than 2 sources")
	}
}

func TestFormat(t *testing.T) {
	r := squash.Result{Dest: "merged", Sources: []string{"a", "b", "c"}, KeyCount: 7}
	got := squash.Format(r)
	if got == "" {
		t.Error("Format returned empty string")
	}
}
