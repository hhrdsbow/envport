package validate_test

import (
	"errors"
	"strings"
	"testing"

	"envport/internal/validate"
)

type memManager struct {
	data map[string]map[string]string
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func newManager() *memManager {
	return &memManager{data: map[string]map[string]string{
		"prod": {
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
			"API_KEY": "",
		},
	}}
}

func TestRunAllPresent(t *testing.T) {
	mgr := newManager()
	results, err := validate.Run(mgr, "prod", []string{"DB_HOST", "DB_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no issues, got %d", len(results))
	}
}

func TestRunMissingKey(t *testing.T) {
	mgr := newManager()
	results, err := validate.Run(mgr, "prod", []string{"DB_HOST", "SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "SECRET" || results[0].Message != "missing" {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestRunEmptyValue(t *testing.T) {
	mgr := newManager()
	results, err := validate.Run(mgr, "prod", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Message != "empty" {
		t.Fatalf("expected empty result, got %+v", results)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	mgr := newManager()
	_, err := validate.Run(mgr, "ghost", []string{"KEY"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	results := []validate.Result{
		{Key: "FOO", Message: "missing"},
		{Key: "BAR", Message: "empty"},
	}
	out := validate.Format(results)
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Fatalf("unexpected format output: %s", out)
	}
}

func TestFormatNoIssues(t *testing.T) {
	out := validate.Format(nil)
	if !strings.Contains(out, "all required") {
		t.Fatalf("unexpected output: %s", out)
	}
}
