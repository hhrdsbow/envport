package store

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestOpenEmpty(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if len(s.List()) != 0 {
		t.Error("expected empty store")
	}
}

func TestSetAndGet(t *testing.T) {
	dir := t.TempDir()
	s, _ := Open(dir)

	if err := s.Set("dev", "/tmp/dev.json"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	path, ok := s.Get("dev")
	if !ok || path != "/tmp/dev.json" {
		t.Errorf("Get: got %q, %v", path, ok)
	}
}

func TestPersistence(t *testing.T) {
	dir := t.TempDir()
	s, _ := Open(dir)
	s.Set("prod", "/tmp/prod.json")

	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	path, ok := s2.Get("prod")
	if !ok || path != "/tmp/prod.json" {
		t.Errorf("persistence: got %q, %v", path, ok)
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	s, _ := Open(dir)
	s.Set("staging", "/tmp/staging.json")
	s.Delete("staging")

	if _, ok := s.Get("staging"); ok {
		t.Error("expected entry to be deleted")
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	s, _ := Open(dir)
	s.Set("a", "a.json")
	s.Set("b", "b.json")

	names := s.List()
	slices.Sort(names)
	if !slices.Equal(names, []string{"a", "b"}) {
		t.Errorf("List: got %v", names)
	}
}

func TestStoreFileCreated(t *testing.T) {
	dir := t.TempDir()
	s, _ := Open(dir)
	s.Set("x", "x.json")

	if _, err := os.Stat(filepath.Join(dir, defaultStoreFile)); err != nil {
		t.Errorf("store file not created: %v", err)
	}
}
