package normalize

import (
	"strings"
	"testing"
)

func TestRunUpperKeys(t *testing.T) {
	vars := map[string]string{"path": "/usr/bin", "home": "/root"}
	opts := DefaultOptions()
	opts.TrimValues = false
	opts.RemoveEmpty = false

	r := Run(vars, opts)

	if _, ok := r.Vars["PATH"]; !ok {
		t.Error("expected PATH key after upper-casing")
	}
	if _, ok := r.Vars["HOME"]; !ok {
		t.Error("expected HOME key after upper-casing")
	}
	if !r.Changed {
		t.Error("expected Changed=true when keys were renamed")
	}
}

func TestRunTrimValues(t *testing.T) {
	vars := map[string]string{"KEY": "  hello  ", "OTHER": "clean"}
	opts := Options{UpperKeys: false, TrimValues: true, RemoveEmpty: false}

	r := Run(vars, opts)

	if r.Vars["KEY"] != "hello" {
		t.Errorf("expected trimmed value 'hello', got %q", r.Vars["KEY"])
	}
	if r.Vars["OTHER"] != "clean" {
		t.Errorf("expected unchanged value 'clean', got %q", r.Vars["OTHER"])
	}
	if !r.Changed {
		t.Error("expected Changed=true when value was trimmed")
	}
}

func TestRunRemoveEmpty(t *testing.T) {
	vars := map[string]string{"PRESENT": "yes", "EMPTY": ""}
	opts := Options{UpperKeys: false, TrimValues: false, RemoveEmpty: true}

	r := Run(vars, opts)

	if _, ok := r.Vars["EMPTY"]; ok {
		t.Error("expected EMPTY key to be removed")
	}
	if r.Vars["PRESENT"] != "yes" {
		t.Errorf("expected PRESENT to remain, got %q", r.Vars["PRESENT"])
	}
	if !r.Changed {
		t.Error("expected Changed=true when empty key removed")
	}
}

func TestRunNoChange(t *testing.T) {
	vars := map[string]string{"KEY": "value"}
	opts := Options{UpperKeys: false, TrimValues: false, RemoveEmpty: false}

	r := Run(vars, opts)

	if r.Changed {
		t.Error("expected Changed=false when nothing changed")
	}
	if r.Vars["KEY"] != "value" {
		t.Errorf("unexpected value: %q", r.Vars["KEY"])
	}
}

func TestFormatSummary(t *testing.T) {
	before := map[string]string{"path": "/usr/bin", "empty": ""}
	opts := DefaultOptions()
	opts.RemoveEmpty = true
	r := Run(before, opts)

	out := Format(before, r.Vars)
	if !strings.Contains(out, "removed") && !strings.Contains(out, "renamed") && out != "" {
		t.Errorf("unexpected format output: %q", out)
	}
}
