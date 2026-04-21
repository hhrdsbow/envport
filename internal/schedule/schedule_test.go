package schedule_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"envport/internal/schedule"
)

// memStore is an in-memory Store implementation for tests.
type memStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemStore() *memStore { return &memStore{data: make(map[string][]byte)} }

func (s *memStore) Set(name string, data []byte) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.data[name] = data
	return nil
}
func (s *memStore) Get(name string) ([]byte, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	d, ok := s.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return d, nil
}
func (s *memStore) Delete(name string) error {
	s.mu.Lock(); defer s.mu.Unlock()
	delete(s.data, name)
	return nil
}
func (s *memStore) List() ([]string, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	var keys []string
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func newManager() *schedule.Manager { return schedule.New(newMemStore()) }

func TestAddAndGet(t *testing.T) {
	m := newManager()
	e := schedule.Entry{Name: "nightly", Profile: "prod", Interval: 24 * time.Hour}
	if err := m.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}
	got, err := m.Get("nightly")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Profile != "prod" {
		t.Errorf("profile = %q, want %q", got.Profile, "prod")
	}
	if !got.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestAddValidation(t *testing.T) {
	m := newManager()
	if err := m.Add(schedule.Entry{Profile: "p", Interval: time.Hour}); err == nil {
		t.Error("expected error for empty name")
	}
	if err := m.Add(schedule.Entry{Name: "x", Interval: time.Hour}); err == nil {
		t.Error("expected error for empty profile")
	}
	if err := m.Add(schedule.Entry{Name: "x", Profile: "p", Interval: -1}); err == nil {
		t.Error("expected error for non-positive interval")
	}
}

func TestRemove(t *testing.T) {
	m := newManager()
	_ = m.Add(schedule.Entry{Name: "s", Profile: "dev", Interval: time.Hour})
	_ = m.Remove("s")
	if _, err := m.Get("s"); err == nil {
		t.Error("expected error after remove")
	}
}

func TestList(t *testing.T) {
	m := newManager()
	_ = m.Add(schedule.Entry{Name: "a", Profile: "dev", Interval: time.Hour})
	_ = m.Add(schedule.Entry{Name: "b", Profile: "prod", Interval: 12 * time.Hour})
	entries, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("len = %d, want 2", len(entries))
	}
}

func TestUpdateLastRun(t *testing.T) {
	m := newManager()
	_ = m.Add(schedule.Entry{Name: "s", Profile: "dev", Interval: time.Hour})
	now := time.Now().Truncate(time.Second)
	if err := m.UpdateLastRun("s", now); err != nil {
		t.Fatalf("UpdateLastRun: %v", err)
	}
	got, _ := m.Get("s")
	if !got.LastRun.Equal(now) {
		t.Errorf("LastRun = %v, want %v", got.LastRun, now)
	}
}
