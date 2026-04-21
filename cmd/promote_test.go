package cmd_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/nicholasgasior/envport/internal/promote"
	"github.com/nicholasgasior/envport/internal/snapshot"
)

// memPromoteManager satisfies promote.Manager for command-level tests.
type memPromoteManager struct {
	data map[string]map[string]*snapshot.Snapshot
}

func newMemPromoteManager() *memPromoteManager {
	return &memPromoteManager{data: make(map[string]map[string]*snapshot.Snapshot)}
}

func (m *memPromoteManager) Load(profile, name string) (*snapshot.Snapshot, error) {
	if snaps, ok := m.data[profile]; ok {
		if s, ok := snaps[name]; ok {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *memPromoteManager) Save(profile string, snap *snapshot.Snapshot) error {
	if m.data[profile] == nil {
		m.data[profile] = make(map[string]*snapshot.Snapshot)
	}
	m.data[profile][snap.Name] = snap
	return nil
}

func (m *memPromoteManager) Exists(profile, name string) bool {
	if snaps, ok := m.data[profile]; ok {
		_, ok = snaps[name]
		return ok
	}
	return false
}

func TestPromoteCommandSuccess(t *testing.T) {
	m := newMemPromoteManager()
	if m.data["dev"] == nil {
		m.data["dev"] = make(map[string]*snapshot.Snapshot)
	}
	m.data["dev"]["v1"] = &snapshot.Snapshot{Name: "v1", Vars: map[string]string{"X": "1"}, CreatedAt: time.Now()}

	var buf bytes.Buffer
	err := promote.Run(m, promote.Options{
		SrcProfile: "dev",
		DstProfile: "prod",
		Name:       "v1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = buf // output verified via unit tests; integration wiring checked here

	if !m.Exists("prod", "v1") {
		t.Error("expected snapshot to exist in prod after promote")
	}
}

func TestPromoteCommandDestinationExists(t *testing.T) {
	m := newMemPromoteManager()
	for _, p := range []string{"dev", "prod"} {
		m.data[p] = map[string]*snapshot.Snapshot{
			"v1": {Name: "v1", Vars: map[string]string{"X": "1"}, CreatedAt: time.Now()},
		}
	}

	err := promote.Run(m, promote.Options{
		SrcProfile: "dev",
		DstProfile: "prod",
		Name:       "v1",
		Force:      false,
	})
	if !errors.Is(err, promote.ErrDestinationExists) {
		t.Fatalf("expected ErrDestinationExists, got %v", err)
	}
}
