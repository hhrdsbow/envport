package cmd_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"envport/internal/clone"
	"envport/internal/snapshot"
)

type memCloneManager struct {
	data map[string]*snapshot.Snapshot
}

func newMemCloneManager() *memCloneManager {
	return &memCloneManager{data: make(map[string]*snapshot.Snapshot)}
}

func (m *memCloneManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.data[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return s, nil
}

func (m *memCloneManager) Save(name string, s *snapshot.Snapshot) error {
	m.data[name] = s
	return nil
}

func (m *memCloneManager) List() ([]string, error) {
	out := make([]string, 0, len(m.data))
	for k := range m.data {
		out = append(out, k)
	}
	return out, nil
}

func TestCloneSnapshot(t *testing.T) {
	m := newMemCloneManager()
	m.data["base"] = &snapshot.Snapshot{Name: "base", Vars: map[string]string{"K": "V"}, CreatedAt: time.Now()}

	dest, err := clone.Run(m, "base", "copy", clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "copy" {
		t.Fatalf("expected copy, got %s", dest)
	}
	if m.data["copy"].Vars["K"] != "V" {
		t.Fatal("vars mismatch")
	}
}

func TestCloneOutputMessage(t *testing.T) {
	m := newMemCloneManager()
	m.data["alpha"] = &snapshot.Snapshot{Name: "alpha", Vars: map[string]string{}, CreatedAt: time.Now()}

	var buf bytes.Buffer
	dest, err := clone.Run(m, "alpha", "beta", clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Fprintf(&buf, "cloned %q → %q\n", "alpha", dest)
	if buf.String() != "cloned \"alpha\" → \"beta\"\n" {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}
