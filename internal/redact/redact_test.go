package redact_test

import (
	"testing"

	"envport/internal/redact"
)

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"AWS_SECRET_KEY", true},
		{"GITHUB_TOKEN", true},
		{"API_KEY", true},
		{"PRIVATE_KEY", true},
		{"AUTH_HEADER", true},
		{"HOME", false},
		{"PATH", false},
		{"PORT", false},
		{"DATABASE_URL", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := redact.IsSensitive(tc.key, redact.DefaultPatterns)
			if got != tc.want {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestApply(t *testing.T) {
	vars := map[string]string{
		"HOME":        "/home/user",
		"DB_PASSWORD": "s3cr3t",
		"GITHUB_TOKEN": "ghp_abc123",
		"PORT":        "8080",
	}

	result := redact.ApplyDefault(vars)

	if result["HOME"] != "/home/user" {
		t.Errorf("HOME should not be redacted, got %q", result["HOME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be redacted, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["GITHUB_TOKEN"] != "***" {
		t.Errorf("GITHUB_TOKEN should be redacted, got %q", result["GITHUB_TOKEN"])
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	vars := map[string]string{
		"API_KEY": "original",
	}
	_ = redact.ApplyDefault(vars)
	if vars["API_KEY"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApplyCustomPatterns(t *testing.T) {
	vars := map[string]string{
		"MY_CUSTOM_SENSITIVE": "value",
		"NORMAL_VAR":         "hello",
	}
	result := redact.Apply(vars, []string{"CUSTOM"})
	if result["MY_CUSTOM_SENSITIVE"] != "***" {
		t.Errorf("expected redaction, got %q", result["MY_CUSTOM_SENSITIVE"])
	}
	if result["NORMAL_VAR"] != "hello" {
		t.Errorf("expected no redaction, got %q", result["NORMAL_VAR"])
	}
}
