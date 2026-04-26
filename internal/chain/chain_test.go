package chain_test

import (
	"errors"
	"testing"

	"envport/internal/chain"
)

type memManager struct {
	data map[string]map[string]string
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

func newManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func TestRunMergesInOrder(t *testing.T) {
	m := newManager()
	m.data["base"] = map[string]string{"A": "1", "B": "2"}
	m.data["override"] = map[string]string{"B": "99", "C": "3"}

	r, err := chain.Run(m, []string{"base", "override"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Vars["A"] != "1" || r.Vars["B"] != "99" || r.Vars["C"] != "3" {
		t.Errorf("unexpected vars: %v", r.Vars)
	}
	if len(r.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(r.Applied))
	}
}

func TestRunMissingErrorsWithoutSkip(t *testing.T) {
	m := newManager()
	m.data["base"] = map[string]string{"X": "1"}

	_, err := chain.Run(m, []string{"base", "missing"}, false)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunSkipsMissing(t *testing.T) {
	m := newManager()
	m.data["base"] = map[string]string{"X": "1"}

	r, err := chain.Run(m, []string{"base", "ghost"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "ghost" {
		t.Errorf("expected ghost in skipped, got %v", r.Skipped)
	}
	if r.Vars["X"] != "1" {
		t.Errorf("expected X=1, got %v", r.Vars["X"])
	}
}

func TestRunEmptyNamesError(t *testing.T) {
	m := newManager()
	_, err := chain.Run(m, []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty names")
	}
}

func TestFormat(t *testing.T) {
	r := chain.Result{
		Vars:    map[string]string{"A": "1", "B": "2"},
		Applied: []string{"base", "prod"},
		Skipped: []string{"ghost"},
	}
	out := chain.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
