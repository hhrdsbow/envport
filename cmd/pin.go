package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envport/internal/pin"
	"envport/internal/store"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage pinned snapshots",
	}

	pinCmd.AddCommand(
		&cobra.Command{
			Use:   "add <name>",
			Short: "Pin a snapshot to protect it from pruning",
			Args:  cobra.ExactArgs(1),
			RunE:  runPinAdd,
		},
		&cobra.Command{
			Use:   "remove <name>",
			Short: "Unpin a snapshot",
			Args:  cobra.ExactArgs(1),
			RunE:  runPinRemove,
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all pinned snapshots",
			Args:  cobra.NoArgs,
			RunE:  runPinList,
		},
	)

	rootCmd.AddCommand(pinCmd)
}

func openPinManager() (*pin.Manager, error) {
	s, err := store.Open(defaultStorePath())
	if err != nil {
		return nil, err
	}
	return pin.New(s), nil
}

func runPinAdd(_ *cobra.Command, args []string) error {
	m, err := openPinManager()
	if err != nil {
		return err
	}
	if err := m.Add(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "pinned %q\n", args[0])
	return nil
}

func runPinRemove(_ *cobra.Command, args []string) error {
	m, err := openPinManager()
	if err != nil {
		return err
	}
	if err := m.Remove(args[0]); err != nil {
		if errors.Is(err, pin.ErrNotPinned) {
			return fmt.Errorf("%q is not pinned", args[0])
		}
		return err
	}
	fmt.Fprintf(os.Stdout, "unpinned %q\n", args[0])
	return nil
}

func runPinList(_ *cobra.Command, _ []string) error {
	m, err := openPinManager()
	if err != nil {
		return err
	}
	names, err := m.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no pinned snapshots")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME")
	for _, n := range names {
		fmt.Fprintln(w, n)
	}
	return w.Flush()
}
