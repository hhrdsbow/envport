package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/group"
	"envport/internal/store"
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage snapshot groups",
	}

	addCmd := &cobra.Command{
		Use:   "add <group> <snapshot>",
		Short: "Add a snapshot to a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupAdd(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <group> <snapshot>",
		Short: "Remove a snapshot from a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <group>",
		Short: "List snapshots in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroupList(args[0])
		},
	}

	groupCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(groupCmd)
}

func openGroupManager() (*group.Manager, error) {
	s, err := store.Open(defaultStorePath())
	if err != nil {
		return nil, err
	}
	return group.New(s), nil
}

func runGroupAdd(grp, snapshot string) error {
	m, err := openGroupManager()
	if err != nil {
		return err
	}
	if err := m.Add(grp, snapshot); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "added %q to group %q\n", snapshot, grp)
	return nil
}

func runGroupRemove(grp, snapshot string) error {
	m, err := openGroupManager()
	if err != nil {
		return err
	}
	if err := m.Remove(grp, snapshot); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "removed %q from group %q\n", snapshot, grp)
	return nil
}

func runGroupList(grp string) error {
	m, err := openGroupManager()
	if err != nil {
		return err
	}
	mems, err := m.Members(grp)
	if err != nil {
		return err
	}
	if len(mems) == 0 {
		fmt.Fprintln(os.Stdout, "(empty group)")
		return nil
	}
	for _, mem := range mems {
		fmt.Fprintln(os.Stdout, mem)
	}
	return nil
}
