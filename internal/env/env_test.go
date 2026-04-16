package env

import (
	"os"
	"testing"
)

func TestCapture(t *testing.T) {
	os.Setenv("ENVPORT_TEST_VAR", "hello")
	t.Cleanup(func() { os.Unsetenv("ENVPORT_TEST_VAR") })

	env := Capture()
	if v, ok := env["ENVPORT_TEST_VAR"]; !ok || v != "hello" {
		t.Fatalf("expected ENVPORT_TEST_VAR=hello, got %q", v)
	}
}

func TestApply(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("ENVPORT_APPLY_A") })

	vars := map[string]string{"ENVPORT_APPLY_A": "42"}
	if err := Apply(vars); err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if got := os.Getenv("ENVPORT_APPLY_A"); got != "42" {
		t.Fatalf("expected 42, got %q", got)
	}
}

func TestDiff(t *testing.T) {
	current := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "1", "B": "99", "D": "4"}

	added, removed, changed := Diff(current, target)

	if len(added) != 1 || added[0] != "D" {
		t.Errorf("expected added=[D], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "C" {
		t.Errorf("expected removed=[C], got %v", removed)
	}
	if len(changed) != 1 || changed[0] != "B" {
		t.Errorf("expected changed=[B], got %v", changed)
	}
}

func TestFilterKeys(t *testing.T) {
	vars := map[string]string{"X": "1", "Y": "2", "Z": "3"}
	result := FilterKeys(vars, []string{"X", "Z", "MISSING"})

	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if result["X"] != "1" || result["Z"] != "3" {
		t.Errorf("unexpected values: %v", result)
	}
	if _, ok := result["MISSING"]; ok {
		t.Error("MISSING key should not be present")
	}
}
