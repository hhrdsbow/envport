package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envport/internal/profile"
	"envport/internal/shadow"

	"github.com/spf13/cobra"
)

func init() {
	var shadowCmd = &cobra.Command{
		Use:   "shadow <base> <override> [override...]",
		Short: "Show keys in base that are shadowed by override snapshots",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runShadow,
	}
	rootCmd.AddCommand(shadowCmd)
}

func runShadow(cmd *cobra.Command, args []string) error {
	dir := filepath.Join(os.Getenv("HOME"), ".envport")
	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	base := args[0]
	overrides := args[1:]

	results, err := shadow.Run(mgr, base, overrides)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), shadow.Format(results))
	return nil
}
