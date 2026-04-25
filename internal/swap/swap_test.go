package swap_test

import (
	"errors"
	"testing"

	"envport/internal/swap"
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

func TestRunSwapsCommonKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"DB_HOST": "dev-db", "PORT": "3000"})
	seed(m, "prod", map[string]string{"DB_HOST": "prod-db", "PORT": "8080"})

	r, err := swap.Run(m, "dev", "prod", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 swapped keys, got %d", len(r.Keys))
	}

	dev, _ := m.Load("dev")
	prod, _ := m.Load("prod")

	if dev["DB_HOST"] != "prod-db" {
		t.Errorf("dev DB_HOST: want prod-db, got %s", dev["DB_HOST"])
	}
	if prod["DB_HOST"] != "dev-db" {
		t.Errorf("prod DB_HOST: want dev-db, got %s", prod["DB_HOST"])
	}
	if dev["PORT"] != "8080" {
		t.Errorf("dev PORT: want 8080, got %s", dev["PORT"])
	}
	if prod["PORT"] != "3000" {
		t.Errorf("prod PORT: want 3000, got %s", prod["PORT"])
	}
}

func TestRunSwapsExplicitKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "a", map[string]string{"FOO": "foo-a", "BAR": "bar-a"})
	seed(m, "b", map[string]string{"FOO": "foo-b", "BAR": "bar-b"})

	_, err := swap.Run(m, "a", "b", []string{"FOO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, _ := m.Load("a")
	b, _ := m.Load("b")

	if a["FOO"] != "foo-b" {
		t.Errorf("a FOO: want foo-b, got %s", a["FOO"])
	}
	if b["FOO"] != "foo-a" {
		t.Errorf("b FOO: want foo-a, got %s", b["FOO"])
	}
	// BAR should be untouched
	if a["BAR"] != "bar-a" {
		t.Errorf("a BAR should be unchanged, got %s", a["BAR"])
	}
}

func TestRunNoCommonKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "x", map[string]string{"ALPHA": "1"})
	seed(m, "y", map[string]string{"BETA": "2"})

	r, err := swap.Run(m, "x", "y", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestRunMissingSrc(t *testing.T) {
	m := newMemManager()
	seed(m, "dst", map[string]string{"K": "v"})

	_, err := swap.Run(m, "missing", "dst", nil)
	if err == nil {
		t.Fatal("expected error for missing src")
	}
}

func TestFormatOutput(t *testing.T) {
	r := swap.Result{SrcName: "dev", DstName: "prod", Keys: []string{"DB_HOST"}}
	out := swap.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
