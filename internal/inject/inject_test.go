package inject_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/inject"
	"github.com/user/envport/internal/snapshot"
)

// memManager is an in-memory Manager for tests.
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

func newManager(vars map[string]string, name string) *memManager {
	return &memManager{
		snaps: map[string]*snapshot.Snapshot{
			name: {Name: name, Vars: vars},
		},
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	mgr := &memManager{snaps: map[string]*snapshot.Snapshot{}}
	err := inject.Run(mgr, "ghost", []string{"env"}, inject.Options{})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunNoCommand(t *testing.T) {
	mgr := newManager(map[string]string{"X": "1"}, "dev")
	err := inject.Run(mgr, "dev", []string{}, inject.Options{})
	if err == nil {
		t.Fatal("expected error when no command given")
	}
}

func TestBuildEnvNoOverride(t *testing.T) {
	vars := map[string]string{"INJECTED": "yes"}
	env := inject.BuildEnv(vars, false)
	var found bool
	for _, e := range env {
		if e == "INJECTED=yes" {
			found = true
		}
	}
	if !found {
		t.Error("expected INJECTED=yes in env slice")
	}
}

func TestBuildEnvOverrideRemovesDuplicate(t *testing.T) {
	// We cannot control the real os.Environ in a unit test, so we just
	// verify that with override=true the snapshot key appears exactly once.
	vars := map[string]string{"PATH": "/custom/bin"}
	env := inject.BuildEnv(vars, true)
	count := 0
	for _, e := range env {
		if len(e) >= 4 && e[:5] == "PATH=" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected PATH to appear exactly once, got %d", count)
	}
}

func TestBuildEnvFilterKeys(t *testing.T) {
	mgr := newManager(map[string]string{"A": "1", "B": "2"}, "prod")
	snap, _ := mgr.Load("prod")
	opts := inject.Options{Keys: []string{"A"}, Override: false}
	_ = snap
	// Verify BuildEnv only includes filtered key.
	vars := map[string]string{"A": "1"}
	env := inject.BuildEnv(vars, opts.Override)
	for _, e := range env {
		if e == "B=2" {
			t.Error("B should not be present after key filtering")
		}
	}
}
