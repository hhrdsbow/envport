package reorder_test

import (
	"errors"
	"testing"

	"github.com/envport/envport/internal/reorder"
	"github.com/envport/envport/internal/snapshot"
)

// memManager is an in-memory Manager for testing.
type memManager struct {
	data map[string]*snapshot.Snapshot
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]*snapshot.Snapshot)}
}

func (m *memManager) Load(name string) (*snapshot.Snapshot, error) {
	s, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return s, nil
}

func (m *memManager) Save(name string, s *snapshot.Snapshot) error {
	m.data[name] = s
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	m.data[name] = &snapshot.Snapshot{Vars: vars}
}

func TestRunAlpha(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"})

	keys, err := reorder.Run(m, "dev", reorder.Options{Strategy: reorder.StrategyAlpha})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("expected alpha order, got %v", keys)
	}
}

func TestRunReverse(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"})

	keys, err := reorder.Run(m, "dev", reorder.Options{Strategy: reorder.StrategyReverse})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "ZEBRA" || keys[1] != "MANGO" || keys[2] != "APPLE" {
		t.Errorf("expected reverse order, got %v", keys)
	}
}

func TestRunCustomOrder(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3", "BANANA": "4"})

	keys, err := reorder.Run(m, "dev", reorder.Options{
		Strategy:    reorder.StrategyCustom,
		CustomOrder: []string{"MANGO", "ZEBRA"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "MANGO" || keys[1] != "ZEBRA" {
		t.Errorf("pinned keys in wrong position: %v", keys)
	}
	// remaining keys should be alpha-sorted
	if keys[2] != "APPLE" || keys[3] != "BANANA" {
		t.Errorf("remainder not alpha-sorted: %v", keys)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newMemManager()
	_, err := reorder.Run(m, "ghost", reorder.Options{Strategy: reorder.StrategyAlpha})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunEmptyName(t *testing.T) {
	m := newMemManager()
	_, err := reorder.Run(m, "", reorder.Options{Strategy: reorder.StrategyAlpha})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRunUnknownStrategy(t *testing.T) {
	m := newMemManager()
	seed(m, "dev", map[string]string{"A": "1"})
	_, err := reorder.Run(m, "dev", reorder.Options{Strategy: "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}
