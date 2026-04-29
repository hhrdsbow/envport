package cmd

import (
	"fmt"
	"os"

	"envport/internal/profile"
	"envport/internal/sortkeys"

	"github.com/spf13/cobra"
)

func init() {
	var strategy string

	cmd := &cobra.Command{
		Use:   "sortkeys <profile>",
		Short: "Display snapshot keys in a specified sort order",
		Long: `Sort and display the keys of a saved snapshot.

Strategies:
  alpha     - alphabetical (default)
  reverse   - reverse alphabetical
  keylen    - shortest key first
  valuelen  - shortest value first`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSortKeys(args[0], sortkeys.Strategy(strategy))
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "alpha",
		"sort strategy: alpha, reverse, keylen, valuelen")

	rootCmd.AddCommand(cmd)
}

func runSortKeys(profileName string, strategy sortkeys.Strategy) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	result, err := sortkeys.Run(mgr, profileName, strategy)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, sortkeys.Format(result))
	return nil
}
