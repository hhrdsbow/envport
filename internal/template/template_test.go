package template_test

import (
	"errors"
	"strings"
	"testing"

	"envport/internal/template"
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

func TestRunSubstitution(t *testing.T) {
	mgr := newMemManager()
	mgr.data["src"] = map[string]string{
		"DB_HOST": "{{HOST}}.example.com",
		"APP_ENV": "{{ENV}}",
		"STATIC":  "unchanged",
	}

	res, err := template.Run(mgr, "src", "dst", map[string]string{
		"HOST": "db-prod",
		"ENV":  "production",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Applied != 2 {
		t.Errorf("expected 2 applied, got %d", res.Applied)
	}
	if res.Vars["DB_HOST"] != "db-prod.example.com" {
		t.Errorf("unexpected DB_HOST: %s", res.Vars["DB_HOST"])
	}
	if res.Vars["STATIC"] != "unchanged" {
		t.Errorf("static value changed: %s", res.Vars["STATIC"])
	}
}

func TestRunMissingSrc(t *testing.T) {
	mgr := newMemManager()
	_, err := template.Run(mgr, "missing", "dst", nil)
	if err == nil {
		t.Fatal("expected error for missing src")
	}
}

func TestRunSavesDst(t *testing.T) {
	mgr := newMemManager()
	mgr.data["src"] = map[string]string{"KEY": "val"}
	_, err := template.Run(mgr, "src", "dst", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := mgr.data["dst"]; !ok {
		t.Error("dst snapshot not saved")
	}
}

func TestFormat(t *testing.T) {
	r := template.Result{Name: "prod", Applied: 3}
	out := template.Format(r)
	if !strings.Contains(out, "prod") || !strings.Contains(out, "3") {
		t.Errorf("unexpected format output: %s", out)
	}
}
