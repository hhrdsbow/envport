package cmd_test

import (
	"bytes"
	"testing"
	"time"

	"envport/internal/expire"
)

// inMemExpireStore satisfies expire.Store for tests.
type inMemExpireStore struct {
	data map[string]time.Time
}

func newInMemExpireStore() *inMemExpireStore {
	return &inMemExpireStore{data: make(map[string]time.Time)}
}

func (s *inMemExpireStore) Set(name string, t time.Time) error {
	s.data[name] = t
	return nil
}
func (s *inMemExpireStore) Get(name string) (time.Time, error) {
	v, ok := s.data[name]
	if !ok {
		return time.Time{}, expire.ErrNotFound
	}
	return v, nil
}
func (s *inMemExpireStore) Delete(name string) error {
	delete(s.data, name)
	return nil
}
func (s *inMemExpireStore) List() (map[string]time.Time, error) {
	copy := make(map[string]time.Time, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy, nil
}

func TestExpireSetAndList(t *testing.T) {
	s := newInMemExpireStore()
	m := expire.New(s)

	if err := m.Set("mysnap", 30*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}

	all, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if _, ok := all["mysnap"]; !ok {
		t.Error("expected mysnap in expiry list")
	}
}

func TestExpireFormatOutput(t *testing.T) {
	records := map[string]time.Time{
		"snap-a": time.Now().Add(time.Hour),
	}
	out := expire.Format(records)
	if !bytes.Contains([]byte(out), []byte("snap-a")) {
		t.Errorf("expected snap-a in output, got: %s", out)
	}
	if !bytes.Contains([]byte(out), []byte("active")) {
		t.Errorf("expected 'active' status in output, got: %s", out)
	}
}

func TestExpireExpiredDetection(t *testing.T) {
	s := newInMemExpireStore()
	s.data["stale"] = time.Now().Add(-2 * time.Hour)
	s.data["live"] = time.Now().Add(2 * time.Hour)
	m := expire.New(s)

	expired, err := m.Expired()
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	if len(expired) != 1 || expired[0] != "stale" {
		t.Errorf("expected [stale], got %v", expired)
	}
}
