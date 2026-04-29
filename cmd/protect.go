package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"envport/internal/protect"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func openProtectManager() (*protect.Manager, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".envport", "protect")
	s, err := store.Open(dir)
	if err != nil {
		return nil, err
	}
	return protect.New(s), nil
}

func init() {
	protectCmd := &cobra.Command{
		Use:   "protect",
		Short: "Manage snapshot write-protection",
	}

	addCmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Protect a snapshot from modification or deletion",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProtectAdd(args[0])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove protection from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProtectRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all protected snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProtectList()
		},
	}

	protectCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(protectCmd)
}

func runProtectAdd(name string) error {
	m, err := openProtectManager()
	if err != nil {
		return err
	}
	if err := m.Protect(name); err != nil {
		return err
	}
	fmt.Printf("snapshot %q is now protected\n", name)
	return nil
}

func runProtectRemove(name string) error {
	m, err := openProtectManager()
	if err != nil {
		return err
	}
	if err := m.Unprotect(name); err != nil {
		return err
	}
	fmt.Printf("protection removed from %q\n", name)
	return nil
}

func runProtectList() error {
	m, err := openProtectManager()
	if err != nil {
		return err
	}
	names, err := m.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Println("no protected snapshots")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME")
	for _, n := range names {
		fmt.Fprintln(w, n)
	}
	return w.Flush()
}
