package summarize_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/envport/internal/summarize"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	data map[string]map[string]string
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func newManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = vars
}

func TestRunBasicCounts(t *testing.T) {
	mgr := newManager()
	seed(mgr, "prod", map[string]string{
		"APP_HOST":     "localhost",
		"APP_PORT":     "8080",
		"DB_PASSWORD":  "s3cr3t",
		"DB_HOST":      "",
		"API_TOKEN":    "tok123",
	})

	r, err := summarize.Run(mgr, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Total != 5 {
		t.Errorf("Total: want 5, got %d", r.Total)
	}
	if r.Empty != 1 {
		t.Errorf("Empty: want 1, got %d", r.Empty)
	}
	if r.Sensitive != 2 {
		t.Errorf("Sensitive: want 2, got %d", r.Sensitive)
	}
}

func TestRunPrefixGrouping(t *testing.T) {
	mgr := newManager()
	seed(mgr, "dev", map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "3000",
		"DB_HOST":  "127.0.0.1",
	})

	r, err := summarize.Run(mgr, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Prefixes["APP"] != 2 {
		t.Errorf("APP prefix count: want 2, got %d", r.Prefixes["APP"])
	}
	if r.Prefixes["DB"] != 1 {
		t.Errorf("DB prefix count: want 1, got %d", r.Prefixes["DB"])
	}
}

func TestRunMissing(t *testing.T) {
	mgr := newManager()
	_, err := summarize.Run(mgr, "ghost")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := summarize.Result{
		Name:      "staging",
		Total:     4,
		Empty:     1,
		Sensitive: 2,
		Prefixes:  map[string]int{"APP": 3, "DB": 1},
	}
	out := summarize.Format(r)
	for _, want := range []string{"staging", "Total", "4", "Empty", "1", "Sensitive", "2", "APP", "DB"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output missing %q", want)
		}
	}
}
