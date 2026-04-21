package expire_test

import (
	"errors"
	"testing"
	"time"

	"envport/internal/expire"
)

// memStore is an in-memory Store implementation for tests.
type memStore struct {
	data map[string]time.Time
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]time.Time)}
}

func (s *memStore) Set(name string, t time.Time) error {
	s.data[name] = t
	return nil
}

func (s *memStore) Get(name string) (time.Time, error) {
	t, ok := s.data[name]
	if !ok {
		return time.Time{}, expire.ErrNotFound
	}
	return t, nil
}

func (s *memStore) Delete(name string) error {
	delete(s.data, name)
	return nil
}

func (s *memStore) List() (map[string]time.Time, error) {
	copy := make(map[string]time.Time, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy, nil
}

func TestSetAndGet(t *testing.T) {
	m := expire.New(newMemStore())
	if err := m.Set("snap1", 10*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("snap1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if time.Until(got) <= 0 {
		t.Errorf("expected future expiry, got %v", got)
	}
}

func TestSetInvalidTTL(t *testing.T) {
	m := expire.New(newMemStore())
	if err := m.Set("snap1", -1*time.Second); err == nil {
		t.Error("expected error for negative ttl")
	}
}

func TestExpired(t *testing.T) {
	s := newMemStore()
	s.data["old"] = time.Now().Add(-1 * time.Hour)
	s.data["fresh"] = time.Now().Add(1 * time.Hour)
	m := expire.New(s)
	expired, err := m.Expired()
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	if len(expired) != 1 || expired[0] != "old" {
		t.Errorf("expected [old], got %v", expired)
	}
}

func TestRemove(t *testing.T) {
	m := expire.New(newMemStore())
	_ = m.Set("snap1", time.Minute)
	if err := m.Remove("snap1"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := m.Get("snap1")
	if !errors.Is(err, expire.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestFormat(t *testing.T) {
	out := expire.Format(map[string]time.Time{})
	if out == "" {
		t.Error("expected non-empty output for empty map")
	}
	records := map[string]time.Time{
		"snap1": time.Now().Add(time.Hour),
		"snap2": time.Now().Add(-time.Hour),
	}
	out = expire.Format(records)
	if len(out) == 0 {
		t.Error("expected formatted output")
	}
}
