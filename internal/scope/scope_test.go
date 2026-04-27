package scope_test

import (
	"testing"

	"envport/internal/scope"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string][]string
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string][]string)}
}

func (m *memStore) Get(name string) ([]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, scope.ErrNotFound
	}
	return v, nil
}

func (m *memStore) Set(name string, keys []string) error {
	m.data[name] = keys
	return nil
}

func (m *memStore) Delete(name string) error {
	delete(m.data, name)
	return nil
}

func (m *memStore) List() ([]string, error) {
	out := make([]string, 0, len(m.data))
	for k := range m.data {
		out = append(out, k)
	}
	return out, nil
}

func TestAddAndGet(t *testing.T) {
	mgr := scope.New(newMemStore())
	if err := mgr.Add("web", []string{"PORT", "HOST"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	sc, err := mgr.Get("web")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if sc.Name != "web" || len(sc.Keys) != 2 {
		t.Fatalf("unexpected scope: %+v", sc)
	}
}

func TestAddEmptyName(t *testing.T) {
	mgr := scope.New(newMemStore())
	if err := mgr.Add("", []string{"A"}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestDelete(t *testing.T) {
	mgr := scope.New(newMemStore())
	_ = mgr.Add("db", []string{"DB_HOST"})
	_ = mgr.Delete("db")
	if _, err := mgr.Get("db"); err == nil {
		t.Fatal("expected ErrNotFound after delete")
	}
}

func TestList(t *testing.T) {
	mgr := scope.New(newMemStore())
	_ = mgr.Add("a", []string{"X"})
	_ = mgr.Add("b", []string{"Y"})
	names, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(names))
	}
}

func TestApply(t *testing.T) {
	mgr := scope.New(newMemStore())
	_ = mgr.Add("web", []string{"PORT", "HOST"})
	vars := map[string]string{"PORT": "8080", "HOST": "localhost", "SECRET": "s3cr3t"}
	out, err := mgr.Apply("web", vars)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["SECRET"]; ok {
		t.Fatal("SECRET should not appear in scoped output")
	}
}
