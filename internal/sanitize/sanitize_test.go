package sanitize_test

import (
	"testing"

	"envport/internal/sanitize"
)

func TestTrimSpace(t *testing.T) {
	input := map[string]string{
		"FOO": "  hello  ",
		"BAR": "no-spaces",
	}
	opts := sanitize.DefaultOptions()
	out, res := sanitize.Run(input, opts)

	if out["FOO"] != "hello" {
		t.Fatalf("expected trimmed value, got %q", out["FOO"])
	}
	if out["BAR"] != "no-spaces" {
		t.Fatalf("unexpected change to BAR: %q", out["BAR"])
	}
	if len(res.TrimmedKeys) != 1 || res.TrimmedKeys[0] != "FOO" {
		t.Fatalf("unexpected TrimmedKeys: %v", res.TrimmedKeys)
	}
}

func TestRemoveEmpty(t *testing.T) {
	input := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
		"SPACES":  "   ",
	}
	opts := sanitize.Options{TrimSpace: true, RemoveEmpty: true}
	out, res := sanitize.Run(input, opts)

	if _, ok := out["EMPTY"]; ok {
		t.Fatal("EMPTY should have been dropped")
	}
	if _, ok := out["SPACES"]; ok {
		t.Fatal("SPACES should have been dropped after trim")
	}
	if out["PRESENT"] != "value" {
		t.Fatalf("PRESENT should be retained")
	}
	if len(res.DroppedKeys) != 2 {
		t.Fatalf("expected 2 dropped keys, got %d", len(res.DroppedKeys))
	}
}

func TestUpperKeys(t *testing.T) {
	input := map[string]string{
		"my_var": "value",
		"ALREADY": "up",
	}
	opts := sanitize.Options{UpperKeys: true}
	out, res := sanitize.Run(input, opts)

	if _, ok := out["MY_VAR"]; !ok {
		t.Fatal("expected MY_VAR in output")
	}
	if len(res.RenamedKeys) != 1 {
		t.Fatalf("expected 1 renamed key, got %d", len(res.RenamedKeys))
	}
}

func TestStripPrefix(t *testing.T) {
	input := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"OTHER":    "x",
	}
	opts := sanitize.Options{StripPrefix: "APP_"}
	out, _ := sanitize.Run(input, opts)

	if _, ok := out["HOST"]; !ok {
		t.Fatal("expected HOST after prefix strip")
	}
	if _, ok := out["PORT"]; !ok {
		t.Fatal("expected PORT after prefix strip")
	}
	if _, ok := out["OTHER"]; !ok {
		t.Fatal("OTHER should remain unchanged")
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Fatal("APP_HOST should have been stripped")
	}
}

func TestResultChanged(t *testing.T) {
	input := map[string]string{"KEY": "val"}
	_, res := sanitize.Run(input, sanitize.DefaultOptions())
	if res.Changed() {
		t.Fatal("expected no changes")
	}
}
