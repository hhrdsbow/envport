package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/store"
	"github.com/user/envport/internal/tag"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags for snapshots",
	}

	addCmd := &cobra.Command{
		Use:   "add <tag> <snapshot>",
		Short: "Tag a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagAdd(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <tag> <snapshot>",
		Short: "Remove a tag from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagRemove(args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [tag]",
		Short: "List tags or snapshots under a tag",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagList(args)
		},
	}

	tagCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(tagCmd)
}

func openTagManager() (*tag.Manager, error) {
	s, err := store.Open(defaultStorePath())
	if err != nil {
		return nil, err
	}
	return tag.New(s), nil
}

func runTagAdd(t, snapshot string) error {
	mgr, err := openTagManager()
	if err != nil {
		return err
	}
	if err := mgr.Add(t, snapshot); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "tagged %q with %q\n", snapshot, t)
	return nil
}

func runTagRemove(t, snapshot string) error {
	mgr, err := openTagManager()
	if err != nil {
		return err
	}
	return mgr.Remove(t, snapshot)
}

func runTagList(args []string) error {
	mgr, err := openTagManager()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		snaps, err := mgr.Get(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, strings.Join(snaps, "\n"))
		return nil
	}
	tags, err := mgr.List()
	if err != nil {
		return err
	}
	for _, tg := range tags {
		fmt.Fprintln(os.Stdout, tg)
	}
	return nil
}
