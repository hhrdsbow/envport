package select

import (
	"errors"
	"testing"
)

// memManager is an in-memory Manager for tests.
type memManager struct {
	data map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = vars
}

func TestRunSelectsExactKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"A": "1", "B": "2", "C": "3"})

	r, err := Run(m, "src", "dst", []string{"A", "C"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Kept) != 2 || len(r.Dropped) != 1 {
		t.Fatalf("expected kept=2 dropped=1, got kept=%d dropped=%d", len(r.Kept), len(r.Dropped))
	}
	dst := m.data["dst"]
	if dst["A"] != "1" || dst["C"] != "3" {
		t.Errorf("unexpected dst contents: %v", dst)
	}
	if _, ok := dst["B"]; ok {
		t.Error("B should have been dropped")
	}
}

func TestRunSelectsGlobPattern(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"DB_HOST": "h", "DB_PORT": "5432", "APP_ENV": "prod"})

	r, err := Run(m, "src", "dst", []string{"DB_*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Kept) != 2 {
		t.Fatalf("expected 2 kept, got %d", len(r.Kept))
	}
}

func TestRunNoPatterns(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"A": "1"})

	_, err := Run(m, "src", "dst", nil)
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := Run(m, "missing", "dst", []string{"A"})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestFormatOutput(t *testing.T) {
	r := Result{Kept: []string{"A", "B"}, Dropped: []string{"C"}}
	out := Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
