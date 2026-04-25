package cmd_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"envport/internal/subset"

	"github.com/spf13/cobra"
)

// memSubsetManager satisfies subset.Manager in-process.
type memSubsetManager struct {
	data map[string]map[string]string
}

func newMemSubsetManager() *memSubsetManager {
	return &memSubsetManager{data: make(map[string]map[string]string)}
}

func (m *memSubsetManager) Load(name string) (map[string]string, error) {
	v, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func (m *memSubsetManager) Save(name string, vars map[string]string) error {
	m.data[name] = vars
	return nil
}

func buildSubsetCmd(m subset.Manager) *cobra.Command {
	var keys string
	c := &cobra.Command{
		Use:  "subset",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			parts := strings.Split(keys, ",")
			r, err := subset.Run(m, args[0], args[1], parts)
			if err != nil {
				return err
			}
			cmd.Println(subset.Format(r))
			return nil
		},
	}
	c.Flags().StringVar(&keys, "keys", "", "keys")
	return c
}

func TestSubsetCommandSuccess(t *testing.T) {
	m := newMemSubsetManager()
	m.data["prod"] = map[string]string{"DB": "pg", "PORT": "5432", "SECRET": "s"}

	buf := &bytes.Buffer{}
	cmd := buildSubsetCmd(m)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"prod", "prod-lite", "--keys", "DB,PORT"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "2 key(s)") {
		t.Errorf("expected summary in output, got: %s", buf.String())
	}
	if _, ok := m.data["prod-lite"]; !ok {
		t.Error("destination snapshot not saved")
	}
}

func TestSubsetCommandMissingSource(t *testing.T) {
	m := newMemSubsetManager()
	cmd := buildSubsetCmd(m)
	cmd.SetArgs([]string{"ghost", "out", "--keys", "A"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing source")
	}
}
