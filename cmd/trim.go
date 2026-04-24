package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/profile"
	"envport/internal/trim"

	"github.com/spf13/cobra"
)

func init() {
	var dryRun bool

	trimCmd := &cobra.Command{
		Use:   "trim <profile> <KEY> [KEY...]",
		Short: "Remove specific keys from a saved snapshot",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTrim(args[0], args[1:], dryRun)
		},
	}

	trimCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview keys that would be removed without modifying the snapshot")
	RootCmd.AddCommand(trimCmd)
}

func runTrim(profileName string, keys []string, dryRun bool) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	res, err := trim.Run(m, profileName, keys, dryRun)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	if dryRun {
		fmt.Println("[dry-run] keys that would be removed:")
	} else {
		fmt.Println("Keys removed:")
	}

	if len(res.Removed) == 0 {
		fmt.Println("  (none matched)")
	} else {
		for _, k := range res.Removed {
			fmt.Println(" -", k)
		}
	}

	if !dryRun {
		fmt.Printf("Snapshot %q now contains %d key(s): %s\n",
			profileName, res.Kept, strings.Join(keys, ", "))
	}
	return nil
}
