package supersede_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/supersede"
)

// memManager is a simple in-memory Manager for testing.
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
	cp := make(map[string]string, len(vars))
	for k, v := range vars {
		cp[k] = v
	}
	m.data[name] = cp
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	_ = m.Save(name, vars)
}

func TestRunWithExplicitKeys(t *testing.T) {
	m := newMemManager()
	seed(m, "prod", map[string]string{"HOST": "old.host", "PORT": "8080"})

	res, err := supersede.Run(m, "prod", supersede.Options{
		Keys: map[string]string{"HOST": "new.host"},
		Dest:  "prod-patched",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Dest != "prod-patched" {
		t.Errorf("dest = %q, want prod-patched", res.Dest)
	}
	if res.Applied["HOST"] != "new.host" {
		t.Errorf("applied HOST = %q, want new.host", res.Applied["HOST"])
	}
	loaded, _ := m.Load("prod-patched")
	if loaded["PORT"] != "8080" {
		t.Errorf("PORT should be preserved, got %q", loaded["PORT"])
	}
}

func TestRunOverwritesBase(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"DB": "localhost"})

	_, err := supersede.Run(m, "dev", supersede.Options{
		Keys: map[string]string{"DB": "remotedb"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, _ := m.Load("dev")
	if loaded["DB"] != "remotedb" {
		t.Errorf("DB = %q, want remotedb", loaded["DB"])
	}
}

func TestRunFromOverrideProfile(t *testing.T) {
	m := newMemManager()
	seed(m, "base", map[string]string{"A": "1", "B": "2"})
	seed(m, "overrides", map[string]string{"B": "99", "C": "3"})

	res, err := supersede.Run(m, "base", supersede.Options{
		OverrideProfile: "overrides",
		Dest:            "merged",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 2 {
		t.Errorf("applied count = %d, want 2", len(res.Applied))
	}
	loaded, _ := m.Load("merged")
	if loaded["A"] != "1" || loaded["B"] != "99" || loaded["C"] != "3" {
		t.Errorf("unexpected merged vars: %v", loaded)
	}
}

func TestRunMissingBase(t *testing.T) {
	m := newMemManager()
	_, err := supersede.Run(m, "ghost", supersede.Options{
		Keys: map[string]string{"X": "1"},
	})
	if err == nil {
		t.Fatal("expected error for missing base, got nil")
	}
}

func TestRunEmptyBaseName(t *testing.T) {
	m := newMemManager()
	_, err := supersede.Run(m, "", supersede.Options{})
	if err == nil {
		t.Fatal("expected error for empty base name")
	}
}
