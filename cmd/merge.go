package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/merge"
	"github.com/user/envport/internal/profile"
)

func init() {
	var strategy string

	cmd := &cobra.Command{
		Use:   "merge <base> <other> <dest>",
		Short: "Merge two snapshots into a new snapshot",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(args[0], args[1], args[2], strategy)
		},
	}
	cmd.Flags().StringVarP(&strategy, "strategy", "s", "base", "conflict strategy: base, other, error")
	rootCmd.AddCommand(cmd)
}

func runMerge(base, other, dest, strategyStr string) error {
	dir, err := snapshotDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return err
	}

	var strategy merge.Strategy
	switch strings.ToLower(strategyStr) {
	case "base":
		strategy = merge.StrategyBase
	case "other":
		strategy = merge.StrategyOther
	case "error":
		strategy = merge.StrategyError
	default:
		return fmt.Errorf("unknown strategy %q (base|other|error)", strategyStr)
	}

	res, err := merge.Run(mgr, base, other, dest, strategy)
	if err != nil {
		return err
	}

	if len(res.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "conflicts resolved (%s): %s\n", strategyStr, strings.Join(res.Conflicts, ", "))
	}
	fmt.Printf("merged %d vars into snapshot %q\n", len(res.Vars), dest)
	return nil
}
