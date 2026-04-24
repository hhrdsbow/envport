package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/profile"
	"envport/internal/sanitize"
)

func init() {
	var (
		trimSpace   bool
		removeEmpty bool
		upperKeys   bool
		stripPrefix string
		dryRun      bool
	)

	cmd := &cobra.Command{
		Use:   "sanitize <profile>",
		Short: "Clean and normalise a stored snapshot's environment variables",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSanitize(args[0], sanitize.Options{
				TrimSpace:   trimSpace,
				RemoveEmpty: removeEmpty,
				UpperKeys:   upperKeys,
				StripPrefix: stripPrefix,
			}, dryRun, cmd)
		},
	}

	cmd.Flags().BoolVar(&trimSpace, "trim", true, "trim whitespace from values")
	cmd.Flags().BoolVar(&removeEmpty, "remove-empty", false, "drop keys with empty values")
	cmd.Flags().BoolVar(&upperKeys, "upper", false, "normalise keys to upper-case")
	cmd.Flags().StringVar(&stripPrefix, "strip-prefix", "", "strip a prefix from every key")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print changes without saving")

	rootCmd.AddCommand(cmd)
}

func runSanitize(profileName string, opts sanitize.Options, dryRun bool, cmd *cobra.Command) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	snap, err := mgr.Load(profileName)
	if err != nil {
		return fmt.Errorf("load profile %q: %w", profileName, err)
	}

	cleaned, res := sanitize.Run(snap.Vars, opts)

	if !res.Changed() {
		fmt.Fprintln(cmd.OutOrStdout(), "no changes")
		return nil
	}

	if len(res.TrimmedKeys) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "trimmed: %v\n", res.TrimmedKeys)
	}
	if len(res.RenamedKeys) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "renamed: %v\n", res.RenamedKeys)
	}
	if len(res.DroppedKeys) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "dropped: %v\n", res.DroppedKeys)
	}

	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "(dry-run: no changes saved)")
		return nil
	}

	snap.Vars = cleaned
	if err := mgr.Save(profileName, snap); err != nil {
		fmt.Fprintf(os.Stderr, "error saving profile: %v\n", err)
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), "profile saved")
	return nil
}
