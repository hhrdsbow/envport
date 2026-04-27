package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"envport/internal/scope"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func openScopeManager() (*scope.Manager, error) {
	dir, err := defaultStoreDir()
	if err != nil {
		return nil, err
	}
	s, err := store.Open(filepath.Join(dir, "scopes.db"))
	if err != nil {
		return nil, err
	}
	return scope.New(s), nil
}

func init() {
	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage named key scopes",
	}

	addCmd := &cobra.Command{
		Use:   "add <name> <KEY>...",
		Short: "Create or replace a scope with the given keys",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScopeAdd(args[0], args[1:])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScopeRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all scopes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScopeList()
		},
	}

	scopeCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(scopeCmd)
}

func runScopeAdd(name string, keys []string) error {
	mgr, err := openScopeManager()
	if err != nil {
		return err
	}
	if err := mgr.Add(name, keys); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "scope %q saved with keys: %s\n", name, strings.Join(keys, ", "))
	return nil
}

func runScopeRemove(name string) error {
	mgr, err := openScopeManager()
	if err != nil {
		return err
	}
	if err := mgr.Delete(name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "scope %q removed\n", name)
	return nil
}

func runScopeList() error {
	mgr, err := openScopeManager()
	if err != nil {
		return err
	}
	names, err := mgr.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no scopes defined")
		return nil
	}
	for _, n := range names {
		sc, _ := mgr.Get(n)
		fmt.Fprintf(os.Stdout, "%-20s %s\n", n, strings.Join(sc.Keys, ", "))
	}
	return nil
}
