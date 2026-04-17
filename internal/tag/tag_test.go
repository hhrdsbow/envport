package tag_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/user/envport/internal/tag"
)

// memStore is a simple in-memory store for testing.
type memStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemStore() *memStore { return &memStore{data: map[string][]byte{}} }

func (m *memStore) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
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
func (m *memStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}
func (m *memStore) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestAddAndGet(t *testing.T) {
	mgr := tag.New(newMemStore())
	if err := mgr.Add("prod", "snap1"); err != nil {
		t.Fatal(err)
	}
	snaps, err := mgr.Get("prod")
	if err != nil || len(snaps) != 1 || snaps[0] != "snap1" {
		t.Fatalf("expected [snap1], got %v %v", snaps, err)
	}
}

func TestAddDuplicate(t *testing.T) {
	mgr := tag.New(newMemStore())
	mgr.Add("prod", "snap1")
	mgr.Add("prod", "snap1")
	snaps, _ := mgr.Get("prod")
	if len(snaps) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snaps))
	}
}

func TestRemove(t *testing.T) {
	mgr := tag.New(newMemStore())
	mgr.Add("prod", "snap1")
	mgr.Add("prod", "snap2")
	if err := mgr.Remove("prod", "snap1"); err != nil {
		t.Fatal(err)
	}
	snaps, _ := mgr.Get("prod")
	if len(snaps) != 1 || snaps[0] != "snap2" {
		t.Fatalf("expected [snap2], got %v", snaps)
	}
}

func TestList(t *testing.T) {
	mgr := tag.New(newMemStore())
	mgr.Add("prod", "snap1")
	mgr.Add("staging", "snap2")
	tags, err := mgr.List()
	if err != nil || len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v %v", tags, err)
	}
}
