package pivot_test

import (
	"errors"
	"testing"

	"github.com/your-org/envport/internal/pivot"
)

// --- in-memory manager ---

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
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	cp := make(map[string]string, len(vars))
	for k, v := range vars {
		cp[k] = v
	}
	m.store[name] = cp
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	_ = m.Save(name, vars)
}

// --- tests ---

func TestRunPivotsKeyValues(t *testing.T) {
	m := newMemManager()
	seed(m, "p", map[string]string{"HOST": "localhost", "PORT": "8080"})

	got, err := pivot.Run(m, "p", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["localhost"] != "HOST" {
		t.Errorf("expected localhost→HOST, got %q", got["localhost"])
	}
	if got["8080"] != "PORT" {
		t.Errorf("expected 8080→PORT, got %q", got["8080"])
	}
}

func TestRunSavesToDst(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"A": "1"})

	_, err := pivot.Run(m, "src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, err := m.Load("dst")
	if err != nil {
		t.Fatalf("dst not saved: %v", err)
	}
	if v["1"] != "A" {
		t.Errorf("expected 1→A in dst, got %q", v["1"])
	}
}

func TestRunSkipsBlankValues(t *testing.T) {
	m := newMemManager()
	seed(m, "p", map[string]string{"EMPTY": "", "KEY": "val"})

	got, err := pivot.Run(m, "p", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got[""]; ok {
		t.Error("blank value should not become a key")
	}
	if got["val"] != "KEY" {
		t.Errorf("expected val→KEY, got %q", got["val"])
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := pivot.Run(m, "missing", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunEmptySourceName(t *testing.T) {
	m := newMemManager()
	_, err := pivot.Run(m, "", "dst")
	if err == nil {
		t.Fatal("expected error for empty source name")
	}
}
