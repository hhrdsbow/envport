package transform_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envport/internal/transform"
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

func TestRunPrefixAdd(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"HOST": "localhost", "PORT": "8080"})

	ops := []transform.Op{{Kind: "prefix-add", Value: "APP_"}}
	r, out, err := transform.Run(m, "dev", ops, false)
	if err != nil {
		t.Fatal(err)
	}
	if r.Changed != 2 || r.Total != 2 {
		t.Errorf("expected 2/2 changed, got %d/%d", r.Changed, r.Total)
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected key APP_HOST")
	}
}

func TestRunPrefixRemove(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"APP_HOST": "localhost"})

	ops := []transform.Op{{Kind: "prefix-remove", Value: "APP_"}}
	_, out, err := transform.Run(m, "dev", ops, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["HOST"]; !ok {
		t.Error("expected key HOST after prefix removal")
	}
}

func TestRunUppercase(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"key": "value"})

	ops := []transform.Op{{Kind: "uppercase"}}
	_, out, err := transform.Run(m, "dev", ops, false)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "VALUE" {
		t.Errorf("expected KEY=VALUE, got %v", out)
	}
}

func TestRunDryRunDoesNotSave(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"HOST": "localhost"})

	ops := []transform.Op{{Kind: "prefix-add", Value: "X_"}}
	_, _, err := transform.Run(m, "dev", ops, true)
	if err != nil {
		t.Fatal(err)
	}
	// original should be unchanged
	loaded, _ := m.Load("dev")
	if _, ok := loaded["HOST"]; !ok {
		t.Error("dry-run must not mutate the store")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, _, err := transform.Run(m, "ghost", nil, false)
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	r := transform.Result{Changed: 3, Total: 5}
	got := transform.Format(r, false)
	if got != "transformed 3/5 variable(s)" {
		t.Errorf("unexpected format: %q", got)
	}
	got = transform.Format(r, true)
	if got != "would transform 3/5 variable(s)" {
		t.Errorf("unexpected dry-run format: %q", got)
	}
}
