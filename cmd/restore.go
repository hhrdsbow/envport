package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/profile"
	"envport/internal/restore"

	"github.com/spf13/cobra"
)

func init() {
	var profileName string
	var dryRun bool
	var filterKeys string

	cmd := &cobra.Command{
		Use:   "restore <snapshot>",
		Short: "Restore environment variables from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore(profileName, args[0], dryRun, filterKeys)
		},
	}

	cmd.Flags().StringVarP(&profileName, "profile", "p", "default", "Profile name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print what would be restored without applying")
	cmd.Flags().StringVar(&filterKeys, "keys", "", "Comma-separated list of keys to restore")

	rootCmd.AddCommand(cmd)
}

func runRestore(profileName, snapshotName string, dryRun bool, filterKeys string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile manager: %w", err)
	}

	var keys []string
	if filterKeys != "" {
		keys = strings.Split(filterKeys, ",")
	}

	res, err := restore.Run(mgr, restore.Options{
		ProfileName:  profileName,
		SnapshotName: snapshotName,
		DryRun:       dryRun,
		FilterKeys:   keys,
	})
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintln(os.Stdout, "# Dry run — would export:")
	}
	for k, v := range res.Applied {
		fmt.Fprintf(os.Stdout, "export %s=%q\n", k, v)
	}
	if len(res.Skipped) > 0 {
		fmt.Fprintf(os.Stderr, "skipped keys: %s\n", strings.Join(res.Skipped, ", "))
	}
	return nil
}
