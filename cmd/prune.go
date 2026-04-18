package cmd

import (
	"fmt"
	"os"
	"time"

	"envport/internal/profile"
	"envport/internal/prune"

	"github.com/spf13/cobra"
)

func init() {
	var (
		keepLast int
		olderThan string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove old snapshots by age or count",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrune(keepLast, olderThan, dryRun)
		},
	}

	cmd.Flags().IntVar(&keepLast, "keep-last", 0, "Keep only the N most recent snapshots")
	cmd.Flags().StringVar(&olderThan, "older-than", "", "Remove snapshots older than duration (e.g. 72h)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print what would be deleted without deleting")

	rootCmd.AddCommand(cmd)
}

func runPrune(keepLast int, olderThan string, dryRun bool) error {
	pm, err := profile.Open(defaultStorePath())
	if err != nil {
		return err
	}

	opts := prune.Options{DryRun: dryRun, KeepLast: keepLast}

	if olderThan != "" {
		d, err := time.ParseDuration(olderThan)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", olderThan, err)
		}
		opts.OlderThan = time.Now().Add(-d)
	}

	if opts.KeepLast == 0 && opts.OlderThan.IsZero() {
		return fmt.Errorf("specify --keep-last or --older-than")
	}

	res, err := prune.Run(pm, opts)
	if err != nil {
		return err
	}

	if len(res.Pruned) == 0 {
		fmt.Fprintln(os.Stdout, "Nothing to prune.")
		return nil
	}

	verb := "Deleted"
	if dryRun {
		verb = "Would delete"
	}
	for _, name := range res.Pruned {
		fmt.Fprintf(os.Stdout, "%s: %s\n", verb, name)
	}
	return nil
}
