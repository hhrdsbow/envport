package cmd

import (
	"bytes"
	"errors"
	"testing"

	"envport/internal/history"
	"envport/internal/rollback"
	"envport/internal/snapshot"
	"time"
)

type memRollbackManager struct {
	snaps map[string]*snapshot.Snapshot
}

func (m *memRollbackManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.snaps[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}
func (m *memRollbackManager) Save(name string, s *snapshot.Snapshot) error {
	m.snaps[name] = s
	return nil
}

type memRollbackHistory struct {
	entries []history.Entry
}

func (h *memRollbackHistory) List(_ string) ([]history.Entry, error) {
	return h.entries, nil
}

func TestRollbackCommand(t *testing.T) {
	mgr := &memRollbackManager{snaps: map[string]*snapshot.Snapshot{
		"snap-a": {Vars: map[string]string{"X": "old"}, CreatedAt: time.Now()},
	}}
	hr := &memRollbackHistory{entries: []history.Entry{
		{SnapshotRef: "snap-a", Timestamp: time.Now()},
	}}

	res, err := rollback.Run(mgr, hr, "myprofile", 1)
	if err != nil {
		t.Fatal(err)
	}
	if res.RolledTo != "snap-a" {
		t.Errorf("unexpected ref: %s", res.RolledTo)
	}
}

func TestRollbackCommandMissingProfile(t *testing.T) {
	mgr := &memRollbackManager{snaps: map[string]*snapshot.Snapshot{}}
	hr := &memRollbackHistory{entries: nil}

	_, err := rollback.Run(mgr, hr, "ghost", 1)
	if err == nil {
		t.Error("expected error")
	}

	var buf bytes.Buffer
	buf.WriteString(err.Error())
	if buf.Len() == 0 {
		t.Error("expected non-empty error message")
	}
}
