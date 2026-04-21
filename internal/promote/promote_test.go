package promote_test

import (
	"errors"
	"testing"
	"time"

	"github.com/nicholasgasior/envport/internal/promote"
	"github.com/nicholasgasior/envport/internal/snapshot"
)

// memManager is an in-memory implementation of promote.Manager.
type memManager struct {
	data map[string]map[string]*snapshot.Snapshot // profile -> name -> snap
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]*snapshot.Snapshot)}
}

func (m *memManager) Load(profile, name string) (*snapshot.Snapshot, error) {
	if snaps, ok := m.data[profile]; ok {
		if s, ok := snaps[name]; ok {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *memManager) Save(profile string, snap *snapshot.Snapshot) error {
	if m.data[profile] == nil {
		m.data[profile] = make(map[string]*snapshot.Snapshot)
	}
	m.data[profile][snap.Name] = snap
	return nil
}

func (m *memManager) Exists(profile, name string) bool {
	if snaps, ok := m.data[profile]; ok {
		_, ok = snaps[name]
		return ok
	}
	return false
}

func seed(m *memManager, profile, name string, vars map[string]string) {
	if m.data[profile] == nil {
		m.data[profile] = make(map[string]*snapshot.Snapshot)
	}
	m.data[profile][name] = &snapshot.Snapshot{Name: name, Vars: vars, CreatedAt: time.Now()}
}

func TestRunPromotesSnapshot(t *testing.T) {
	m := newMemManager()
	seed(m, "staging", "v1", map[string]string{"FOO": "bar"})

	err := promote.Run(m, promote.Options{SrcProfile: "staging", DstProfile: "production", Name: "v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := m.Load("production", "v1")
	if err != nil {
		t.Fatalf("snapshot not found in destination: %v", err)
	}
	if snap.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", snap.Vars["FOO"])
	}
}

func TestRunDestinationExistsNoForce(t *testing.T) {
	m := newMemManager()
	seed(m, "staging", "v1", map[string]string{"A": "1"})
	seed(m, "production", "v1", map[string]string{"A": "old"})

	err := promote.Run(m, promote.Options{SrcProfile: "staging", DstProfile: "production", Name: "v1"})
	if !errors.Is(err, promote.ErrDestinationExists) {
		t.Fatalf("expected ErrDestinationExists, got %v", err)
	}
}

func TestRunDestinationExistsWithForce(t *testing.T) {
	m := newMemManager()
	seed(m, "staging", "v1", map[string]string{"A": "new"})
	seed(m, "production", "v1", map[string]string{"A": "old"})

	err := promote.Run(m, promote.Options{SrcProfile: "staging", DstProfile: "production", Name: "v1", Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap, _ := m.Load("production", "v1")
	if snap.Vars["A"] != "new" {
		t.Errorf("expected A=new after force promote, got %s", snap.Vars["A"])
	}
}

func TestRunSameProfileError(t *testing.T) {
	m := newMemManager()
	seed(m, "staging", "v1", map[string]string{})

	err := promote.Run(m, promote.Options{SrcProfile: "staging", DstProfile: "staging", Name: "v1"})
	if err == nil {
		t.Fatal("expected error for same src/dst profile")
	}
}

func TestRunMissingSourceSnapshot(t *testing.T) {
	m := newMemManager()

	err := promote.Run(m, promote.Options{SrcProfile: "staging", DstProfile: "production", Name: "missing"})
	if err == nil {
		t.Fatal("expected error for missing source snapshot")
	}
}
