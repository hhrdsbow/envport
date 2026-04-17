package cmd_test

import (
	"bytes"
	"testing"

	"github.com/user/envport/internal/tag"
)

type memTagStore struct {
	data map[string][]byte
}

func newMemTagStore() *memTagStore { return &memTagStore{data: map[string][]byte{}} }
func (m *memTagStore) Get(k string) ([]byte, error) {
	v, ok := m.data[k]
	if !ok {
		return nil, bytes.ErrTooLarge
	}
	return v, nil
}
func (m *memTagStore) Set(k string, v []byte) error { m.data[k] = v; return nil }
func (m *memTagStore) Delete(k string) error        { delete(m.data, k); return nil }
func (m *memTagStore) List() ([]string, error) {
	keys := make([]string, 0)
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func TestTagAddAndGet(t *testing.T) {
	mgr := tag.New(newMemTagStore())
	if err := mgr.Add("release", "snap-v1"); err != nil {
		t.Fatal(err)
	}
	snaps, err := mgr.Get("release")
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 1 || snaps[0] != "snap-v1" {
		t.Fatalf("unexpected snaps: %v", snaps)
	}
}

func TestTagRemoveCleansUp(t *testing.T) {
	mgr := tag.New(newMemTagStore())
	mgr.Add("release", "snap-v1")
	if err := mgr.Remove("release", "snap-v1"); err != nil {
		t.Fatal(err)
	}
	tags, _ := mgr.List()
	if len(tags) != 0 {
		t.Fatalf("expected no tags, got %v", tags)
	}
}

func TestTagListMultiple(t *testing.T) {
	mgr := tag.New(newMemTagStore())
	mgr.Add("alpha", "snap1")
	mgr.Add("beta", "snap2")
	mgr.Add("alpha", "snap3")
	tags, err := mgr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
}
