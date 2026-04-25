package cmd_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"envport/internal/patch"

	"github.com/spf13/cobra"
)

// memPatchManager satisfies patch.Manager in tests.
type memPatchManager struct {
	store map[string]map[string]string
}

func newMemPatchManager() *memPatchManager {
	return &memPatchManager{store: make(map[string]map[string]string)}
}

func (m *memPatchManager) Load(name string) (map[string]string, error) {
	v, ok := m.store[name]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (m *memPatchManager) Save(name string, vars map[string]string) error {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	m.store[name] = copy
	return nil
}

func buildPatchCmd(m *memPatchManager) (*cobra.Command, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cmd := &cobra.Command{
		Use: "patch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("need profile and at least one key=value")
			}
			updates := make(map[string]string)
			for _, pair := range args[1:] {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) != 2 {
					return errors.New("bad pair")
				}
				updates[parts[0]] = parts[1]
			}
			r, err := patch.Run(m, args[0], updates)
			if err != nil {
				return err
			}
			cmd.Println(patch.Format(r))
			return nil
		},
	}
	cmd.SetOut(out)
	return cmd, out
}

func TestPatchCommandSuccess(t *testing.T) {
	m := newMemPatchManager()
	_ = m.Save("dev", map[string]string{"HOST": "localhost", "PORT": "8080"})

	cmd, out := buildPatchCmd(m)
	cmd.SetArgs([]string{"dev", "PORT=9090", "DEBUG=true"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "patched") {
		t.Errorf("expected output to contain 'patched', got: %q", got)
	}
	if !strings.Contains(got, "dev") {
		t.Errorf("expected output to contain profile name, got: %q", got)
	}
}

func TestPatchCommandMissingProfile(t *testing.T) {
	m := newMemPatchManager()
	cmd, _ := buildPatchCmd(m)
	cmd.SetArgs([]string{"ghost", "X=1"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestPatchCommandPreservesUntouched(t *testing.T) {
	m := newMemPatchManager()
	_ = m.Save("staging", map[string]string{"A": "1", "B": "2"})

	cmd, _ := buildPatchCmd(m)
	cmd.SetArgs([]string{"staging", "A=99"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vars, _ := m.Load("staging")
	if vars["B"] != "2" {
		t.Errorf("expected B=2 preserved, got %q", vars["B"])
	}
	if vars["A"] != "99" {
		t.Errorf("expected A=99, got %q", vars["A"])
	}
}
