package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envport/internal/lowercase"
	"envport/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var keys []string

	cmd := &cobra.Command{
		Use:   "lowercase <snapshot>",
		Short: "Lowercase all (or selected) values in a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLowercase(cmd, args, keys)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Only lowercase these keys (comma-separated)")

	rootCmd.AddCommand(cmd)
}

func runLowercase(cmd *cobra.Command, args []string, keys []string) error {
	name := args[0]

	dir := filepath.Join(os.Getenv("HOME"), ".envport")
	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	r, err := lowercase.Run(mgr, name, keys)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), lowercase.Format(r))
	return nil
}
