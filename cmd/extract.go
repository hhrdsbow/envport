package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/extract"
	"envport/internal/profile"

	"github.com/spf13/cobra"
)

func init() {
	var keys []string

	cmd := &cobra.Command{
		Use:   "extract <src> <dst> --keys KEY1,KEY2,...",
		Short: "Extract specific keys from a snapshot into a new snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExtract(cmd, args, keys)
		},
	}

	cmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated list of keys to extract (required)")
	_ = cmd.MarkFlagRequired("keys")

	rootCmd.AddCommand(cmd)
}

func runExtract(cmd *cobra.Command, args []string, keys []string) error {
	src, dst := args[0], args[1]

	dir, err := snapshotDir()
	if err != nil {
		return err
	}

	mgr, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	opts := extract.Options{
		Src:  src,
		Dst:  dst,
		Keys: keys,
	}

	out, err := extract.Run(mgr, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Extracted %d key(s) from %q into %q:\n", len(out), src, dst)
	for _, k := range sortedStringKeys(out) {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, out[k])
	}
	return nil
}

// sortedStringKeys returns the keys of m in sorted order.
func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStringsLocal(keys)
	return keys
}

func sortStringsLocal(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && strings.Compare(s[j-1], s[j]) > 0; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
