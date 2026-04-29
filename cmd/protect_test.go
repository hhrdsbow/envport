package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"envport/internal/protect"

	"github.com/spf13/cobra"
)

type memProtectStore struct {
	data map[string]struct{}
}

func newMemProtectStore() *memProtectStore {
	return &memProtectStore{data: make(map[string]struct{})}
}

func (s *memProtectStore) Set(name string) error {
	s.data[name] = struct{}{}
	return nil
}
func (s *memProtectStore) Delete(name string) error {
	delete(s.data, name)
	return nil
}
func (s *memProtectStore) Exists(name string) (bool, error) {
	_, ok := s.data[name]
	return ok, nil
}
func (s *memProtectStore) List() ([]string, error) {
	out := make([]string, 0, len(s.data))
	for k := range s.data {
		out = append(out, k)
	}
	return out, nil
}

func buildProtectCmd(m *protect.Manager) *cobra.Command {
	root := &cobra.Command{Use: "envport"}
	protectCmd := &cobra.Command{Use: "protect"}

	protectCmd.AddCommand(&cobra.Command{
		Use:  "add",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Protect(args[0])
		},
	})
	protectCmd.AddCommand(&cobra.Command{
		Use:  "remove",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Unprotect(args[0])
		},
	})
	protectCmd.AddCommand(&cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := m.List()
			if err != nil {
				return err
			}
			cmd.Println(strings.Join(names, "\n"))
			return nil
		},
	})
	root.AddCommand(protectCmd)
	return root
}

func TestProtectCommandAdd(t *testing.T) {
	m := protect.New(newMemProtectStore())
	root := buildProtectCmd(m)
	root.SetArgs([]string{"protect", "add", "prod"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, _ := m.IsProtected("prod")
	if !ok {
		t.Fatal("expected prod to be protected")
	}
}

func TestProtectCommandAddDuplicate(t *testing.T) {
	m := protect.New(newMemProtectStore())
	root := buildProtectCmd(m)
	root.SetArgs([]string{"protect", "add", "prod"})
	_ = root.Execute()
	root.SetArgs([]string{"protect", "add", "prod"})
	err := root.Execute()
	if !errors.Is(err, protect.ErrProtected) {
		t.Fatalf("expected ErrProtected, got %v", err)
	}
}

func TestProtectCommandList(t *testing.T) {
	m := protect.New(newMemProtectStore())
	_ = m.Protect("prod")
	root := buildProtectCmd(m)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"protect", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "prod") {
		t.Fatalf("expected 'prod' in output, got: %s", buf.String())
	}
}
