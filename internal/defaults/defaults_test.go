package defaults_test

import (
	"testing"

	"envport/internal/defaults"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string]map[string]string
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]map[string]string)}
}

func (s *memStore) Get(profile string) (map[string]string, error) {
	v, ok := s.data[profile]
	if !ok {
		return nil, defaults.ErrNotFound
	}
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (s *memStore) Set(profile string, vals map[string]string) error {
	s.data[profile] = vals
	return nil
}

func (s *memStore) Delete(profile string) error {
	delete(s.data, profile)
	return nil
}

func (s *memStore) List() ([]string, error) {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestSetAndGet(t *testing.T) {
	m := defaults.New(newMemStore())
	if err := m.Set("dev", map[string]string{"FOO": "bar"}); err != nil {
		t.Fatal(err)
	}
	got, err := m.Get("dev")
	if err != nil {
		t.Fatal(err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("expected bar, got %q", got["FOO"])
	}
}

func TestSetMerges(t *testing.T) {
	m := defaults.New(newMemStore())
	_ = m.Set("dev", map[string]string{"A": "1"})
	_ = m.Set("dev", map[string]string{"B": "2"})
	got, _ := m.Get("dev")
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("merge failed: %v", got)
	}
}

func TestRemoveKeys(t *testing.T) {
	m := defaults.New(newMemStore())
	_ = m.Set("dev", map[string]string{"A": "1", "B": "2"})
	if err := m.Remove("dev", []string{"A"}); err != nil {
		t.Fatal(err)
	}
	got, _ := m.Get("dev")
	if _, ok := got["A"]; ok {
		t.Error("expected A to be removed")
	}
	if got["B"] != "2" {
		t.Error("B should remain")
	}
}

func TestApplyDoesNotOverwrite(t *testing.T) {
	m := defaults.New(newMemStore())
	_ = m.Set("dev", map[string]string{"FOO": "default", "BAR": "base"})
	vars := map[string]string{"FOO": "override"}
	out, err := m.Apply("dev", vars)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "override" {
		t.Errorf("expected override, got %q", out["FOO"])
	}
	if out["BAR"] != "base" {
		t.Errorf("expected base, got %q", out["BAR"])
	}
}

func TestApplyMissingProfile(t *testing.T) {
	m := defaults.New(newMemStore())
	vars := map[string]string{"X": "1"}
	out, err := m.Apply("nonexistent", vars)
	if err != nil {
		t.Fatal(err)
	}
	if out["X"] != "1" {
		t.Error("vars should pass through unchanged")
	}
}

func TestClear(t *testing.T) {
	m := defaults.New(newMemStore())
	_ = m.Set("dev", map[string]string{"A": "1"})
	_ = m.Clear("dev")
	_, err := m.Get("dev")
	if err == nil {
		t.Error("expected ErrNotFound after clear")
	}
}

func TestList(t *testing.T) {
	m := defaults.New(newMemStore())
	_ = m.Set("dev", map[string]string{"A": "1"})
	_ = m.Set("prod", map[string]string{"B": "2"})
	names, err := m.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(names))
	}
}
