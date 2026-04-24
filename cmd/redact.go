package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envport/internal/profile"
	"envport/internal/redact"
)

func init() {
	var profileDir string
	var customPatterns []string

	cmd := &cobra.Command{
		Use:   "redact <profile> <snapshot>",
		Short: "Print a snapshot with sensitive values masked",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedact(args[0], args[1], profileDir, customPatterns)
		},
	}

	cmd.Flags().StringVar(&profileDir, "dir", defaultProfileDir(), "profile storage directory")
	cmd.Flags().StringSliceVar(&customPatterns, "pattern", nil, "additional key patterns to redact (comma-separated)")

	rootCmd.AddCommand(cmd)
}

func runRedact(profileName, snapshotName, dir string, extra []string) error {
	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile manager: %w", err)
	}

	snap, err := mgr.Load(profileName, snapshotName)
	if err != nil {
		return fmt.Errorf("load snapshot %q/%q: %w", profileName, snapshotName, err)
	}

	patterns := append([]string{}, redact.DefaultPatterns...)
	patterns = append(patterns, extra...)

	masked := redact.Apply(snap.Vars, patterns)

	keys := make([]string, 0, len(masked))
	for k := range masked {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := os.Stdout
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, masked[k])
	}
	return nil
}
