package sample_test

import (
	"errors"
	"strings"
	"testing"

	"envport/internal/sample"
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
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = vars
}

func TestRunSamplesNKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"})

	r, err := sample.Run(m, "src", "dst", 3, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
	dst, _ := m.Load("dst")
	if len(dst) != 3 {
		t.Fatalf("expected dest to have 3 entries, got %d", len(dst))
	}
}

func TestRunNGreaterThanLen(t *testing.T) {
	m := newMemManager()
	seed(m, "src", map[string]string{"X": "1", "Y": "2"})

	r, err := sample.Run(m, "src", "dst", 100, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
}

func TestRunMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := sample.Run(m, "missing", "dst", 1, 0)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestFormat(t *testing.T) {
	r := sample.Result{Source: "src", Dest: "dst", Keys: []string{"A", "B"}}
	got := sample.Format(r)
	if !strings.Contains(got, "2") || !strings.Contains(got, "src") || !strings.Contains(got, "dst") {
		t.Fatalf("unexpected format output: %q", got)
	}
}
