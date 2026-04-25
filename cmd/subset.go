package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/profile"
	"envport/internal/subset"

	"github.com/spf13/cobra"
)

func init() {
	var keys string

	subsetCmd := &cobra.Command{
		Use:   "subset <src> <dst> --keys KEY1,KEY2,...",
		Short: "Extract a subset of keys from a snapshot into a new snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSubset(cmd, args, keys)
		},
	}

	subsetCmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of keys to extract (required)")
	_ = subsetCmd.MarkFlagRequired("keys")

	rootCmd.AddCommand(subsetCmd)
}

func runSubset(cmd *cobra.Command, args []string, keysFlag string) error {
	src, dst := args[0], args[1]

	keys := strings.Split(keysFlag, ",")
	for i, k := range keys {
		keys[i] = strings.TrimSpace(k)
	}

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	r, err := subset.Run(m, src, dst, keys)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), subset.Format(r))
	return nil
}
