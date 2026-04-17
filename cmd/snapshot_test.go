package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotCommand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DATA_DIR", dir)

	t.Setenv("MY_VAR", "hello")
	t.Setenv("OTHER_VAR", "world")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"snapshot", "testprofile"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "testprofile") {
		t.Errorf("expected profile name in output, got: %s", output)
	}

	entries, err := os.ReadDir(filepath.Join(dir, "profiles"))
	if err != nil {
		t.Fatalf("profiles dir not created: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected at least one profile file")
	}
}

func TestSnapshotCommandFilterKeys(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DATA_DIR", dir)

	t.Setenv("KEEP_ME", "yes")
	t.Setenv("IGNORE_ME", "no")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"snapshot", "filtered", "--keys", "KEEP_ME"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "1 variables") {
		t.Errorf("expected 1 variable captured, output: %s", buf.String())
	}
}
