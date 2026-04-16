package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/diff"
	"envport/internal/env"
	"envport/internal/profile"
)

var diffCmd = &cobra.Command{
	Use:   "diff <profile>",
	Short: "Show differences between current environment and a saved profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		mgr, err := profile.New(defaultStorePath())
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		saved, err := mgr.Load(name)
		if err != nil {
			return fmt.Errorf("load profile %q: %w", name, err)
		}

		current := env.Capture()
		changes := diff.Compute(saved, current)

		if len(changes) == 0 {
			fmt.Println("No differences.")
			return nil
		}

		for _, c := range changes {
			fmt.Fprintln(os.Stdout, diff.Format(c))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
