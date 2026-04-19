package cmd

import (
	"bytes"
	"strings"
	"testing"

	"envport/internal/template"
)

type memTemplateManager struct {
	data map[string]map[string]string
}

func newMemTemplateManager() *memTemplateManager {
	return &memTemplateManager{data: make(map[string]map[string]string)}
}

func (m *memTemplateManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return v, nil
}

func (m *memTemplateManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func TestTemplateSubstitution(t *testing.T) {
	mgr := newMemTemplateManager()
	mgr.data["base"] = map[string]string{
		"URL": "https://{{HOST}}/api",
	}

	res, err := template.Run(mgr, "base", "prod", map[string]string{"HOST": "prod.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["URL"] != "https://prod.example.com/api" {
		t.Errorf("wrong URL: %s", res.Vars["URL"])
	}
}

func TestTemplateFormatOutput(t *testing.T) {
	r := template.Result{Name: "staging", Applied: 5}
	out := template.Format(r)
	var buf bytes.Buffer
	buf.WriteString(out)
	if !strings.Contains(buf.String(), "staging") {
		t.Error("format output missing snapshot name")
	}
	if !strings.Contains(buf.String(), "5") {
		t.Error("format output missing substitution count")
	}
}
