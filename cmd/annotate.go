package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/annotate"
	"github.com/user/envport/internal/store"
)

func init() {
	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Manage notes attached to snapshots",
	}

	setCmd := &cobra.Command{
		Use:   "set <snapshot> <note>",
		Short: "Attach a note to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateSet(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <snapshot>",
		Short: "Remove the note from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all annotated snapshots",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotateList()
		},
	}

	annotateCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(annotateCmd)
}

func openAnnotateManager() (*annotate.Manager, error) {
	s, err := store.Open(defaultStorePath())
	if err != nil {
		return nil, err
	}
	return annotate.New(s), nil
}

func runAnnotateSet(name, note string) error {
	mgr, err := openAnnotateManager()
	if err != nil {
		return err
	}
	if err := mgr.Set(name, note); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "annotation set for %q\n", name)
	return nil
}

func runAnnotateRemove(name string) error {
	mgr, err := openAnnotateManager()
	if err != nil {
		return err
	}
	return mgr.Remove(name)
}

func runAnnotateList() error {
	mgr, err := openAnnotateManager()
	if err != nil {
		return err
	}
	list, err := mgr.List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(os.Stdout, "no annotations found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SNAPSHOT\tNOTE")
	for name, note := range list {
		fmt.Fprintf(w, "%s\t%s\n", name, note)
	}
	return w.Flush()
}
