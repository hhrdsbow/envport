package watch_test

import (
	"errors"
	"testing"
	"time"

	"envport/internal/snapshot"
	"envport/internal/watch"
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

func newManager(vars map[string]string) *memManager {
	return &memManager{
		snaps: map[string]*snapshot.Snapshot{
			"test": {Name: "test", Vars: vars, CreatedAt: time.Now()},
		},
	}
}

func TestCheckNoChanges(t *testing.T) {
	// snapshot matches current env key subset
	m := newManager(map[string]string{})
	res, err := watch.Check(m, "test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Changed {
		t.Errorf("expected no changes, got report: %s", res.Report)
	}
}

func TestCheckMissingSnapshot(t *testing.T) {
	m := &memManager{snaps: map[string]*snapshot.Snapshot{}}
	_, err := watch.Check(m, "missing", nil)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestCheckDetectsChanges(t *testing.T) {
	// snapshot has a key that current env almost certainly won't have with this value
	m := newManager(map[string]string{"ENVPORT_WATCH_SENTINEL_XYZ": "old_value"})
	res, err := watch.Check(m, "test", []string{"ENVPORT_WATCH_SENTINEL_XYZ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// sentinel key is in snapshot but not in current env, so it's a removal → changed
	if !res.Changed {
		t.Error("expected changes to be detected")
	}
	if res.Report == "" {
		t.Error("expected non-empty report")
	}
}
