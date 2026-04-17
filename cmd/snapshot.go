package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/env"
	"envport/internal/profile"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot [profile]",
	Short: "Capture current environment variables into a named profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runSnapshot,
}

var snapshotKeys []string

func init() {
	snapshotCmd.Flags().StringSliceVarP(&snapshotKeys, "keys", "k", nil, "Only capture specific keys (comma-separated)")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	name := args[0]

	captured := env.Capture()
	if len(snapshotKeys) > 0 {
		captured = env.FilterKeys(captured, snapshotKeys)
	}

	dataDir, err := defaultDataDir()
	if err != nil {
		return fmt.Errorf("resolving data dir: %w", err)
	}

	mgr, err := profile.New(dataDir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	if err := mgr.Save(name, captured); err != nil {
		return fmt.Errorf("saving profile %q: %w", name, err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot saved to profile %q (%d variables)\n", name, len(captured))
	return nil
}
