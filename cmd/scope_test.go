package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envport/internal/scope"
)

type memScopeStore struct {
	data map[string][]string
}

func newMemScopeStore() *memScopeStore {
	return &memScopeStore{data: make(map[string][]string)}
}
func (m *memScopeStore) Get(name string) ([]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, scope.ErrNotFound
	}
	return v, nil
}
func (m *memScopeStore) Set(name string, keys []string) error {
	m.data[name] = keys; return nil
}
func (m *memScopeStore) Delete(name string) error {
	delete(m.data, name); return nil
}
func (m *memScopeStore) List() ([]string, error) {
	out := make([]string, 0, len(m.data))
	for k := range m.data {
		out = append(out, k)
	}
	return out, nil
}

func TestScopeAddAndList(t *testing.T) {
	mgr := scope.New(newMemScopeStore())
	if err := mgr.Add("api", []string{"API_KEY", "API_URL"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	names, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 1 || names[0] != "api" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestScopeListOutput(t *testing.T) {
	mgr := scope.New(newMemScopeStore())
	_ = mgr.Add("web", []string{"PORT", "HOST"})
	names, _ := mgr.List()
	var buf bytes.Buffer
	for _, n := range names {
		sc, _ := mgr.Get(n)
		buf.WriteString(n)
		for _, k := range sc.Keys {
			buf.WriteString(" " + k)
		}
	}
	if !strings.Contains(buf.String(), "web") {
		t.Fatal("expected 'web' in output")
	}
	if !strings.Contains(buf.String(), "PORT") {
		t.Fatal("expected 'PORT' in output")
	}
}

func TestScopeRemove(t *testing.T) {
	mgr := scope.New(newMemScopeStore())
	_ = mgr.Add("tmp", []string{"X"})
	_ = mgr.Delete("tmp")
	names, _ := mgr.List()
	if len(names) != 0 {
		t.Fatalf("expected empty list after remove, got %v", names)
	}
}
