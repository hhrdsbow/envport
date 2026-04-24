package quota_test

import (
	"testing"

	"github.com/envport/envport/internal/quota"
)

type memStore struct {
	data map[string]int
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]int)}
}

func (m *memStore) Get(profile string) (int, error) {
	v, ok := m.data[profile]
	if !ok {
		return 0, nil
	}
	return v, nil
}

func (m *memStore) Set(profile string, limit int) error {
	m.data[profile] = limit
	return nil
}

func (m *memStore) Delete(profile string) error {
	delete(m.data, profile)
	return nil
}

func (m *memStore) List() (map[string]int, error) {
	out := make(map[string]int, len(m.data))
	for k, v := range m.data {
		out[k] = v
	}
	return out, nil
}

func TestSetAndGet(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	if err := mgr.Set("prod", 50); err != nil {
		t.Fatalf("Set: %v", err)
	}

	limit, err := mgr.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if limit != 50 {
		t.Errorf("expected 50, got %d", limit)
	}
}

func TestGetMissing(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	limit, err := mgr.Get("missing")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if limit != 0 {
		t.Errorf("expected 0 for missing profile, got %d", limit)
	}
}

func TestCheckWithinLimit(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	_ = mgr.Set("dev", 10)

	if err := mgr.Check("dev", 5); err != nil {
		t.Errorf("expected no error within limit, got: %v", err)
	}
}

func TestCheckExceedsLimit(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	_ = mgr.Set("dev", 3)

	if err := mgr.Check("dev", 10); err == nil {
		t.Error("expected error when exceeding quota, got nil")
	}
}

func TestCheckNoLimitSet(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	// No limit configured — should always pass
	if err := mgr.Check("any", 9999); err != nil {
		t.Errorf("expected no error with no limit set, got: %v", err)
	}
}

func TestDelete(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	_ = mgr.Set("staging", 20)
	if err := mgr.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	limit, _ := mgr.Get("staging")
	if limit != 0 {
		t.Errorf("expected 0 after delete, got %d", limit)
	}
}

func TestList(t *testing.T) {
	store := newMemStore()
	mgr := quota.New(store)

	_ = mgr.Set("a", 10)
	_ = mgr.Set("b", 20)

	all, err := mgr.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
