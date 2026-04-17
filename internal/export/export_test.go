package export

import (
	"strings"
	"testing"
)

var testEnv = map[string]string{
	"FOO": "bar",
	"BAZ": "qux",
}

func TestWriteShell(t *testing.T) {
	var b strings.Builder
	if err := Write(&b, testEnv, FormatShell); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, "export BAZ") || !strings.Contains(out, "export FOO") {
		t.Errorf("unexpected shell output: %s", out)
	}
}

func TestWriteDotenv(t *testing.T) {
	var b strings.Builder
	if err := Write(&b, testEnv, FormatDotenv); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("unexpected dotenv output: %s", out)
	}
}

func TestWriteJSON(t *testing.T) {
	var b strings.Builder
	if err := Write(&b, testEnv, FormatJSON); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, `"FOO"`) || !strings.Contains(out, `"bar"`) {
		t.Errorf("unexpected json output: %s", out)
	}
}

func TestWriteUnknownFormat(t *testing.T) {
	var b strings.Builder
	err := Write(&b, testEnv, Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWriteOrdering(t *testing.T) {
	var b strings.Builder
	_ = Write(&b, testEnv, FormatDotenv)
	lines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(lines) != 2 || !strings.HasPrefix(lines[0], "BAZ") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
