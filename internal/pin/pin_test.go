package pin_test

import (
	"errors"
	"testing"

	"envport/internal/pin"
)

// newMemStore returns a simple in-memory Store implementation for tests.
func newMemStore() pin.Store {
	return &memStore{data: map[string]string{}}
}

type memStore struct{ data map[string]string }

func (m *memStore) Set(k, v string) error          { m.data[k] = v; return nil }
func (m *memStore) Get(k string) (string, error) {
	v, ok := m.data[k]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (m *memStore) Delete(k string) error {
	delete(m.data, k)
	return nil
}
func (m *memStore) List() (map[string]string, error) {
	copy := make(map[string]string, len(m.data))
	for k, v := range m.data {
		copy[k] = v
	}
	return copy, nil
}

func TestAddAndIsPinned(t *testing.T) {
	m := pin.New(newMemStore())
	if m.IsPinned("prod") {
		t.Fatal("expected prod to not be pinned initially")
	}
	if err := m.Add("prod"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !m.IsPinned("prod") {
		t.Fatal("expected prod to be pinned after Add")
	}
}

func TestRemovePinned(t *testing.T) {
	m := pin.New(newMemStore())
	_ = m.Add("staging")
	if err := m.Remove("staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if m.IsPinned("staging") {
		t.Fatal("expected staging to be unpinned after Remove")
	}
}

func TestRemoveNotPinned(t *testing.T) {
	m := pin.New(newMemStore())
	err := m.Remove("ghost")
	if !errors.Is(err, pin.ErrNotPinned) {
		t.Fatalf("expected ErrNotPinned, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := pin.New(newMemStore())
	_ = m.Add("b")
	_ = m.Add("a")
	_ = m.Add("c")
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Fatalf("expected 3 pinned names, got %d", len(names))
	}
	for i, want := range []string{"a", "b", "c"} {
		if names[i] != want {
			t.Errorf("names[%d]: want %q, got %q", i, want, names[i])
		}
	}
}

func TestListEmpty(t *testing.T) {
	m := pin.New(newMemStore())
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected empty list, got %v", names)
	}
}
