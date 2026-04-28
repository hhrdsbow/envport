package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/profile"
	envselect "envport/internal/select"
)

func init() {
	var patterns []string
	var dst string

	cmd := &cobra.Command{
		Use:   "select <snapshot>",
		Short: "Select keys from a snapshot by name or glob pattern",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSelect(args[0], dst, patterns)
		},
	}

	cmd.Flags().StringVarP(&dst, "out", "o", "", "destination snapshot name (required)")
	_ = cmd.MarkFlagRequired("out")
	cmd.Flags().StringArrayVarP(&patterns, "key", "k", nil, "key name or glob pattern (repeatable)")
	_ = cmd.MarkFlagRequired("key")

	rootCmd.AddCommand(cmd)
}

func runSelect(src, dst string, patterns []string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	r, err := envselect.Run(m, src, dst, patterns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	fmt.Print(envselect.Format(r))
	return nil
}
