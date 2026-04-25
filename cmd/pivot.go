package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/pivot"
	"github.com/user/envport/internal/profile"
)

func init() {
	var destName string

	cmd := &cobra.Command{
		Use:   "pivot <snapshot>",
		Short: "Transpose a snapshot so values become keys and keys become values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPivot(cmd, args, destName)
		},
	}

	cmd.Flags().StringVarP(&destName, "dest", "d", "", "name for the pivoted snapshot (default: <snapshot>-pivoted)")
	rootCmd.AddCommand(cmd)
}

func runPivot(cmd *cobra.Command, args []string, destName string) error {
	src := args[0]
	if destName == "" {
		destName = src + "-pivoted"
	}

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	res, err := pivot.Run(mgr, src, destName)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "pivoted %q → %q (%d keys)\n", res.Src, res.Dest, res.Keys)
	return nil
}
