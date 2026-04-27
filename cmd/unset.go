package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envport/internal/profile"
	"envport/internal/unset"

	"github.com/spf13/cobra"
)

func init() {
	var strict bool

	cmd := &cobra.Command{
		Use:   "unset <snapshot> <KEY> [KEY...]",
		Short: "Remove keys from a snapshot",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnset(cmd, args, strict)
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "error if a key does not exist in the snapshot")
	rootCmd.AddCommand(cmd)
}

func runUnset(cmd *cobra.Command, args []string, strict bool) error {
	name := args[0]
	keys := args[1:]

	storeDir := filepath.Join(os.Getenv("HOME"), ".envport")
	m, err := profile.New(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	res, err := unset.Run(m, name, keys, strict)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), unset.Format(res))
	return nil
}
