package lowercase_test

import (
	"errors"
	"testing"

	"envport/internal/lowercase"
)

// --- in-memory manager ---

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

// --- tests ---

func TestRunLowercasesAllValues(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "dev", map[string]string{"FOO": "Hello", "BAR": "WORLD"})

	r, err := lowercase.Run(mgr, "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Changed) != 2 {
		t.Fatalf("expected 2 changed, got %d", len(r.Changed))
	}
	vars, _ := mgr.Load("dev")
	if vars["FOO"] != "hello" || vars["BAR"] != "world" {
		t.Errorf("unexpected values: %v", vars)
	}
}

func TestRunNoChange(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "dev", map[string]string{"KEY": "already"})

	r, err := lowercase.Run(mgr, "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Changed) != 0 {
		t.Fatalf("expected 0 changed, got %d", len(r.Changed))
	}
}

func TestRunFilterKeys(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "dev", map[string]string{"A": "UPPER", "B": "UPPER"})

	_, err := lowercase.Run(mgr, "dev", []string{"A"})
	if err != nil {
		t.Fatal(err)
	}
	vars, _ := mgr.Load("dev")
	if vars["A"] != "upper" {
		t.Errorf("A should be lowercased, got %q", vars["A"])
	}
	if vars["B"] != "UPPER" {
		t.Errorf("B should be unchanged, got %q", vars["B"])
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	mgr := newMemManager()
	_, err := lowercase.Run(mgr, "missing", nil)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := lowercase.Result{Name: "dev", Changed: []string{"FOO", "BAR"}}
	out := lowercase.Format(r)
	if out == "" {
		t.Fatal("expected non-empty format output")
	}
	if !containsStr(out, "dev") || !containsStr(out, "2") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func TestFormatNoChange(t *testing.T) {
	r := lowercase.Result{Name: "prod", Changed: nil}
	out := lowercase.Format(r)
	if !containsStr(out, "no values changed") {
		t.Errorf("unexpected output: %s", out)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
