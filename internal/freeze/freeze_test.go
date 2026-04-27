package freeze_test

import (
	"errors"
	"testing"

	"envport/internal/freeze"
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
func (m *memStore) List() (map[string]string, error) {
	copy := make(map[string]string, len(m.data))
	for k, v := range m.data {
		copy[k] = v
	}
	return copy, nil
}
func (m *memStore) Get(k string) (string, error) {
	v := m.data[k]
	return v, nil
}

func TestFreezeAndIsFrozen(t *testing.T) {
	mgr := freeze.New(newMemStore())
	if err := mgr.Freeze("prod"); err != nil {
		t.Fatalf("Freeze: %v", err)
	}
	ok, err := mgr.IsFrozen("prod")
	if err != nil {
		t.Fatalf("IsFrozen: %v", err)
	}
	if !ok {
		t.Error("expected prod to be frozen")
	}
}

func TestUnfreeze(t *testing.T) {
	mgr := freeze.New(newMemStore())
	_ = mgr.Freeze("staging")
	if err := mgr.Unfreeze("staging"); err != nil {
		t.Fatalf("Unfreeze: %v", err)
	}
	ok, _ := mgr.IsFrozen("staging")
	if ok {
		t.Error("expected staging to be unfrozen")
	}
}

func TestUnfreezeMissing(t *testing.T) {
	mgr := freeze.New(newMemStore())
	err := mgr.Unfreeze("ghost")
	if err == nil {
		t.Fatal("expected error unfreezing non-frozen snapshot")
	}
}

func TestList(t *testing.T) {
	mgr := freeze.New(newMemStore())
	_ = mgr.Freeze("a")
	_ = mgr.Freeze("b")
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 frozen snapshots, got %d", len(list))
	}
}

func TestListEmpty(t *testing.T) {
	mgr := freeze.New(newMemStore())
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

var _ = errors.New // suppress unused import
