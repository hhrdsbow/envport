package import_test

import (
	"strings"
	"testing"

	importer "envport/internal/import"
)

func TestParseShell(t *testing.T) {
	input := `
# comment
export FOO=bar
export BAZ="hello world"
PLAIN=value
`
	got, err := importer.Parse(strings.NewReader(input), importer.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{
		"FOO":   "bar",
		"BAZ":   "hello world",
		"PLAIN": "value",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestParseDotenv(t *testing.T) {
	input := "DB_HOST=localhost\nDB_PORT='5432'\n# ignore me\nSECRET=\"abc123\""
	got, err := importer.Parse(strings.NewReader(input), importer.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: got %q", got["DB_PORT"])
	}
	if got["SECRET"] != "abc123" {
		t.Errorf("SECRET: got %q", got["SECRET"])
	}
}

func TestParseJSON(t *testing.T) {
	input := `{"APP_ENV": "production", "PORT": "8080"}`
	got, err := importer.Parse(strings.NewReader(input), importer.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q", got["APP_ENV"])
	}
	if got["PORT"] != "8080" {
		t.Errorf("PORT: got %q", got["PORT"])
	}
}

func TestParseUnknownFormat(t *testing.T) {
	_, err := importer.Parse(strings.NewReader(""), importer.Format("yaml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestParseJSONInvalid(t *testing.T) {
	_, err := importer.Parse(strings.NewReader("not json"), importer.FormatJSON)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
