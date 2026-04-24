package baseline_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"envport/internal/baseline"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	mu      sync.Mutex
	entries map[string]baseline.Entry
}

func newMemStore() *memStore {
	return &memStore{entries: make(map[string]baseline.Entry)}
}

func (s *memStore) Set(profile, snapshot string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[profile] = baseline.Entry{Profile: profile, Snapshot: snapshot, SetAt: at}
	return nil
}

func (s *memStore) Get(profile string) (baseline.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[profile]
	if !ok {
		return baseline.Entry{}, baseline.ErrNotFound
	}
	return e, nil
}

func (s *memStore) Delete(profile string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, profile)
	return nil
}

func (s *memStore) List() ([]baseline.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]baseline.Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out, nil
}

func TestSetAndGet(t *testing.T) {
	m := baseline.New(newMemStore())
	if err := m.Set("prod", "snap-001"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, err := m.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.Snapshot != "snap-001" {
		t.Errorf("got snapshot %q, want %q", e.Snapshot, "snap-001")
	}
}

func TestGetMissing(t *testing.T) {
	m := baseline.New(newMemStore())
	_, err := m.Get("missing")
	if !errors.Is(err, baseline.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestClear(t *testing.T) {
	m := baseline.New(newMemStore())
	_ = m.Set("dev", "snap-002")
	if err := m.Clear("dev"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	_, err := m.Get("dev")
	if !errors.Is(err, baseline.ErrNotFound) {
		t.Errorf("expected ErrNotFound after clear, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := baseline.New(newMemStore())
	_ = m.Set("prod", "snap-001")
	_ = m.Set("staging", "snap-010")
	entries, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestSetValidation(t *testing.T) {
	m := baseline.New(newMemStore())
	if err := m.Set("", "snap"); err == nil {
		t.Error("expected error for empty profile")
	}
	if err := m.Set("prod", ""); err == nil {
		t.Error("expected error for empty snapshot")
	}
}
