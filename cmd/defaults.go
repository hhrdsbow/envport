package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"envport/internal/defaults"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func init() {
	defaultsCmd := &cobra.Command{
		Use:   "defaults",
		Short: "Manage default key-value pairs for a profile",
	}

	setCmd := &cobra.Command{
		Use:   "set <profile> KEY=VALUE...",
		Short: "Set default values for a profile",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runDefaultsSet,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <profile> KEY...",
		Short: "Remove default keys from a profile",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runDefaultsRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List profiles that have defaults registered",
		Args:  cobra.NoArgs,
		RunE:  runDefaultsList,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <profile>",
		Short: "Clear all defaults for a profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runDefaultsClear,
	}

	defaultsCmd.AddCommand(setCmd, removeCmd, listCmd, clearCmd)
	rootCmd.AddCommand(defaultsCmd)
}

func openDefaultsManager() (*defaults.Manager, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".envport", "defaults")
	s, err := store.Open(dir)
	if err != nil {
		return nil, err
	}
	return defaults.New(s), nil
}

func runDefaultsSet(cmd *cobra.Command, args []string) error {
	profile := args[0]
	pairs := make(map[string]string)
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid pair %q: expected KEY=VALUE", arg)
		}
		pairs[parts[0]] = parts[1]
	}
	m, err := openDefaultsManager()
	if err != nil {
		return err
	}
	if err := m.Set(profile, pairs); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "defaults updated for profile %q\n", profile)
	return nil
}

func runDefaultsRemove(cmd *cobra.Command, args []string) error {
	profile := args[0]
	m, err := openDefaultsManager()
	if err != nil {
		return err
	}
	if err := m.Remove(profile, args[1:]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "removed keys from defaults for profile %q\n", profile)
	return nil
}

func runDefaultsList(cmd *cobra.Command, args []string) error {
	m, err := openDefaultsManager()
	if err != nil {
		return err
	}
	names, err := m.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no profiles with defaults")
		return nil
	}
	for _, n := range names {
		fmt.Fprintln(cmd.OutOrStdout(), n)
	}
	return nil
}

func runDefaultsClear(cmd *cobra.Command, args []string) error {
	m, err := openDefaultsManager()
	if err != nil {
		return err
	}
	if err := m.Clear(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "cleared defaults for profile %q\n", args[0])
	return nil
}
