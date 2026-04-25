package cmd

import (
	"fmt"
	"strings"

	"github.com/envport/envport/internal/profile"
	"github.com/envport/envport/internal/reorder"
	"github.com/spf13/cobra"
)

var (
	reorderStrategy string
	reorderCustom   string
)

func init() {
	reorderCmd := &cobra.Command{
		Use:   "reorder <snapshot>",
		Short: "Reorder environment variable keys within a snapshot",
		Long: `Reorder sorts the keys of a stored snapshot according to a chosen
strategy and saves the result back to the store.

Strategies:
  alpha    – ascending alphabetical order (default)
  reverse  – descending alphabetical order
  custom   – caller-supplied comma-separated key list; unlisted keys follow
             in alphabetical order`,
		Args: cobra.ExactArgs(1),
		RunE: runReorder,
	}

	reorderCmd.Flags().StringVarP(&reorderStrategy, "strategy", "s", "alpha",
		"reorder strategy: alpha | reverse | custom")
	reorderCmd.Flags().StringVarP(&reorderCustom, "order", "o", "",
		"comma-separated key list for custom strategy")

	rootCmd.AddCommand(reorderCmd)
}

func runReorder(cmd *cobra.Command, args []string) error {
	name := args[0]

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	opts := reorder.Options{
		Strategy: reorder.Strategy(reorderStrategy),
	}

	if reorderCustom != "" {
		for _, k := range strings.Split(reorderCustom, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.CustomOrder = append(opts.CustomOrder, k)
			}
		}
	}

	keys, err := reorder.Run(m, name, opts)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "reordered %d keys in snapshot %q:\n", len(keys), name)
	for i, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s\n", i+1, k)
	}
	return nil
}
