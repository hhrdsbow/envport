package lock_test

import (
	"testing"

	"envport/internal/lock"
)

type memStore struct {
	entries map[string]lock.Entry
}

func newMemStore() *memStore {
	return &memStore{entries: make(map[string]lock.Entry)}
}

func (m *memStore) Set(name string, e lock.Entry) error {
	m.entries[name] = e
	return nil
}

func (m *memStore) Get(name string) (lock.Entry, bool, error) {
	e, ok := m.entries[name]
	return e, ok, nil
}

func (m *memStore) Delete(name string) error {
	delete(m.entries, name)
	return nil
}

func (m *memStore) List() ([]lock.Entry, error) {
	out := make([]lock.Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	return out, nil
}

func TestLockAndUnlock(t *testing.T) {
	mgr := lock.New(newMemStore())
	if err := mgr.Lock("prod", "deploying"); err != nil {
		t.Fatal(err)
	}
	ok, err := mgr.IsLocked("prod")
	if err != nil || !ok {
		t.Fatal("expected locked")
	}
	if err := mgr.Unlock("prod"); err != nil {
		t.Fatal(err)
	}
	ok, _ = mgr.IsLocked("prod")
	if ok {
		t.Fatal("expected unlocked")
	}
}

func TestDoubleLock(t *testing.T) {
	mgr := lock.New(newMemStore())
	_ = mgr.Lock("prod", "")
	if err := mgr.Lock("prod", ""); err != lock.ErrAlreadyLocked {
		t.Fatalf("expected ErrAlreadyLocked, got %v", err)
	}
}

func TestUnlockMissing(t *testing.T) {
	mgr := lock.New(newMemStore())
	if err := mgr.Unlock("missing"); err != lock.ErrNotLocked {
		t.Fatalf("expected ErrNotLocked, got %v", err)
	}
}

func TestList(t *testing.T) {
	mgr := lock.New(newMemStore())
	_ = mgr.Lock("a", "reason a")
	_ = mgr.Lock("b", "reason b")
	list, err := mgr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 locks, got %d", len(list))
	}
}
