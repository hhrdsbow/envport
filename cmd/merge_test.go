package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envport/internal/profile"
)

func setupMergeProfiles(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	mgr, err := profile.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.Save("base", map[string]string{"A": "1", "B": "2"}); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Save("other", map[string]string{"B": "99", "C": "3"}); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestMergeCommandStrategyBase(t *testing.T) {
	dir := setupMergeProfiles(t)

	rootCmd.SetArgs([]string{"merge", "base", "other", "merged", "--strategy", "base"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	mgr, _ := profile.New(dir)
	vars, err := mgr.Load("merged")
	if err != nil {
		t.Fatal(err)
	}
	if vars["B"] != "2" {
		t.Errorf("base strategy: expected B=2, got %q", vars["B"])
	}
	if vars["C"] != "3" {
		t.Errorf("expected C=3, got %q", vars["C"])
	}
}

func TestMergeCommandStrategyOther(t *testing.T) {
	dir := setupMergeProfiles(t)

	rootCmd.SetArgs([]string{"merge", "base", "other", "merged2", "--strategy", "other"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	mgr, _ := profile.New(dir)
	vars, err := mgr.Load("merged2")
	if err != nil {
		t.Fatal(err)
	}
	if vars["B"] != "99" {
		t.Errorf("other strategy: expected B=99, got %q", vars["B"])
	}
}

func TestMergeCommandInvalidStrategy(t *testing.T) {
	setupMergeProfiles(t)

	rootCmd.SetArgs([]string{"merge", "base", "other", "dest", "--strategy", "invalid"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func init() {
	_ = os.MkdirAll(filepath.Join(os.TempDir(), "envport-test"), 0o755)
}
