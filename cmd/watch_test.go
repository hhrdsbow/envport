package cmd

import (
	"bytes"
	"testing"
	"time"

	"envport/internal/snapshot"
	"envport/internal/watch"
)

type stubWatchManager struct {
	snaps map[string]*snapshot.Snapshot
}

func (s *stubWatchManager) Load(name string) (*snapshot.Snapshot, error) {
	snap, ok := s.snaps[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return snap, nil
}

func TestWatchNoDrift(t *testing.T) {
	m := &stubWatchManager{
		snaps: map[string]*snapshot.Snapshot{
			"base": {Name: "base", Vars: map[string]string{}, CreatedAt: time.Now()},
		},
	}
	res, err := watch.Check(m, "base", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Changed {
		t.Errorf("expected no drift, report: %s", res.Report)
	}
}

func TestWatchDriftDetected(t *testing.T) {
	m := &stubWatchManager{
		snaps: map[string]*snapshot.Snapshot{
			"base": {
				Name:      "base",
				Vars:      map[string]string{"ENVPORT_DRIFT_KEY": "expected"},
				CreatedAt: time.Now(),
			},
		},
	}
	res, err := watch.Check(m, "base", []string{"ENVPORT_DRIFT_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Changed {
		t.Error("expected drift to be detected")
	}
	buf := bytes.NewBufferString(res.Report)
	if buf.Len() == 0 {
		t.Error("expected non-empty drift report")
	}
}
