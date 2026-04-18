package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/pin"
	"envport/internal/store"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage pinned snapshot aliases",
	}

	addCmd := &cobra.Command{
		Use:   "add <alias> <snapshot>",
		Short: "Pin a snapshot under an alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinAdd(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove a pinned alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pinned aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinList()
		},
	}

	pinCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(pinCmd)
}

func openPinManager() (*pin.Manager, error) {
	s, err := store.Open("pins")
	if err != nil {
		return nil, err
	}
	return pin.New(s), nil
}

func runPinAdd(alias, snapshot string) error {
	m, err := openPinManager()
	if err != nil {
		return err
	}
	if err := m.Pin(alias, snapshot); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Pinned %q -> %q\n", alias, snapshot)
	return nil
}

func runPinRemove(alias string) error {
	m, err := openPinManager()
	if err != nil {
		return err
	}
	return m.Unpin(alias)
}

func runPinList() error {
	m, err := openPinManager()
	if err != nil {
		return err
	}
	list, err := m.List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(os.Stdout, "No pins defined.")
		return nil
	}
	for _, pair := range list {
		fmt.Fprintf(os.Stdout, "%-20s %s\n", pair[0], pair[1])
	}
	return nil
}
