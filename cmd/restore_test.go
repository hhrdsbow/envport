package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"envport/internal/profile"
	"envport/internal/snapshot"
)

func TestRestoreCommand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	mgr, err := profile.NewManager(filepath.Join(dir, "profiles"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	snap := snapshot.New(map[string]string{"HELLO": "world"})
	if err := mgr.Save("default", "mysnap", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"restore", "--dry-run", "mysnap"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Error("expected output, got empty")
	}
}

func TestRestoreCommandMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	rootCmd.SetArgs([]string{"restore", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing snapshot")
	}
	_ = os.Unsetenv("ENVPORT_DIR")
}
