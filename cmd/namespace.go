package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/namespace"
	"envport/internal/store"
)

func init() {
	nsCmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage snapshot namespaces",
	}

	addCmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Register a new namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNamespaceAdd(args[0])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete an existing namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNamespaceRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all namespaces",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNamespaceList()
		},
	}

	nsCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(nsCmd)
}

func openNamespaceManager() (*namespace.Manager, error) {
	s, err := store.Open("")
	if err != nil {
		return nil, err
	}
	return namespace.New(s), nil
}

func runNamespaceAdd(name string) error {
	mgr, err := openNamespaceManager()
	if err != nil {
		return err
	}
	if err := mgr.Add(name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "namespace %q added\n", name)
	return nil
}

func runNamespaceRemove(name string) error {
	mgr, err := openNamespaceManager()
	if err != nil {
		return err
	}
	if err := mgr.Remove(name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "namespace %q removed\n", name)
	return nil
}

func runNamespaceList() error {
	mgr, err := openNamespaceManager()
	if err != nil {
		return err
	}
	names, err := mgr.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no namespaces defined")
		return nil
	}
	for _, n := range names {
		fmt.Fprintln(os.Stdout, n)
	}
	return nil
}
