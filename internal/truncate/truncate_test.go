package truncate_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/truncate"
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
	m.data[name] = copy
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	_ = m.Save(name, vars)
}

func TestRunTruncatesLongValues(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{
		"TOKEN": "abcdefghij",
		"HOST":  "localhost",
	})

	results, err := truncate.Run(m, "prod", truncate.Options{MaxLen: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "TOKEN" {
		t.Errorf("expected key TOKEN, got %s", results[0].Key)
	}
	if results[0].Truncated != "abcde..." {
		t.Errorf("unexpected truncated value: %q", results[0].Truncated)
	}

	vars, _ := m.Load("prod")
	if vars["TOKEN"] != "abcde..." {
		t.Errorf("saved value mismatch: %q", vars["TOKEN"])
	}
	if vars["HOST"] != "localhost" {
		t.Errorf("short value should be unchanged, got %q", vars["HOST"])
	}
}

func TestRunNoChangesNeeded(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "hi", "B": "ok"})

	results, err := truncate.Run(m, "dev", truncate.Options{MaxLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestRunFilterKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "stage", map[string]string{
		"SECRET": "verylongsecretvalue",
		"URL":    "verylongurl",
	})

	results, err := truncate.Run(m, "stage", truncate.Options{
		MaxLen: 5,
		Keys:   []string{"SECRET"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "SECRET" {
		t.Errorf("expected only SECRET to be truncated")
	}
	vars, _ := m.Load("stage")
	if vars["URL"] != "verylongurl" {
		t.Errorf("URL should be unchanged, got %q", vars["URL"])
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := truncate.Run(m, "missing", truncate.Options{MaxLen: 5})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunInvalidMaxLen(t *testing.T) {
	m := newMemManager()
	seed(m, "x", map[string]string{"K": "v"})
	_, err := truncate.Run(m, "x", truncate.Options{MaxLen: 0})
	if err == nil {
		t.Fatal("expected error for MaxLen=0")
	}
}

func TestRunCustomEllipsis(t *testing.T) {
	m := newMemManager()
	seed(m, "e", map[string]string{"VAL": "hello world"})

	results, err := truncate.Run(m, "e", truncate.Options{MaxLen: 5, Ellipsis: "~~"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Truncated != "hello~~" {
		t.Errorf("unexpected truncated value: %q", results[0].Truncated)
	}
}
