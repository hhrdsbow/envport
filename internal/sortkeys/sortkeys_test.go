package sortkeys_test

import (
	"errors"
	"strings"
	"testing"

	"envport/internal/sortkeys"
)

type memManager struct {
	data map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func seed(m *memManager, name string, vars map[string]string) {
	_ = m.Save(name, vars)
}

func TestRunAlpha(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "p", map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"})

	r, err := sortkeys.Run(mgr, "p", sortkeys.StrategyAlpha)
	if err != nil {
		t.Fatal(err)
	}
	if r.Keys[0] != "APPLE" || r.Keys[1] != "MANGO" || r.Keys[2] != "ZEBRA" {
		t.Fatalf("unexpected order: %v", r.Keys)
	}
}

func TestRunReverse(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "p", map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"})

	r, err := sortkeys.Run(mgr, "p", sortkeys.StrategyReverse)
	if err != nil {
		t.Fatal(err)
	}
	if r.Keys[0] != "ZEBRA" || r.Keys[2] != "APPLE" {
		t.Fatalf("unexpected reverse order: %v", r.Keys)
	}
}

func TestRunKeyLen(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "p", map[string]string{"A": "1", "BB": "2", "CCC": "3"})

	r, err := sortkeys.Run(mgr, "p", sortkeys.StrategyKeyLen)
	if err != nil {
		t.Fatal(err)
	}
	if r.Keys[0] != "A" || r.Keys[1] != "BB" || r.Keys[2] != "CCC" {
		t.Fatalf("unexpected keylen order: %v", r.Keys)
	}
}

func TestRunValueLen(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "p", map[string]string{"X": "hi", "Y": "hello world", "Z": "hey"})

	r, err := sortkeys.Run(mgr, "p", sortkeys.StrategyValueLen)
	if err != nil {
		t.Fatal(err)
	}
	if r.Keys[0] != "X" || r.Keys[len(r.Keys)-1] != "Y" {
		t.Fatalf("unexpected valuelen order: %v", r.Keys)
	}
}

func TestRunUnknownStrategy(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "p", map[string]string{"A": "1"})

	_, err := sortkeys.Run(mgr, "p", "bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	mgr := newMemManager()
	_, err := sortkeys.Run(mgr, "missing", sortkeys.StrategyAlpha)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestFormat(t *testing.T) {
	mgr := newMemManager()
	seed(mgr, "p", map[string]string{"KEY": "val"})

	r, _ := sortkeys.Run(mgr, "p", sortkeys.StrategyAlpha)
	out := sortkeys.Format(r)
	if !strings.Contains(out, "KEY=val") {
		t.Fatalf("expected KEY=val in output, got: %s", out)
	}
}
