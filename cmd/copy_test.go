package cmd

import (
	"testing"

	icoppy "envport/internal/copy"
)

// memCopyManager is an in-memory Manager for copy tests.
type memCopyManager struct {
	data map[string]map[string]string
}

func newMemCopyManager() *memCopyManager {
	return &memCopyManager{data: make(map[string]map[string]string)}
}

func (m *memCopyManager) Load(name string) (map[string]string, error) {
	env, ok := m.data[name]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return copy, nil
}

func (m *memCopyManager) Save(name string, env map[string]string) error {
	m.data[name] = env
	return nil
}

func (m *memCopyManager) List() ([]string, error) {
	names := make([]string, 0, len(m.data))
	for k := range m.data {
		names = append(names, k)
	}
	return names, nil
}

func TestCopySnapshot(t *testing.T) {
	m := newMemCopyManager()
	m.data["src"] = map[string]string{"FOO": "bar"}

	if err := icoppy.Run(m, "src", "dst", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, _ := m.Load("dst")
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
}

func TestCopySnapshotMissing(t *testing.T) {
	m := newMemCopyManager()
	err := icoppy.Run(m, "missing", "dst", false)
	if err == nil {
		t.Fatal("expected error for missing src")
	}
}

func TestCopySnapshotNoOverwrite(t *testing.T) {
	m := newMemCopyManager()
	m.data["src"] = map[string]string{"A": "1"}
	m.data["dst"] = map[string]string{"B": "2"}

	err := icoppy.Run(m, "src", "dst", false)
	if err == nil {
		t.Fatal("expected error when dst exists and overwrite=false")
	}
}

func TestCopySnapshotOverwrite(t *testing.T) {
	m := newMemCopyManager()
	m.data["src"] = map[string]string{"A": "1"}
	m.data["dst"] = map[string]string{"B": "2"}

	if err := icoppy.Run(m, "src", "dst", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, _ := m.Load("dst")
	if _, ok := env["B"]; ok {
		t.Error("expected old key B to be gone after overwrite")
	}
}
