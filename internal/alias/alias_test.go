package alias

import (
	"errors"
	"testing"
)

// memStore is an in-memory Store used in tests.
type memStore struct {
	data map[string]string
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]string)}
}

func (s *memStore) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (s *memStore) Set(key, value string) error {
	s.data[key] = value
	return nil
}

func (s *memStore) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *memStore) List() (map[string]string, error) {
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out, nil
}

func TestAddAndResolve(t *testing.T) {
	m := New(newMemStore())
	if err := m.Add("prod", "snapshot-2024"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	target, err := m.Resolve("prod")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if target != "snapshot-2024" {
		t.Errorf("got %q, want %q", target, "snapshot-2024")
	}
}

func TestResolveMissing(t *testing.T) {
	m := New(newMemStore())
	_, err := m.Resolve("missing")
	if err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestRemove(t *testing.T) {
	m := New(newMemStore())
	_ = m.Add("dev", "snap-dev")
	_ = m.Remove("dev")
	_, err := m.Resolve("dev")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestList(t *testing.T) {
	m := New(newMemStore())
	_ = m.Add("a", "snap-a")
	_ = m.Add("b", "snap-b")
	list, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("got %d entries, want 2", len(list))
	}
}

func TestAddEmptyAlias(t *testing.T) {
	m := New(newMemStore())
	if err := m.Add("", "target"); err == nil {
		t.Fatal("expected error for empty alias")
	}
}
