package shadow_test

import (
	"errors"
	"testing"

	"envport/internal/shadow"
)

type memManager struct {
	data map[string]map[string]string
}

func (m *memManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func newManager() *memManager {
	return &memManager{data: make(map[string]map[string]string)}
}

func TestRunNoShadows(t *testing.T) {
	mgr := newManager()
	mgr.data["base"] = map[string]string{"A": "1", "B": "2"}
	mgr.data["prod"] = map[string]string{"A": "1", "B": "2"}

	results, err := shadow.Run(mgr, "base", []string{"prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestRunDetectsShadow(t *testing.T) {
	mgr := newManager()
	mgr.data["base"] = map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	mgr.data["prod"] = map[string]string{"DB_HOST": "db.prod.example.com", "PORT": "5432"}

	results, err := shadow.Run(mgr, "base", []string{"prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", results[0].Key)
	}
	if results[0].Source != "prod" {
		t.Errorf("expected source prod, got %s", results[0].Source)
	}
}

func TestRunMissingBase(t *testing.T) {
	mgr := newManager()
	_, err := shadow.Run(mgr, "missing", []string{"prod"})
	if err == nil {
		t.Fatal("expected error for missing base")
	}
}

func TestRunMissingOverride(t *testing.T) {
	mgr := newManager()
	mgr.data["base"] = map[string]string{"A": "1"}
	_, err := shadow.Run(mgr, "base", []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing override")
	}
}

func TestFormat(t *testing.T) {
	results := []shadow.Result{
		{Key: "DB_HOST", BaseVal: "localhost", OverVal: "db.prod.com", Source: "prod"},
	}
	out := shadow.Format(results)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormatEmpty(t *testing.T) {
	out := shadow.Format(nil)
	if out != "no shadowed keys found\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
