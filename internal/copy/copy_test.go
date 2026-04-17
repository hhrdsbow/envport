package copy_test

import (
	"errors"
	"testing"

	copy_ "github.com/envport/envport/internal/copy"
)

// memManager is an in-memory implementation of copy.Manager.
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
	return v, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func (m *memManager) List() ([]string, error) {
	names := make([]string, 0, len(m.data))
	for k := range m.data {
		names = append(names, k)
	}
	return names, nil
}

func TestRunCopiesSnapshot(t *testing.T) {
	mgr := newMemManager()
	mgr.data["src"] = map[string]string{"FOO": "bar", "BAZ": "qux"}

	if err := copy_.Run(mgr, "src", "dst", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := mgr.Load("dst")
	if err != nil {
		t.Fatalf("dst not saved: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", got)
	}
}

func TestRunMissingSource(t *testing.T) {
	mgr := newMemManager()
	err := copy_.Run(mgr, "missing", "dst", false)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRunNoOverwrite(t *testing.T) {
	mgr := newMemManager()
	mgr.data["src"] = map[string]string{"A": "1"}
	mgr.data["dst"] = map[string]string{"B": "2"}

	err := copy_.Run(mgr, "src", "dst", false)
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestRunOverwrite(t *testing.T) {
	mgr := newMemManager()
	mgr.data["src"] = map[string]string{"A": "1"}
	mgr.data["dst"] = map[string]string{"B": "2"}

	if err := copy_.Run(mgr, "src", "dst", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := mgr.Load("dst")
	if got["A"] != "1" {
		t.Errorf("expected overwritten value, got %v", got)
	}
}

func TestRunSameSourceDest(t *testing.T) {
	mgr := newMemManager()
	err := copy_.Run(mgr, "x", "x", false)
	if err == nil {
		t.Fatal("expected error when src == dst")
	}
}
