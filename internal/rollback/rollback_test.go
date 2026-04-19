package rollback_test

import (
	"errors"
	"testing"
	"time"

	"envport/internal/history"
	"envport/internal/rollback"
	"envport/internal/snapshot"
)

type memManager struct {
	snaps map[string]*snapshot.Snapshot
}

func (m *memManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.snaps[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *memManager) Save(name string, s *snapshot.Snapshot) error {
	m.snaps[name] = s
	return nil
}

type memHistory struct {
	entries map[string][]history.Entry
}

func (h *memHistory) List(profile string) ([]history.Entry, error) {
	return h.entries[profile], nil
}

func newSetup() (*memManager, *memHistory) {
	mgr := &memManager{snaps: map[string]*snapshot.Snapshot{
		"snap-old": {Vars: map[string]string{"A": "1"}, CreatedAt: time.Now()},
		"snap-new": {Vars: map[string]string{"A": "2"}, CreatedAt: time.Now()},
	}}
	hr := &memHistory{entries: map[string][]history.Entry{
		"dev": {
			{SnapshotRef: "snap-new", Timestamp: time.Now()},
			{SnapshotRef: "snap-old", Timestamp: time.Now().Add(-time.Hour)},
		},
	}}
	return mgr, hr
}

func TestRunRollbackOffset1(t *testing.T) {
	mgr, hr := newSetup()
	res, err := rollback.Run(mgr, hr, "dev", 1)
	if err != nil {
		t.Fatal(err)
	}
	if res.RolledTo != "snap-new" {
		t.Errorf("expected snap-new, got %s", res.RolledTo)
	}
	if mgr.snaps["dev"].Vars["A"] != "2" {
		t.Error("expected A=2 after rollback")
	}
}

func TestRunRollbackOffset2(t *testing.T) {
	mgr, hr := newSetup()
	res, err := rollback.Run(mgr, hr, "dev", 2)
	if err != nil {
		t.Fatal(err)
	}
	if res.RolledTo != "snap-old" {
		t.Errorf("expected snap-old, got %s", res.RolledTo)
	}
}

func TestRunRollbackNoHistory(t *testing.T) {
	mgr, hr := newSetup()
	_, err := rollback.Run(mgr, hr, "prod", 1)
	if err == nil {
		t.Error("expected error for missing history")
	}
}

func TestRunRollbackOffsetTooLarge(t *testing.T) {
	mgr, hr := newSetup()
	_, err := rollback.Run(mgr, hr, "dev", 99)
	if err == nil {
		t.Error("expected error for offset too large")
	}
}
