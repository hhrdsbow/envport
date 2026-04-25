package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/normalize"
	"envport/internal/profile"
)

func init() {
	var (
		noUpperKeys   bool
		noTrimValues  bool
		removeEmpty   bool
		dryRun        bool
	)

	cmd := &cobra.Command{
		Use:   "normalize <profile>",
		Short: "Normalize variable keys and values in a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNormalize(args[0], normalize.Options{
				UpperKeys:   !noUpperKeys,
				TrimValues:  !noTrimValues,
				RemoveEmpty: removeEmpty,
			}, dryRun, cmd)
		},
	}

	cmd.Flags().BoolVar(&noUpperKeys, "no-upper", false, "skip uppercasing keys")
	cmd.Flags().BoolVar(&noTrimValues, "no-trim", false, "skip trimming values")
	cmd.Flags().BoolVar(&removeEmpty, "remove-empty", false, "remove keys with empty values")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show changes without saving")

	rootCmd.AddCommand(cmd)
}

func runNormalize(profileName string, opts normalize.Options, dryRun bool, cmd *cobra.Command) error {
	m, err := profile.NewManager(defaultStorePath())
	if err != nil {
		return err
	}

	snap, err := m.Load(profileName)
	if err != nil {
		return fmt.Errorf("load profile %q: %w", profileName, err)
	}

	result := normalize.Run(snap.Vars, opts)

	if !result.Changed {
		fmt.Fprintln(cmd.OutOrStdout(), "no changes")
		return nil
	}

	summary := normalize.Format(snap.Vars, result.Vars)
	fmt.Fprint(cmd.OutOrStdout(), summary)

	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "(dry-run: changes not saved)")
		return nil
	}

	snap.Vars = result.Vars
	if err := m.Save(profileName, snap); err != nil {
		return fmt.Errorf("save profile %q: %w", profileName, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "profile %q normalized\n", profileName)
	os.Exit(0)
	return nil
}
