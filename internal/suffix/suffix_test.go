package suffix_test

import (
	"errors"
	"testing"

	"envport/internal/suffix"
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
		return nil, errors.New("snapshot not found: " + name)
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
	m := newMemManager()
	seed(m, "dev", map[string]string{"URL": "http://localhost", "DB": "mydb"})

	res, err := suffix.RunAdd(m, "dev", "_v2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Modified) != 2 {
		t.Errorf("expected 2 modified, got %d", len(res.Modified))
	}
	vars, _ := m.Load("dev")
	if vars["URL"] != "http://localhost_v2" {
		t.Errorf("URL = %q, want %q", vars["URL"], "http://localhost_v2")
	}
}

func TestRunAddSelectedKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"URL": "http://localhost", "DB": "mydb"})

	res, err := suffix.RunAdd(m, "dev", "_test", []string{"DB"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Modified) != 1 || res.Modified[0] != "DB" {
		t.Errorf("expected DB modified, got %v", res.Modified)
	}
	vars, _ := m.Load("dev")
	if vars["URL"] != "http://localhost" {
		t.Errorf("URL should be unchanged, got %q", vars["URL"])
	}
}

func TestRunRemoveSuffix(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"URL": "http://localhost_v2", "DB": "mydb"})

	res, err := suffix.RunRemove(m, "prod", "_v2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Modified) != 1 {
		t.Errorf("expected 1 modified, got %d", len(res.Modified))
	}
	vars, _ := m.Load("prod")
	if vars["URL"] != "http://localhost" {
		t.Errorf("URL = %q, want %q", vars["URL"], "http://localhost")
	}
}

func TestRunAddEmptySuffix(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "val"})
	_, err := suffix.RunAdd(m, "dev", "", nil)
	if err == nil {
		t.Error("expected error for empty suffix")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := suffix.RunAdd(m, "missing", "_x", nil)
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := suffix.Result{Modified: []string{"A", "B"}, Skipped: []string{"C"}}
	out := suffix.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
