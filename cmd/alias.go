package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envport/internal/alias"
	"envport/internal/store"
)

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage snapshot aliases",
	}

	addCmd := &cobra.Command{
		Use:   "add <alias> <snapshot>",
		Short: "Create an alias pointing to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE:  runAliasAdd,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		RunE:  runAliasRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Args:  cobra.NoArgs,
		RunE:  runAliasList,
	}

	aliasCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(aliasCmd)
}

func openAliasManager() (*alias.Manager, error) {
	s, err := store.Open("aliases")
	if err != nil {
		return nil, err
	}
	return alias.New(s), nil
}

func runAliasAdd(cmd *cobra.Command, args []string) error {
	m, err := openAliasManager()
	if err != nil {
		return err
	}
	if err := m.Add(args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "alias %q → %q saved\n", args[0], args[1])
	return nil
}

func runAliasRemove(cmd *cobra.Command, args []string) error {
	m, err := openAliasManager()
	if err != nil {
		return err
	}
	if err := m.Remove(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "alias %q removed\n", args[0])
	return nil
}

func runAliasList(cmd *cobra.Command, args []string) error {
	m, err := openAliasManager()
	if err != nil {
		return err
	}
	list, err := m.List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no aliases defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ALIAS\tSNAPSHOT")
	for a, t := range list {
		fmt.Fprintf(w, "%s\t%s\n", a, t)
	}
	return w.Flush()
}
