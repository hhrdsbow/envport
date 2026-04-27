package cmd_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"envport/internal/unset"

	"github.com/spf13/cobra"
)

// memUnsetManager is a minimal in-memory Manager for cmd-level tests.
type memUnsetManager struct {
	data map[string]map[string]string
}

func newMemUnsetManager() *memUnsetManager {
	return &memUnsetManager{data: make(map[string]map[string]string)}
}

func (m *memUnsetManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (m *memUnsetManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func buildUnsetCmd(m *memUnsetManager, strict bool) (*cobra.Command, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{
		Use: "unset",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := unset.Run(m, args[0], args[1:], strict)
			if err != nil {
				return err
			}
			cmd.Println(unset.Format(res))
			return nil
		},
	}
	cmd.SetOut(buf)
	return cmd, buf
}

func TestUnsetCommandSuccess(t *testing.T) {
	m := newMemUnsetManager()
	m.data["dev"] = map[string]string{"FOO": "bar", "BAZ": "qux"}

	cmd, buf := buildUnsetCmd(m, false)
	cmd.SetArgs([]string{"dev", "FOO"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "FOO") {
		t.Errorf("expected FOO in output, got: %s", buf.String())
	}
	if _, ok := m.data["dev"]["FOO"]; ok {
		t.Error("FOO should have been removed")
	}
}

func TestUnsetCommandMissingSnapshot(t *testing.T) {
	m := newMemUnsetManager()
	cmd, _ := buildUnsetCmd(m, false)
	cmd.SetArgs([]string{"ghost", "KEY"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
