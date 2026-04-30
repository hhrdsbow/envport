package slice_test

import (
	"errors"
	"testing"

	"envport/internal/slice"
)

// --- in-memory manager ---

type memManager struct {
	data map[string]map[string]string
}

func newMemManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (m *memManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func seed(m *memManager, name string) {
	m.data[name] = map[string]string{
		"ALPHA": "1",
		"BETA":  "2",
		"GAMMA": "3",
		"DELTA": "4",
		"ECHO":  "5",
	}
}

func TestRunSlicesMiddle(t *testing.T) {
	m := newMemManager()
	seed(m, "src")

	// sorted: ALPHA BETA DELTA ECHO GAMMA  → indices 1..3 = BETA DELTA ECHO
	r, err := slice.Run(m, "src", "dst", 1, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
	if r.Keys[0] != "BETA" || r.Keys[1] != "DELTA" || r.Keys[2] != "ECHO" {
		t.Errorf("unexpected keys: %v", r.Keys)
	}
	if _, saved := m.data["dst"]; !saved {
		t.Error("dst snapshot not saved")
	}
}

func TestRunSlicesAll(t *testing.T) {
	m := newMemManager()
	seed(m, "src")

	r, err := slice.Run(m, "src", "", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(r.Keys))
	}
}

func TestRunSlicesMissingSource(t *testing.T) {
	m := newMemManager()
	_, err := slice.Run(m, "missing", "", 0, 2)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRunSlicesFromGTTo(t *testing.T) {
	m := newMemManager()
	seed(m, "src")
	_, err := slice.Run(m, "src", "", 3, 1)
	if err == nil {
		t.Fatal("expected error when from > to")
	}
}

func TestFormat(t *testing.T) {
	r := slice.Result{Keys: []string{"A", "B"}, Src: "base", Dst: "out"}
	s := slice.Format(r)
	if s == "" {
		t.Error("expected non-empty format string")
	}
}

func TestFormatEmpty(t *testing.T) {
	r := slice.Result{}
	s := slice.Format(r)
	if s == "" {
		t.Error("expected non-empty format string for empty result")
	}
}
