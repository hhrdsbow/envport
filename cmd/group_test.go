package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"envport/internal/group"
)

// memGroupStore is an in-memory group.Store for cmd-level tests.
type memGroupStore struct {
	data map[string][]string
}

func newMemGroupStore() *memGroupStore {
	return &memGroupStore{data: make(map[string][]string)}
}

func (s *memGroupStore) Get(name string) ([]string, error) {
	v, ok := s.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (s *memGroupStore) Set(name string, members []string) error {
	s.data[name] = members
	return nil
}

func (s *memGroupStore) Delete(name string) error {
	delete(s.data, name)
	return nil
}

func (s *memGroupStore) List() ([]string, error) {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestGroupAddAndList(t *testing.T) {
	st := newMemGroupStore()
	m := group.New(st)

	_ = m.Add("prod", "snap-1")
	_ = m.Add("prod", "snap-2")

	mems, err := m.Members("prod")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if len(mems) != 2 {
		t.Fatalf("expected 2, got %d", len(mems))
	}
}

func TestGroupRemove(t *testing.T) {
	st := newMemGroupStore()
	m := group.New(st)

	_ = m.Add("dev", "snap-a")
	_ = m.Add("dev", "snap-b")
	_ = m.Remove("dev", "snap-a")

	mems, _ := m.Members("dev")
	if len(mems) != 1 || mems[0] != "snap-b" {
		t.Fatalf("unexpected: %v", mems)
	}
}

func TestGroupListOutput(t *testing.T) {
	st := newMemGroupStore()
	m := group.New(st)
	_ = m.Add("staging", "alpha")
	_ = m.Add("staging", "beta")

	mems, _ := m.Members("staging")
	var buf bytes.Buffer
	for _, mem := range mems {
		buf.WriteString(mem + "\n")
	}
	out := buf.String()
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Fatalf("output missing members: %q", out)
	}
}

func TestGroupMissingReturnsError(t *testing.T) {
	st := newMemGroupStore()
	m := group.New(st)
	_, err := m.Members("ghost")
	if !errors.Is(err, group.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
