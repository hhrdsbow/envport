package group_test

import (
	"errors"
	"testing"

	"envport/internal/group"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string][]string
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string][]string)}
}

func (s *memStore) Get(name string) ([]string, error) {
	v, ok := s.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (s *memStore) Set(name string, members []string) error {
	s.data[name] = members
	return nil
}

func (s *memStore) Delete(name string) error {
	delete(s.data, name)
	return nil
}

func (s *memStore) List() ([]string, error) {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestCreateAndMembers(t *testing.T) {
	m := group.New(newMemStore())
	if err := m.Create("prod"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	mems, err := m.Members("prod")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if len(mems) != 0 {
		t.Fatalf("expected empty group, got %v", mems)
	}
}

func TestAddAndRemove(t *testing.T) {
	m := group.New(newMemStore())
	_ = m.Create("staging")
	_ = m.Add("staging", "snap-a")
	_ = m.Add("staging", "snap-b")
	_ = m.Add("staging", "snap-a") // idempotent

	mems, _ := m.Members("staging")
	if len(mems) != 2 {
		t.Fatalf("expected 2 members, got %d", len(mems))
	}

	_ = m.Remove("staging", "snap-a")
	mems, _ = m.Members("staging")
	if len(mems) != 1 || mems[0] != "snap-b" {
		t.Fatalf("unexpected members after remove: %v", mems)
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := group.New(newMemStore())
	_ = m.Create("dev")
	err := m.Create("dev")
	if !errors.Is(err, group.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestDeleteMissing(t *testing.T) {
	m := group.New(newMemStore())
	err := m.Delete("ghost")
	if !errors.Is(err, group.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := group.New(newMemStore())
	_ = m.Create("alpha")
	_ = m.Create("beta")
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(names))
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	members := []string{"a", "b", "c"}
	data, err := group.MarshalMembers(members)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	out, err := group.UnmarshalMembers(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(out) != len(members) {
		t.Fatalf("length mismatch: %v", out)
	}
}
