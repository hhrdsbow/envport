package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/history"
	"envport/internal/profile"
	"envport/internal/rollback"
)

func init() {
	var offset int

	cmd := &cobra.Command{
		Use:   "rollback <profile>",
		Short: "Restore a profile to a previous snapshot from its history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRollback(args[0], offset)
		},
	}
	cmd.Flags().IntVarP(&offset, "offset", "n", 1, "How many steps back to roll (1 = most recent history entry)")
	rootCmd.AddCommand(cmd)
}

func runRollback(profileName string, offset int) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	hr, err := history.New(dir)
	if err != nil {
		return fmt.Errorf("opening history store: %w", err)
	}

	res, err := rollback.Run(mgr, hr, profileName, offset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rollback failed: %v\n", err)
		return err
	}

	fmt.Printf("Rolled back %q to snapshot %q\n", res.Profile, res.RolledTo)
	return nil
}
