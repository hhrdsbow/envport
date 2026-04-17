package restore_test

import (
	"os"
	"testing"

	"envport/internal/profile"
	"envport/internal/restore"
	"envport/internal/snapshot"
)

func newTestManager(t *testing.T) *profile.Manager {
	t.Helper()
	dir := t.TempDir()
	mgr, err := profile.NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return mgr
}

func TestRunDryRun(t *testing.T) {
	mgr := newTestManager(t)
	snap := snapshot.New(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err := mgr.Save("proj", "base", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	res, err := restore.Run(mgr, restore.Options{
		ProfileName:  "proj",
		SnapshotName: "base",
		DryRun:       true,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true")
	}
	if res.Applied["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Applied["FOO"])
	}
	if _, set := os.LookupEnv("FOO"); set {
		t.Error("FOO should not be set in dry-run mode")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	mgr := newTestManager(t)
	_, err := restore.Run(mgr, restore.Options{
		ProfileName:  "proj",
		SnapshotName: "missing",
	})
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRunFilterKeys(t *testing.T) {
	mgr := newTestManager(t)
	snap := snapshot.New(map[string]string{"A": "1", "B": "2", "C": "3"})
	if err := mgr.Save("proj", "base", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	res, err := restore.Run(mgr, restore.Options{
		ProfileName:  "proj",
		SnapshotName: "base",
		DryRun:       true,
		FilterKeys:   []string{"A", "C"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied keys, got %d", len(res.Applied))
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Errorf("expected B skipped, got %v", res.Skipped)
	}
}
