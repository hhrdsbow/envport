package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	environ := []string{"HOME=/root", "PATH=/usr/bin:/bin", "EDITOR=vim"}
	s := New("test", environ)

	if s.Name != "test" {
		t.Errorf("expected name 'test', got %q", s.Name)
	}
	if len(s.Env) != 3 {
		t.Errorf("expected 3 env vars, got %d", len(s.Env))
	}
	if s.Env["HOME"] != "/root" {
		t.Errorf("expected HOME=/root, got %q", s.Env["HOME"])
	}
	if s.Env["PATH"] != "/usr/bin:/bin" {
		t.Errorf("expected PATH=/usr/bin:/bin, got %q", s.Env["PATH"])
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := New("mysnap", []string{"FOO=bar", "BAZ=qux"})
	if err := original.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Name != original.Name {
		t.Errorf("name mismatch: got %q, want %q", loaded.Name, original.Name)
	}
	if loaded.Env["FOO"] != "bar" {
		t.Errorf("FOO mismatch: got %q", loaded.Env["FOO"])
	}
}

func TestLoadMissing(t *testing.T) {
	_, err := Load("/nonexistent/path/snap.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestToExports(t *testing.T) {
	s := &Snapshot{Env: map[string]string{"KEY": "value"}}
	exports := s.ToExports()
	if len(exports) != 1 {
		t.Fatalf("expected 1 export line, got %d", len(exports))
	}
	if exports[0] != "export KEY=value" {
		t.Errorf("unexpected export line: %q", exports[0])
	}
}
