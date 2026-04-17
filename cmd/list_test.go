package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envport/internal/profile"
	"envport/internal/snapshot"
)

func TestListCommand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_STORE_DIR", dir)

	mgr, err := profile.New(dir)
	if err != nil {
		t.Fatalf("profile.New: %v", err)
	}

	snap := snapshot.New(map[string]string{"FOO": "bar"})
	if err := mgr.Save("alpha", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := mgr.Save("beta", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected 'alpha' in output, got: %s", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected 'beta' in output, got: %s", out)
	}
}

func TestListCommandEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_STORE_DIR", dir)

	_ = filepath.Join(dir, "unused")
	_ = os.MkdirAll(dir, 0755)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No snapshots") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
