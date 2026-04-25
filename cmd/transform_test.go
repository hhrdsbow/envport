package cmd_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/yourorg/envport/internal/transform"
)

// memTransformManager satisfies transform.Manager for cmd-level tests.
type memTransformManager struct {
	data map[string]map[string]string
}

func newMemTransformManager() *memTransformManager {
	return &memTransformManager{data: make(map[string]map[string]string)}
}

func (m *memTransformManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (m *memTransformManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func TestTransformPrefixAdd(t *testing.T) {
	m := newMemTransformManager()
	m.data["dev"] = map[string]string{"HOST": "localhost", "PORT": "9000"}

	ops := []transform.Op{{Kind: "prefix-add", Value: "SVC_"}}
	r, out, err := transform.Run(m, "dev", ops, false)
	if err != nil {
		t.Fatal(err)
	}
	if r.Changed != 2 {
		t.Errorf("expected 2 changed, got %d", r.Changed)
	}
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", out["SVC_HOST"])
	}
}

func TestTransformFormatOutput(t *testing.T) {
	var buf bytes.Buffer
	r := transform.Result{Changed: 1, Total: 3}
	buf.WriteString(transform.Format(r, false))
	if !bytes.Contains(buf.Bytes(), []byte("1/3")) {
		t.Errorf("expected '1/3' in output, got %q", buf.String())
	}
}

func TestTransformDryRunNoSave(t *testing.T) {
	m := newMemTransformManager()
	m.data["prod"] = map[string]string{"DB": "postgres"}

	ops := []transform.Op{{Kind: "lowercase"}}
	_, _, err := transform.Run(m, "prod", ops, true)
	if err != nil {
		t.Fatal(err)
	}
	// original key must survive
	if _, ok := m.data["prod"]["DB"]; !ok {
		t.Error("dry-run should not modify stored data")
	}
}
