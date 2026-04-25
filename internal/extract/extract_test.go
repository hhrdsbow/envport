package extract_test

import (
	"errors"
	"testing"

	"envport/internal/extract"
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
		return nil, errors.New("not found: " + name)
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

func TestRunExtractsKeys(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "prod", map[string]string{"DB_HOST": "db.prod", "API_KEY": "secret", "PORT": "5432"})

	got, err := extract.Run(mgr, extract.Options{Src: "prod", Dst: "prod-db", Keys: []string{"DB_HOST", "PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "db.prod" {
		t.Errorf("DB_HOST: got %q", got["DB_HOST"])
	}
	if got["PORT"] != "5432" {
		t.Errorf("PORT: got %q", got["PORT"])
	}
	if _, ok := got["API_KEY"]; ok {
		t.Error("API_KEY should not be in extracted snapshot")
	}
}

func TestRunSavesDestination(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "src", map[string]string{"X": "1", "Y": "2"})

	_, err := extract.Run(mgr, extract.Options{Src: "src", Dst: "dst", Keys: []string{"X"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dst, _ := mgr.Load("dst")
	if dst["X"] != "1" {
		t.Errorf("saved value wrong: %q", dst["X"])
	}
}

func TestRunNoKeys(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "src", map[string]string{"A": "1"})

	_, err := extract.Run(mgr, extract.Options{Src: "src", Dst: "dst", Keys: nil})
	if !errors.Is(err, extract.ErrNoKeys) {
		t.Fatalf("expected ErrNoKeys, got %v", err)
	}
}

func TestRunMissingKey(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "src", map[string]string{"A": "1"})

	_, err := extract.Run(mgr, extract.Options{Src: "src", Dst: "dst", Keys: []string{"MISSING"}})
	if !errors.Is(err, extract.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestRunMissingSource(t *testing.T) {
	mgr := newMemManager()

	_, err := extract.Run(mgr, extract.Options{Src: "ghost", Dst: "dst", Keys: []string{"A"}})
	if err == nil {
		t.Fatal("expected error for missing source snapshot")
	}
}
