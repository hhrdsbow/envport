package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"envport/internal/lock"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func init() {
	lockCmd := &cobra.Command{Use: "lock", Short: "Manage snapshot locks"}

	addCmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Lock a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reason, _ := cmd.Flags().GetString("reason")
			return runLockAdd(args[0], reason)
		},
	}
	addCmd.Flags().String("reason", "", "reason for locking")

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Unlock a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List locked snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockList()
		},
	}

	lockCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(lockCmd)
}

func openLockManager() (*lock.Manager, error) {
	s, err := store.Open(defaultStorePath())
	if err != nil {
		return nil, err
	}
	return lock.New(s), nil
}

func runLockAdd(name, reason string) error {
	mgr, err := openLockManager()
	if err != nil {
		return err
	}
	if err := mgr.Lock(name, reason); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "locked: %s\n", name)
	return nil
}

func runLockRemove(name string) error {
	mgr, err := openLockManager()
	if err != nil {
		return err
	}
	if err := mgr.Unlock(name); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "unlocked: %s\n", name)
	return nil
}

func runLockList() error {
	mgr, err := openLockManager()
	if err != nil {
		return err
	}
	entries, err := mgr.List()
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tLOCKED AT\tREASON")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Name, e.LockedAt.Format("2006-01-02 15:04:05"), e.Reason)
	}
	return w.Flush()
}
