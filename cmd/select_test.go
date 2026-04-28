package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"

	envselect "envport/internal/select"
)

type memSelectManager struct {
	data map[string]map[string]string
}

func newMemSelectManager() *memSelectManager {
	return &memSelectManager{data: make(map[string]map[string]string)}
}

func (m *memSelectManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func (m *memSelectManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func buildSelectCmd(m *memSelectManager) *cobra.Command {
	var patterns []string
	var dst string

	cmd := &cobra.Command{
		Use: "select <snapshot>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := envselect.Run(m, args[0], dst, patterns)
			if err != nil {
				return err
			}
			cmd.Print(envselect.Format(r))
			return nil
		},
	}
	cmd.Flags().StringVarP(&dst, "out", "o", "", "")
	cmd.Flags().StringArrayVarP(&patterns, "key", "k", nil, "")
	return cmd
}

func TestSelectCommandSuccess(t *testing.T) {
	m := newMemSelectManager()
	m.data["src"] = map[string]string{"FOO": "bar", "BAZ": "qux", "OTHER": "val"}

	buf := &bytes.Buffer{}
	cmd := buildSelectCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"src", "-o", "dst", "-k", "FOO", "-k", "BAZ"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Error("expected output from select command")
	}
	if _, ok := m.data["dst"]["FOO"]; !ok {
		t.Error("expected FOO in dst")
	}
	if _, ok := m.data["dst"]["OTHER"]; ok {
		t.Error("OTHER should not be in dst")
	}
}

func TestSelectCommandMissingSource(t *testing.T) {
	m := newMemSelectManager()
	cmd := buildSelectCmd(m)
	cmd.SetArgs([]string{"missing", "-o", "dst", "-k", "A"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing source")
	}
}
