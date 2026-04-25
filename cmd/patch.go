package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/patch"
	"envport/internal/profile"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func init() {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "patch <profile> KEY=VALUE [KEY=VALUE ...]",
		Short: "Partially update a snapshot with specific key-value pairs",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPatch(cmd, args, dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would change without saving")
	rootCmd.AddCommand(cmd)
}

func runPatch(cmd *cobra.Command, args []string, dryRun bool) error {
	profileName := args[0]
	updates := make(map[string]string)

	for _, pair := range args[1:] {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid key=value pair: %q", pair)
		}
		updates[parts[0]] = parts[1]
	}

	db, err := store.Open(defaultStorePath())
	if err != nil {
		return err
	}
	defer db.Close()

	m := profile.New(db)

	if dryRun {
		vars, err := m.Load(profileName)
		if err != nil {
			return err
		}
		for k, v := range updates {
			if _, exists := vars[k]; exists {
				fmt.Fprintf(cmd.OutOrStdout(), "~ %s=%s\n", k, v)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", k, v)
			}
		}
		return nil
	}

	result, err := patch.Run(m, profileName, updates)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintln(cmd.OutOrStdout(), patch.Format(result))
	return nil
}
