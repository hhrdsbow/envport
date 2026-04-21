package namespace_test

import (
	"errors"
	"testing"

	"envport/internal/namespace"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string]string
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]string)}
}

func (m *memStore) Get(key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (m *memStore) Set(key, value string) error {
	m.data[key] = value
	return nil
}

func (m *memStore) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func (m *memStore) List() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestAddAndList(t *testing.T) {
	mgr := namespace.New(newMemStore())
	if err := mgr.Add("production"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := mgr.Add("staging"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	names, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 || names[0] != "production" || names[1] != "staging" {
		t.Fatalf("unexpected list: %v", names)
	}
}

func TestAddDuplicate(t *testing.T) {
	mgr := namespace.New(newMemStore())
	_ = mgr.Add("dev")
	if err := mgr.Add("dev"); err == nil {
		t.Fatal("expected error for duplicate namespace")
	}
}

func TestAddEmpty(t *testing.T) {
	mgr := namespace.New(newMemStore())
	if err := mgr.Add(""); err == nil {
		t.Fatal("expected error for empty namespace name")
	}
}

func TestRemove(t *testing.T) {
	mgr := namespace.New(newMemStore())
	_ = mgr.Add("temp")
	if err := mgr.Remove("temp"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	names, _ := mgr.List()
	if len(names) != 0 {
		t.Fatalf("expected empty list, got %v", names)
	}
}

func TestRemoveMissing(t *testing.T) {
	mgr := namespace.New(newMemStore())
	err := mgr.Remove("ghost")
	if !errors.Is(err, namespace.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestExists(t *testing.T) {
	mgr := namespace.New(newMemStore())
	_ = mgr.Add("alpha")
	if !mgr.Exists("alpha") {
		t.Fatal("expected Exists to return true")
	}
	if mgr.Exists("beta") {
		t.Fatal("expected Exists to return false")
	}
}
