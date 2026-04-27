package split_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envport/internal/split"
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
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	m.data[name] = copy
	return nil
}

func seed(m *memManager) {
	m.data["src"] = map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_PORT": "8080",
		"APP_ENV":  "prod",
	}
}

func TestRunSplitsIntoTargets(t *testing.T) {
	m := newMemManager()
	seed(m)

	targets := map[string][]string{
		"db":  {"DB_HOST", "DB_PORT"},
		"app": {"APP_PORT", "APP_ENV"},
	}

	results, err := split.Run(m, "src", targets, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	db := m.data["db"]
	if db["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", db["DB_HOST"])
	}
}

func TestRunRemainder(t *testing.T) {
	m := newMemManager()
	seed(m)

	targets := map[string][]string{
		"db": {"DB_HOST", "DB_PORT"},
	}

	_, err := split.Run(m, "src", targets, "rest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rest := m.data["rest"]
	if _, ok := rest["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in remainder")
	}
	if _, ok := rest["DB_HOST"]; ok {
		t.Error("DB_HOST should not appear in remainder")
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := split.Run(m, "missing", map[string][]string{"a": {"X"}}, "")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestFormat(t *testing.T) {
	results := []split.Result{
		{Name: "db", Keys: []string{"DB_HOST", "DB_PORT"}},
		{Name: "app", Keys: []string{"APP_PORT"}},
	}
	out := split.Format(results)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormatEmpty(t *testing.T) {
	out := split.Format(nil)
	if out != "no snapshots created\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
