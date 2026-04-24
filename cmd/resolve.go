package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envport/internal/alias"
	"github.com/yourorg/envport/internal/profile"
	"github.com/yourorg/envport/internal/resolve"
	"github.com/yourorg/envport/internal/store"
)

func init() {
	resolveCmd := &cobra.Command{
		Use:   "resolve <name>",
		Short: "Resolve an alias or snapshot name to its canonical snapshot name",
		Args:  cobra.ExactArgs(1),
		RunE:  runResolve,
	}
	rootCmd.AddCommand(resolveCmd)
}

func runResolve(cmd *cobra.Command, args []string) error {
	input := args[0]

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	aliasStore, err := store.Open(dir, "aliases")
	if err != nil {
		return fmt.Errorf("open alias store: %w", err)
	}
	aliasManager := alias.New(aliasStore)

	profileStore, err := store.Open(dir, "profiles")
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}
	profileManager := profile.New(profileStore)

	resolver := resolve.New(aliasManager, profileManager)

	canonical, err := resolver.Resolve(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), canonical)
	return nil
}
