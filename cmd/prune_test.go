package cmd

import (
	"bytes"
	"testing"
	"time"

	"envport/internal/prune"
)

// stubPruneManager implements prune.Manager for testing.
type stubPruneManager struct {
	snaps   []prune.Snapshot
	deleted []string
}

type stubSnap struct {
	name string
	at   time.Time
}

func (s stubSnap) Name() string         { return s.name }
func (s stubSnap) CreatedAt() time.Time { return s.at }

func (m *stubPruneManager) List() ([]prune.Snapshot, error) { return m.snaps, nil }
func (m *stubPruneManager) Delete(name string) error {
	m.deleted = append(m.deleted, name)
	return nil
}

func TestPruneOutputDryRun(t *testing.T) {
	now := time.Now()
	m := &stubPruneManager{
		snaps: []prune.Snapshot{
			stubSnap{"recent", now},
			stubSnap{"stale", now.Add(-100 * time.Hour)},
		},
	}

	res, err := prune.Run(m, prune.Options{
		OlderThan: now.Add(-24 * time.Hour),
		DryRun:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Pruned) != 1 || res.Pruned[0] != "stale" {
		t.Fatalf("unexpected pruned list: %v", res.Pruned)
	}
	if len(m.deleted) != 0 {
		t.Fatal("dry run must not delete")
	}
}

func TestPruneOutputNothing(t *testing.T) {
	m := &stubPruneManager{
		snaps: []prune.Snapshot{
			stubSnap{"snap", time.Now()},
		},
	}
	res, err := prune.Run(m, prune.Options{KeepLast: 5})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if len(res.Pruned) != 0 {
		t.Fatal("expected nothing pruned")
	}
	_ = buf
}
