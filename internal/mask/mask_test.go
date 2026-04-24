package mask

import (
	"strings"
	"testing"
)

func TestValueShortString(t *testing.T) {
	opts := DefaultOptions()
	// value shorter than VisibleChars — length must not be leaked
	result := Value("ab", opts)
	if strings.Contains(result, "ab") {
		t.Errorf("short value should be fully masked, got %q", result)
	}
	if len(result) != opts.MaskLen {
		t.Errorf("expected mask length %d, got %d", opts.MaskLen, len(result))
	}
}

func TestValueLongString(t *testing.T) {
	opts := DefaultOptions()
	v := "supersecrettoken"
	result := Value(v, opts)
	if !strings.HasPrefix(result, v[:opts.VisibleChars]) {
		t.Errorf("expected prefix %q in result %q", v[:opts.VisibleChars], result)
	}
	suffix := result[opts.VisibleChars:]
	for _, ch := range suffix {
		if ch != opts.MaskChar {
			t.Errorf("expected mask char %q, got %q", opts.MaskChar, ch)
		}
	}
}

func TestValueCustomMaskChar(t *testing.T) {
	opts := Options{VisibleChars: 2, MaskChar: '#', MaskLen: 6}
	result := Value("hello", opts)
	if !strings.HasPrefix(result, "he") {
		t.Errorf("expected prefix 'he', got %q", result)
	}
	if !strings.HasSuffix(result, "######") {
		t.Errorf("expected six '#' chars, got %q", result)
	}
}

func TestApplyMasksAllValues(t *testing.T) {
	vars := map[string]string{
		"API_KEY": "abc123secret",
		"TOKEN":   "xyz987token",
	}
	opts := DefaultOptions()
	masked := Apply(vars, opts)
	for k, v := range masked {
		original := vars[k]
		if v == original {
			t.Errorf("key %s: value was not masked", k)
		}
	}
	// original map must be unmodified
	if vars["API_KEY"] != "abc123secret" {
		t.Error("Apply mutated the input map")
	}
}

func TestApplyKeysOnlyMasksListed(t *testing.T) {
	vars := map[string]string{
		"SECRET": "topsecret",
		"PUBLIC": "openvalue",
	}
	opts := DefaultOptions()
	masked := ApplyKeys(vars, []string{"SECRET"}, opts)
	if masked["SECRET"] == vars["SECRET"] {
		t.Error("SECRET should have been masked")
	}
	if masked["PUBLIC"] != vars["PUBLIC"] {
		t.Errorf("PUBLIC should be unchanged, got %q", masked["PUBLIC"])
	}
}

func TestApplyKeysEmptyList(t *testing.T) {
	vars := map[string]string{"A": "alpha", "B": "beta"}
	opts := DefaultOptions()
	masked := ApplyKeys(vars, nil, opts)
	for k, v := range masked {
		if v != vars[k] {
			t.Errorf("key %s should be unchanged", k)
		}
	}
}
