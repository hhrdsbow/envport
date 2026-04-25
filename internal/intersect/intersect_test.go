package intersect_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/intersect"
)

// memManager is an in-memory Manager for tests.
type memManager struct {
	snaps map[string]map[string]string
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.snaps[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func newManager() *memManager {
	return &memManager{snaps: make(map[string]map[string]string)}
}

func seed(m *memManager, name string, vars map[string]string) {
	m.snaps[name] = vars
}

func TestRunCommonKeys(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"})
	seed(m, "b", map[string]string{"FOO": "10", "BAR": "20", "QUX": "4"})

	res, err := intersect.Run(m, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 2 {
		t.Fatalf("expected 2 common keys, got %d: %v", len(res.Keys), res.Keys)
	}
	if res.Keys[0] != "BAR" || res.Keys[1] != "FOO" {
		t.Errorf("unexpected key order: %v", res.Keys)
	}
	if res.Values["FOO"]["a"] != "1" || res.Values["FOO"]["b"] != "10" {
		t.Errorf("unexpected FOO values: %v", res.Values["FOO"])
	}
}

func TestRunNoCommonKeys(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"ALPHA": "1"})
	seed(m, "b", map[string]string{"BETA": "2"})

	res, err := intersect.Run(m, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 0 {
		t.Errorf("expected no common keys, got %v", res.Keys)
	}
}

func TestRunThreeSnapshots(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"X": "1", "Y": "2", "Z": "3"})
	seed(m, "b", map[string]string{"X": "10", "Y": "20"})
	seed(m, "c", map[string]string{"X": "100", "Y": "200", "Z": "300"})

	res, err := intersect.Run(m, []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Z is missing from "b", so only X and Y should be common.
	if len(res.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(res.Keys), res.Keys)
	}
}

func TestRunMissingSnapshot(t *testing.T) {
	m := newManager()
	seed(m, "a", map[string]string{"FOO": "1"})

	_, err := intersect.Run(m, []string{"a", "missing"})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunTooFewNames(t *testing.T) {
	m := newManager()
	_, err := intersect.Run(m, []string{"only-one"})
	if err == nil {
		t.Fatal("expected error when fewer than two names provided")
	}
}
