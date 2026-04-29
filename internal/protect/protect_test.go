package protect_test

import (
	"errors"
	"testing"

	"envport/internal/protect"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string]struct{}
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]struct{})}
}

func (s *memStore) Set(name string) error {
	s.data[name] = struct{}{}
	return nil
}

func (s *memStore) Delete(name string) error {
	delete(s.data, name)
	return nil
}

func (s *memStore) Exists(name string) (bool, error) {
	_, ok := s.data[name]
	return ok, nil
}

func (s *memStore) List() ([]string, error) {
	out := make([]string, 0, len(s.data))
	for k := range s.data {
		out = append(out, k)
	}
	return out, nil
}

func TestProtectAndIsProtected(t *testing.T) {
	m := protect.New(newMemStore())
	if err := m.Protect("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, err := m.IsProtected("prod")
	if err != nil || !ok {
		t.Fatalf("expected prod to be protected, got ok=%v err=%v", ok, err)
	}
}

func TestProtectAlreadyProtected(t *testing.T) {
	m := protect.New(newMemStore())
	_ = m.Protect("prod")
	err := m.Protect("prod")
	if !errors.Is(err, protect.ErrProtected) {
		t.Fatalf("expected ErrProtected, got %v", err)
	}
}

func TestUnprotect(t *testing.T) {
	m := protect.New(newMemStore())
	_ = m.Protect("staging")
	if err := m.Unprotect("staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, _ := m.IsProtected("staging")
	if ok {
		t.Fatal("expected staging to be unprotected")
	}
}

func TestUnprotectMissing(t *testing.T) {
	m := protect.New(newMemStore())
	if err := m.Unprotect("ghost"); err == nil {
		t.Fatal("expected error for unprotecting non-existent entry")
	}
}

func TestGuardAllowsUnprotected(t *testing.T) {
	m := protect.New(newMemStore())
	if err := m.Guard("dev"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGuardBlocksProtected(t *testing.T) {
	m := protect.New(newMemStore())
	_ = m.Protect("prod")
	if err := m.Guard("prod"); !errors.Is(err, protect.ErrProtected) {
		t.Fatalf("expected ErrProtected, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := protect.New(newMemStore())
	_ = m.Protect("a")
	_ = m.Protect("b")
	list, err := m.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
}
