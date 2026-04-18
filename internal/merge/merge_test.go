package merge_test

import (
	"errors"
	"testing"

	"github.com/user/envport/internal/merge"
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

func TestRunMergeNoConflicts(t *testing.T) {
	mgr := newMemManager()
	mgr.data["base"] = map[string]string{"A": "1", "B": "2"}
	mgr.data["other"] = map[string]string{"C": "3"}

	res, err := merge.Run(mgr, "base", "other", "dest", merge.StrategyBase)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if res.Vars["A"] != "1" || res.Vars["B"] != "2" || res.Vars["C"] != "3" {
		t.Errorf("unexpected vars: %v", res.Vars)
	}
}

func TestRunMergeStrategyBase(t *testing.T) {
	mgr := newMemManager()
	mgr.data["base"] = map[string]string{"A": "base"}
	mgr.data["other"] = map[string]string{"A": "other"}

	res, err := merge.Run(mgr, "base", "other", "dest", merge.StrategyBase)
	if err != nil {
		t.Fatal(err)
	}
	if res.Vars["A"] != "base" {
		t.Errorf("expected base value, got %q", res.Vars["A"])
	}
}

func TestRunMergeStrategyOther(t *testing.T) {
	mgr := newMemManager()
	mgr.data["base"] = map[string]string{"A": "base"}
	mgr.data["other"] = map[string]string{"A": "other"}

	res, err := merge.Run(mgr, "base", "other", "dest", merge.StrategyOther)
	if err != nil {
		t.Fatal(err)
	}
	if res.Vars["A"] != "other" {
		t.Errorf("expected other value, got %q", res.Vars["A"])
	}
}

func TestRunMergeStrategyError(t *testing.T) {
	mgr := newMemManager()
	mgr.data["base"] = map[string]string{"A": "base"}
	mgr.data["other"] = map[string]string{"A": "other"}

	_, err := merge.Run(mgr, "base", "other", "dest", merge.StrategyError)
	if err == nil {
		t.Fatal("expected error on conflict")
	}
}

func TestRunMergeMissingBase(t *testing.T) {
	mgr := newMemManager()
	mgr.data["other"] = map[string]string{"A": "1"}

	_, err := merge.Run(mgr, "base", "other", "dest", merge.StrategyBase)
	if err == nil {
		t.Fatal("expected error for missing base")
	}
}
