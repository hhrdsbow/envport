package uppercase_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envport/internal/uppercase"
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
		return nil, errors.New("not found")
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

func TestRunUppercasesAllValues(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"APP_ENV": "production", "LOG_LEVEL": "debug"})

	r, err := uppercase.Run(m, "dev", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changed) != 2 {
		t.Fatalf("expected 2 changed keys, got %d", len(r.Changed))
	}

	vars, _ := m.Load("dev")
	if vars["APP_ENV"] != "PRODUCTION" {
		t.Errorf("APP_ENV: want PRODUCTION, got %s", vars["APP_ENV"])
	}
	if vars["LOG_LEVEL"] != "DEBUG" {
		t.Errorf("LOG_LEVEL: want DEBUG, got %s", vars["LOG_LEVEL"])
	}
}

func TestRunFilterKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"APP_ENV": "production", "LOG_LEVEL": "debug"})

	_, err := uppercase.Run(m, "dev", []string{"LOG_LEVEL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vars, _ := m.Load("dev")
	if vars["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %s", vars["APP_ENV"])
	}
	if vars["LOG_LEVEL"] != "DEBUG" {
		t.Errorf("LOG_LEVEL: want DEBUG, got %s", vars["LOG_LEVEL"])
	}
}

func TestRunAlreadyUppercase(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"MODE": "RELEASE"})

	r, err := uppercase.Run(m, "prod", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changed keys, got %v", r.Changed)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := uppercase.Run(m, "missing", nil)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := uppercase.Result{Name: "dev", Changed: []string{"FOO", "BAR"}}
	out := uppercase.Format(r)
	if out == "" {
		t.Error("Format returned empty string")
	}

	r2 := uppercase.Result{Name: "dev", Changed: nil}
	out2 := uppercase.Format(r2)
	if out2 == "" {
		t.Error("Format returned empty string for no changes")
	}
}
