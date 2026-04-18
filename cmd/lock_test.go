package cmd

import (
	"bytes"
	"testing"

	"envport/internal/lock"
)

type memLockStore struct {
	entries map[string]lock.Entry
}

func newMemLockStore() *memLockStore {
	return &memLockStore{entries: make(map[string]lock.Entry)}
}
func (m *memLockStore) Set(name string, e lock.Entry) error {
	m.entries[name] = e; return nil
}
func (m *memLockStore) Get(name string) (lock.Entry, bool, error) {
	e, ok := m.entries[name]; return e, ok, nil
}
func (m *memLockStore) Delete(name string) error {
	delete(m.entries, name); return nil
}
func (m *memLockStore) List() ([]lock.Entry, error) {
	out := make([]lock.Entry, 0)
	for _, e := range m.entries { out = append(out, e) }
	return out, nil
}

func TestLockAddAndList(t *testing.T) {
	mgr := lock.New(newMemLockStore())
	if err := mgr.Lock("staging", "freeze"); err != nil {
		t.Fatal(err)
	}
	list, err := mgr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "staging" {
		t.Fatalf("unexpected list: %+v", list)
	}
}

func TestLockRemove(t *testing.T) {
	mgr := lock.New(newMemLockStore())
	_ = mgr.Lock("staging", "")
	if err := mgr.Unlock("staging"); err != nil {
		t.Fatal(err)
	}
	ok, _ := mgr.IsLocked("staging")
	if ok {
		t.Fatal("expected unlocked after remove")
	}
}

func TestLockListOutput(t *testing.T) {
	mgr := lock.New(newMemLockStore())
	_ = mgr.Lock("prod", "release")
	list, _ := mgr.List()
	var buf bytes.Buffer
	for _, e := range list {
		buf.WriteString(e.Name)
	}
	if !bytes.Contains(buf.Bytes(), []byte("prod")) {
		t.Fatal("expected prod in output")
	}
}
