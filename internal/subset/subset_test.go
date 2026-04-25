package subset_test

import (
	"errors"
	"testing"

	"envport/internal/subset"
)

// memManager is an in-memory Manager implementation for testing.
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

func TestRunExtractsSubset(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"DB": "pg", "PORT": "5432", "SECRET": "x"})

	r, err := subset.Run(m, "prod", "prod-db", []string{"DB", "PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Extracted) != 2 {
		t.Fatalf("expected 2 extracted, got %d", len(r.Extracted))
	}
	if r.Extracted["DB"] != "pg" || r.Extracted["PORT"] != "5432" {
		t.Errorf("unexpected extracted values: %v", r.Extracted)
	}
	if _, ok := m.data["prod-db"]; !ok {
		t.Error("destination snapshot was not saved")
	}
}

func TestRunMissingKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"DB": "pg"})

	r, err := subset.Run(m, "prod", "out", []string{"DB", "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Missing) != 1 || r.Missing[0] != "MISSING" {
		t.Errorf("expected one missing key, got %v", r.Missing)
	}
}

func TestRunAllKeysMissing(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"DB": "pg"})

	_, err := subset.Run(m, "prod", "out", []string{"NOPE"})
	if err == nil {
		t.Fatal("expected error when all keys missing")
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := subset.Run(m, "ghost", "out", []string{"KEY"})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRunNoKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"A": "1"})
	_, err := subset.Run(m, "prod", "out", nil)
	if err == nil {
		t.Fatal("expected error when no keys provided")
	}
}

func TestFormat(t *testing.T) {
	r := subset.Result{
		Source:      "src",
		Destination: "dst",
		Extracted:   map[string]string{"A": "1"},
		Missing:     []string{"B"},
	}
	out := subset.Format(r)
	if out == "" {
		t.Error("Format returned empty string")
	}
}
