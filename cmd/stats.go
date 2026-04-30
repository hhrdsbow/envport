package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/profile"
	"github.com/user/envport/internal/stats"
)

func init() {
	var pattern string
	var topN int

	cmd := &cobra.Command{
		Use:   "stats <name>",
		Short: "Show statistics for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(args[0], pattern, topN)
		},
	}

	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "filter keys by prefix")
	cmd.Flags().IntVarP(&topN, "top", "n", 5, "number of longest-value keys to display (0 = all)")

	rootCmd.AddCommand(cmd)
}

func runStats(name, pattern string, topN int) error {
	dir := filepath.Join(os.Getenv("HOME"), ".envport")
	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	res, err := stats.Run(mgr, name, pattern, topN)
	if err != nil {
		return err
	}

	fmt.Print(stats.Format(res))
	return nil
}
