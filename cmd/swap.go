package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/profile"
	"envport/internal/swap"

	"github.com/spf13/cobra"
)

func init() {
	var keys []string

	swapCmd := &cobra.Command{
		Use:   "swap <src> <dst>",
		Short: "Exchange variable values between two snapshots",
		Long: `Swap exchanges the values of shared (or specified) keys between two
snapshots in-place. Both snapshots are updated atomically.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSwap(args[0], args[1], keys)
		},
	}

	swapCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "comma-separated list of keys to swap (default: all common keys)")

	rootCmd.AddCommand(swapCmd)
}

func runSwap(src, dst string, keys []string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	result, err := swap.Run(m, src, dst, keys)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, swap.Format(result))
	return nil
}

// defaultStoreDir returns the directory used for profile storage.
func defaultStoreDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return strings.Join([]string{home, ".envport", "profiles"}, string(os.PathSeparator)), nil
}
