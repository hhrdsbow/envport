package pivot_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/pivot"
)

// --- in-memory fakes ---

type memSnap struct{ vars map[string]string }

func (s *memSnap) Vars() map[string]string { return s.vars }

type memManager struct {
	snaps map[string]*memSnap
}

func newMemManager() *memManager {
	return &memManager{snaps: make(map[string]*memSnap)}
}

func (m *memManager) Load(name string) (pivot.Snapshot, error) {
	s, ok := m.snaps[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.snaps[name] = &memSnap{vars: vars}
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.snaps[name] = &memSnap{vars: vars}
}

// --- tests ---

func TestRunPivotsSnapshot(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"HOST": "localhost", "PORT": "8080"})

	res, err := pivot.Run(m, "src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Keys != 2 {
		t.Errorf("expected 2 keys, got %d", res.Keys)
	}

	dst := m.snaps["dst"]
	if dst.vars["localhost"] != "HOST" {
		t.Errorf("expected localhost→HOST, got %q", dst.vars["localhost"])
	}
	if dst.vars["8080"] != "PORT" {
		t.Errorf("expected 8080→PORT, got %q", dst.vars["8080"])
	}
}

func TestRunSkipsBlankValues(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"EMPTY": "", "KEY": "val"})

	res, err := pivot.Run(m, "src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Keys != 1 {
		t.Errorf("expected 1 key (blank skipped), got %d", res.Keys)
	}
	if _, ok := m.snaps["dst"].vars[""]; ok {
		t.Error("blank value should not appear as key in destination")
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := pivot.Run(m, "missing", "dst")
	if err == nil {
		t.Fatal("expected error for missing source snapshot")
	}
}
