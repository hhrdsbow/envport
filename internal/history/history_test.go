package history_test

import (
	"errors"
	"sync"
	"testing"

	"envport/internal/history"
)

type memStore struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newMemStore() *memStore { return &memStore{data: map[string][]byte{}} }

func (m *memStore) Get(key string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *memStore) Set(key string, val []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
	return nil
}

func TestRecordAndList(t *testing.T) {
	mgr := history.New(newMemStore())
	if err := mgr.Record("snapshot", "dev", ""); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Record("restore", "dev", "filtered"); err != nil {
		t.Fatal(err)
	}
	entries, err := mgr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Operation != "snapshot" || entries[1].Operation != "restore" {
		t.Errorf("unexpected operations: %+v", entries)
	}
}

func TestListEmpty(t *testing.T) {
	mgr := history.New(newMemStore())
	entries, err := mgr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestClear(t *testing.T) {
	mgr := history.New(newMemStore())
	_ = mgr.Record("snapshot", "prod", "")
	if err := mgr.Clear(); err != nil {
		t.Fatal(err)
	}
	entries, _ := mgr.List()
	if len(entries) != 0 {
		t.Fatalf("expected 0 after clear, got %d", len(entries))
	}
}
