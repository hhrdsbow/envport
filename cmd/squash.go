package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envport/internal/profile"
	"envport/internal/squash"

	"github.com/spf13/cobra"
)

func init() {
	var dest string

	cmd := &cobra.Command{
		Use:   "squash <src1> <src2> [src...] --dest <name>",
		Short: "Merge multiple snapshots into one",
		Long: `Squash loads each source snapshot in order and merges their
variables into a new snapshot. When the same key appears in multiple
sources the last source wins.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSquash(cmd, args, dest)
		},
	}

	cmd.Flags().StringVar(&dest, "dest", "", "name of the output snapshot (required)")
	_ = cmd.MarkFlagRequired("dest")

	rootCmd.AddCommand(cmd)
}

func runSquash(cmd *cobra.Command, sources []string, dest string) error {
	dir := filepath.Join(os.Getenv("HOME"), ".envport")
	m, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	r, err := squash.Run(m, sources, dest)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), squash.Format(r))
	return nil
}
