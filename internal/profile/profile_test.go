package profile_test

import (
	"os"
	"testing"

	"github.com/user/envport/internal/profile"
	"github.com/user/envport/internal/store"
)

func newTestManager(t *testing.T) *profile.Manager {
	t.Helper()
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	return profile.New(st)
}

func TestSaveAndLoad(t *testing.T) {
	m := newTestManager(t)
	os.Setenv("ENVPORT_TEST_VAR", "hello")
	t.Cleanup(func() { os.Unsetenv("ENVPORT_TEST_VAR") })

	if err := m.Save("dev"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := m.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Name != "dev" {
		t.Errorf("expected name 'dev', got %q", snap.Name)
	}
	if snap.Vars["ENVPORT_TEST_VAR"] != "hello" {
		t.Errorf("expected ENVPORT_TEST_VAR=hello in snapshot")
	}
}

func TestLoadMissing(t *testing.T) {
	m := newTestManager(t)
	_, err := m.Load("nonexistent")
	if err == nil {
		t.Error("expected error loading missing profile")
	}
}

func TestDelete(t *testing.T) {
	m := newTestManager(t)
	if err := m.Save("tmp"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := m.Delete("tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Load("tmp")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestDeleteMissing(t *testing.T) {
	m := newTestManager(t)
	if err := m.Delete("nonexistent"); err == nil {
		t.Error("expected error deleting nonexistent profile")
	}
}

func TestList(t *testing.T) {
	m := newTestManager(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := m.Save(name); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(names))
	}
}
