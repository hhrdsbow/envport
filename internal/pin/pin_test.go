package pin_test

import (
	"errors"
	"sync"
	"testing"

	"envport/internal/pin"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	mu   sync.Mutex
	data map[string]string
}

func newMemStore() *memStore { return &memStore{data: map[string]string{}} }

func (s *memStore) Set(k, v string) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.data[k] = v; return nil
}
func (s *memStore) Get(k string) (string, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	v, ok := s.data[k]
	if !ok { return "", errors.New("not found") }
	return v, nil
}
func (s *memStore) Delete(k string) error {
	s.mu.Lock(); defer s.mu.Unlock()
	delete(s.data, k); return nil
}
func (s *memStore) List() (map[string]string, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	copy := map[string]string{}
	for k, v := range s.data { copy[k] = v }
	return copy, nil
}

func TestPinAndResolve(t *testing.T) {
	m := pin.New(newMemStore())
	if err := m.Pin("prod", "snap-001"); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	got, err := m.Resolve("prod")
	if err != nil || got != "snap-001" {
		t.Fatalf("Resolve: got %q, err %v", got, err)
	}
}

func TestResolveNotFound(t *testing.T) {
	m := pin.New(newMemStore())
	_, err := m.Resolve("missing")
	if !errors.Is(err, pin.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUnpin(t *testing.T) {
	m := pin.New(newMemStore())
	_ = m.Pin("dev", "snap-002")
	_ = m.Unpin("dev")
	_, err := m.Resolve("dev")
	if !errors.Is(err, pin.ErrNotFound) {
		t.Fatal("expected ErrNotFound after unpin")
	}
}

func TestList(t *testing.T) {
	m := pin.New(newMemStore())
	_ = m.Pin("b", "snap-b")
	_ = m.Pin("a", "snap-a")
	list, err := m.List()
	if err != nil || len(list) != 2 {
		t.Fatalf("List: %v %v", list, err)
	}
	if list[0][0] != "a" || list[1][0] != "b" {
		t.Fatalf("List not sorted: %v", list)
	}
}

func TestPinValidation(t *testing.T) {
	m := pin.New(newMemStore())
	if err := m.Pin("", "snap"); err == nil {
		t.Fatal("expected error for empty alias")
	}
	if err := m.Pin("alias", ""); err == nil {
		t.Fatal("expected error for empty snapshot")
	}
}
