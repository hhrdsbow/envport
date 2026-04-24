package annotate_test

import (
	"fmt"
	"testing"

	"github.com/user/envport/internal/annotate"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string]string
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]string)}
}

func (m *memStore) Set(k, v string) error        { m.data[k] = v; return nil }
func (m *memStore) Delete(k string) error        { delete(m.data, k); return nil }
func (m *memStore) Get(k string) (string, error) {
	v, ok := m.data[k]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return v, nil
}
func (m *memStore) List() (map[string]string, error) {
	out := make(map[string]string, len(m.data))
	for k, v := range m.data {
		out[k] = v
	}
	return out, nil
}

func TestSetAndGet(t *testing.T) {
	mgr := annotate.New(newMemStore())
	if err := mgr.Set("prod", "production snapshot"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	note, err := mgr.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if note != "production snapshot" {
		t.Errorf("expected %q, got %q", "production snapshot", note)
	}
}

func TestGetMissing(t *testing.T) {
	mgr := annotate.New(newMemStore())
	note, err := mgr.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note != "" {
		t.Errorf("expected empty string, got %q", note)
	}
}

func TestRemove(t *testing.T) {
	mgr := annotate.New(newMemStore())
	_ = mgr.Set("dev", "dev note")
	if err := mgr.Remove("dev"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	note, _ := mgr.Get("dev")
	if note != "" {
		t.Errorf("expected empty after remove, got %q", note)
	}
}

func TestList(t *testing.T) {
	mgr := annotate.New(newMemStore())
	_ = mgr.Set("alpha", "first")
	_ = mgr.Set("beta", "second")
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 entries, got %d", len(list))
	}
	if list["alpha"] != "first" || list["beta"] != "second" {
		t.Errorf("unexpected list contents: %v", list)
	}
}

func TestSetEmptyName(t *testing.T) {
	mgr := annotate.New(newMemStore())
	if err := mgr.Set("", "note"); err == nil {
		t.Error("expected error for empty name")
	}
}
