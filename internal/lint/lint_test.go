package lint

import (
	"errors"
	"testing"
)

// memManager is an in-memory Manager for testing.
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

func newManager(vars map[string]string) *memManager {
	return &memManager{data: map[string]map[string]string{"snap": vars}}
}

func TestRunClean(t *testing.T) {
	m := newManager(map[string]string{"HOME": "/home/user", "PATH": "/usr/bin"})
	issues, err := Run(m, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestRunEmptyValue(t *testing.T) {
	m := newManager(map[string]string{"TOKEN": "", "HOST": "localhost"})
	issues, err := Run(m, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 || issues[0].Key != "TOKEN" {
		t.Errorf("expected 1 empty-value issue for TOKEN, got %v", issues)
	}
}

func TestRunKeyWithSpace(t *testing.T) {
	m := newManager(map[string]string{"MY KEY": "val"})
	issues, err := Run(m, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Error("expected whitespace issue, got none")
	}
}

func TestRunKeyWithEquals(t *testing.T) {
	m := newManager(map[string]string{"BAD=KEY": "val"})
	issues, err := Run(m, "snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Error("expected equals issue, got none")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newManager(map[string]string{})
	_, err := Run(m, "missing")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	issues := []Issue{{Key: "TOKEN", Message: "empty value"}}
	out := Format(issues)
	if out == "no issues found" {
		t.Error("expected formatted issue, got no-issues message")
	}
}

func TestFormatEmpty(t *testing.T) {
	out := Format(nil)
	if out != "no issues found" {
		t.Errorf("expected no-issues message, got %q", out)
	}
}
