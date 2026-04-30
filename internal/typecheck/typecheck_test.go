package typecheck

import (
	"errors"
	"strings"
	"testing"
)

// memManager is an in-memory Manager for tests.
type memManager struct {
	data map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func (m *memManager) save(name string, vars map[string]string) {
	m.data[name] = vars
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func seed(m *memManager) {
	m.save("prod", map[string]string{
		"PORT":     "8080",
		"TIMEOUT":  "30.5",
		"DEBUG":    "true",
		"BASE_URL": "https://example.com",
		"BAD_INT":  "notanumber",
		"BAD_URL":  "not-a-url",
	})
}

func TestRunNoViolations(t *testing.T) {
	m := newMemManager()
	seed(m)
	typeMap := map[string]Type{
		"PORT":     TypeInt,
		"TIMEOUT":  TypeFloat,
		"DEBUG":    TypeBool,
		"BASE_URL": TypeURL,
	}
	res, err := Run(m, "prod", typeMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(res.Violations))
	}
}

func TestRunDetectsViolations(t *testing.T) {
	m := newMemManager()
	seed(m)
	typeMap := map[string]Type{
		"BAD_INT": TypeInt,
		"BAD_URL": TypeURL,
	}
	res, err := Run(m, "prod", typeMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(res.Violations))
	}
}

func TestRunSkipsMissingKeys(t *testing.T) {
	m := newMemManager()
	m.save("empty", map[string]string{})
	typeMap := map[string]Type{"PORT": TypeInt}
	res, err := Run(m, "empty", typeMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("expected no violations for missing key")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := Run(m, "missing", map[string]Type{"X": TypeInt})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormatClean(t *testing.T) {
	out := Format(Result{})
	if !strings.Contains(out, "all values") {
		t.Errorf("unexpected clean output: %q", out)
	}
}

func TestFormatWithViolations(t *testing.T) {
	res := Result{Violations: []Violation{
		{Key: "PORT", Value: "abc", Expected: TypeInt, Reason: `"abc" is not a valid integer`},
	}}
	out := Format(res)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output: %q", out)
	}
	if !strings.Contains(out, "1 type violation") {
		t.Errorf("expected violation count in output: %q", out)
	}
}
