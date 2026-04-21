package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envport/internal/alias"
)

type memAliasStore struct {
	data map[string]string
}

func newMemAliasStore() *memAliasStore {
	return &memAliasStore{data: make(map[string]string)}
}
func (s *memAliasStore) Get(k string) (string, error) {
	v, ok := s.data[k]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (s *memAliasStore) Set(k, v string) error { s.data[k] = v; return nil }
func (s *memAliasStore) Delete(k string) error  { delete(s.data, k); return nil }
func (s *memAliasStore) List() (map[string]string, error) {
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out, nil
}

func buildAliasCmd(m *alias.Manager) *cobra.Command {
	root := &cobra.Command{Use: "envport"}
	ac := &cobra.Command{Use: "alias"}
	ac.AddCommand(
		&cobra.Command{
			Use:  "add <alias> <snapshot>",
			Args: cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := m.Add(args[0], args[1]); err != nil {
					return err
				}
				cmd.Printf("alias %q → %q saved\n", args[0], args[1])
				return nil
			},
		},
		&cobra.Command{
			Use:  "list",
			Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				list, err := m.List()
				if err != nil {
					return err
				}
				for a, t := range list {
					cmd.Printf("%s\t%s\n", a, t)
				}
				return nil
			},
		},
	)
	root.AddCommand(ac)
	return root
}

func TestAliasAdd(t *testing.T) {
	m := alias.New(newMemAliasStore())
	root := buildAliasCmd(m)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"alias", "add", "prod", "snap-prod"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(buf.String(), "prod") {
		t.Errorf("expected output to mention alias name, got: %s", buf.String())
	}
}

func TestAliasList(t *testing.T) {
	m := alias.New(newMemAliasStore())
	_ = m.Add("staging", "snap-staging")
	root := buildAliasCmd(m)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"alias", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(buf.String(), "staging") {
		t.Errorf("expected staging in list output, got: %s", buf.String())
	}
}
