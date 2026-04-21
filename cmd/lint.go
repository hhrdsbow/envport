package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/lint"
	"envport/internal/profile"
)

func init() {
	var dir string

	cmd := &cobra.Command{
		Use:   "lint <snapshot>",
		Short: "Check a snapshot for common issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(args[0], dir)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", defaultProfileDir(), "profile storage directory")
	rootCmd.AddCommand(cmd)
}

func runLint(name, dir string) error {
	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("lint: open profile store: %w", err)
	}

	issues, err := lint.Run(mgr, name)
	if err != nil {
		return err
	}

	out := lint.Format(issues)
	fmt.Fprintln(os.Stdout, out)

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}
