package prune_test

import (
	"errors"
	"testing"
	"time"

	"envport/internal/prune"
)

type memSnap struct {
	name      string
	createdAt time.Time
}

func (s memSnap) Name() string          { return s.name }
func (s memSnap) CreatedAt() time.Time  { return s.createdAt }

type memManager struct {
	snaps   []prune.Snapshot
	deleted []string
	listErr error
}

func (m *memManager) List() ([]prune.Snapshot, error) { return m.snaps, m.listErr }
func (m *memManager) Delete(name string) error {
	m.deleted = append(m.deleted, name)
	return nil
}

var (
	now   = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	old   = now.Add(-48 * time.Hour)
	older = now.Add(-96 * time.Hour)
)

func newManager() *memManager {
	return &memManager{
		snaps: []prune.Snapshot{
			memSnap{"snap-new", now},
			memSnap{"snap-old", old},
			memSnap{"snap-older", older},
		},
	}
}

func TestPruneOlderThan(t *testing.T) {
	m := newManager()
	res, err := prune.Run(m, prune.Options{OlderThan: now.Add(-24 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(res.Pruned))
	}
	if len(m.deleted) != 2 {
		t.Fatalf("expected 2 deletions, got %d", len(m.deleted))
	}
}

func TestPruneKeepLast(t *testing.T) {
	m := newManager()
	res, err := prune.Run(m, prune.Options{KeepLast: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(res.Pruned))
	}
}

func TestPruneDryRun(t *testing.T) {
	m := newManager()
	res, err := prune.Run(m, prune.Options{KeepLast: 1, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Pruned) != 2 {
		t.Fatalf("expected 2 in result")
	}
	if len(m.deleted) != 0 {
		t.Fatal("dry run should not delete")
	}
}

func TestPruneListError(t *testing.T) {
	m := &memManager{listErr: errors.New("boom")}
	_, err := prune.Run(m, prune.Options{KeepLast: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}
