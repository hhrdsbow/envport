package patch_test

import (
	"errors"
	"sort"
	"testing"

	"envport/internal/patch"
)

// memManager is an in-memory Manager for testing.
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
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	m.store[name] = copy
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	_ = m.Save(name, vars)
}

func TestRunAddsNewKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})

	r, err := patch.Run(m, "dev", map[string]string{"B": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Added) != 1 || r.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", r.Added)
	}
	if len(r.Updated) != 0 {
		t.Errorf("expected no updates, got %v", r.Updated)
	}
}

func TestRunUpdatesExistingKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1", "B": "2"})

	r, err := patch.Run(m, "dev", map[string]string{"A": "99"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Updated) != 1 || r.Updated[0] != "A" {
		t.Errorf("expected Updated=[A], got %v", r.Updated)
	}
	if len(r.Added) != 0 {
		t.Errorf("expected no adds, got %v", r.Added)
	}
}

func TestRunPreservesUntouchedKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1", "B": "2"})

	_, err := patch.Run(m, "dev", map[string]string{"A": "99"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vars, _ := m.Load("dev")
	if vars["B"] != "2" {
		t.Errorf("expected B=2 to be preserved, got %q", vars["B"])
	}
}

func TestRunMixedAddAndUpdate(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"HOST": "localhost"})

	r, err := patch.Run(m, "prod", map[string]string{"HOST": "prod.example.com", "PORT": "443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sort.Strings(r.Added)
	sort.Strings(r.Updated)

	if len(r.Added) != 1 || r.Added[0] != "PORT" {
		t.Errorf("expected Added=[PORT], got %v", r.Added)
	}
	if len(r.Updated) != 1 || r.Updated[0] != "HOST" {
		t.Errorf("expected Updated=[HOST], got %v", r.Updated)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := patch.Run(m, "missing", map[string]string{"X": "1"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunEmptyName(t *testing.T) {
	m := newMemManager()
	_, err := patch.Run(m, "", map[string]string{"X": "1"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRunEmptyUpdates(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})
	_, err := patch.Run(m, "dev", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty updates")
	}
}

func TestFormat(t *testing.T) {
	r := patch.Result{Name: "dev", Added: []string{"X"}, Updated: []string{"Y", "Z"}}
	got := patch.Format(r)
	if got == "" {
		t.Error("expected non-empty format output")
	}
}
