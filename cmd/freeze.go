package cmd

import (
	"fmt"
	"os"
	"sort"

	"envport/internal/freeze"
	"envport/internal/store"
	"github.com/spf13/cobra"
)

func openFreezeManager() (*freeze.Manager, error) {
	s, err := store.Open(defaultStoreDir(), "freeze")
	if err != nil {
		return nil, err
	}
	return freeze.New(s), nil
}

func init() {
	freezeCmd := &cobra.Command{
		Use:   "freeze",
		Short: "Manage frozen (immutable) snapshots",
	}

	addCmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Freeze a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFreezeAdd(args[0])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Unfreeze a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFreezeRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all frozen snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFreezeList()
		},
	}

	freezeCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(freezeCmd)
}

func runFreezeAdd(name string) error {
	mgr, err := openFreezeManager()
	if err != nil {
		return err
	}
	if err := mgr.Freeze(name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "snapshot %q frozen\n", name)
	return nil
}

func runFreezeRemove(name string) error {
	mgr, err := openFreezeManager()
	if err != nil {
		return err
	}
	if err := mgr.Unfreeze(name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "snapshot %q unfrozen\n", name)
	return nil
}

func runFreezeList() error {
	mgr, err := openFreezeManager()
	if err != nil {
		return err
	}
	names, err := mgr.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no frozen snapshots")
		return nil
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(os.Stdout, n)
	}
	return nil
}
