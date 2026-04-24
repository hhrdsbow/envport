package archive_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/envport/internal/archive"
	"github.com/user/envport/internal/snapshot"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	data map[string]*snapshot.Snapshot
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]*snapshot.Snapshot)}
}

func (m *memManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *memManager) Save(name string, s *snapshot.Snapshot) error {
	m.data[name] = s
	return nil
}

func (m *memManager) List() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = snapshot.New(vars)
}

func TestPackAndUnpack(t *testing.T) {
	src := newMemManager()
	seed(src, "dev", map[string]string{"FOO": "bar", "PORT": "8080"})
	seed(src, "prod", map[string]string{"FOO": "baz", "PORT": "443"})

	ar, err := archive.Pack(src, []string{"dev", "prod"})
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	if len(ar.Snapshots) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(ar.Snapshots))
	}

	dst := newMemManager()
	restored, err := archive.Unpack(dst, ar)
	if err != nil {
		t.Fatalf("Unpack: %v", err)
	}
	if len(restored) != 2 {
		t.Fatalf("expected 2 restored, got %d", len(restored))
	}
	snap, _ := dst.Load("dev")
	if snap.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", snap.Vars["FOO"])
	}
}

func TestPackMissing(t *testing.T) {
	m := newMemManager()
	_, err := archive.Pack(m, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	ar := &archive.Archive{
		CreatedAt: time.Now().UTC(),
		Snapshots: map[string]*snapshot.Snapshot{
			"staging": snapshot.New(map[string]string{"ENV": "staging"}),
		},
	}

	data, err := archive.Marshal(ar)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got, err := archive.Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Snapshots["staging"].Vars["ENV"] != "staging" {
		t.Errorf("unexpected value after round-trip")
	}
}

func TestUnmarshalInvalid(t *testing.T) {
	_, err := archive.Unmarshal([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
