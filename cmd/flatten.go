package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/flatten"
	"github.com/user/envport/internal/profile"
)

func init() {
	var dest string

	cmd := &cobra.Command{
		Use:   "flatten <name> [name...] --dest <dest>",
		Short: "Merge multiple snapshots into one",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if dest == "" {
				return fmt.Errorf("--dest is required")
			}
			dir := filepath.Join(os.Getenv("HOME"), ".envport")
			m, err := profile.New(dir)
			if err != nil {
				return err
			}
			return runFlatten(cmd, m, args, dest)
		},
	}

	cmd.Flags().StringVar(&dest, "dest", "", "name of the destination snapshot")
	rootCmd.AddCommand(cmd)
}

func runFlatten(cmd *cobra.Command, m flatten.Manager, names []string, dest string) error {
	res, err := flatten.Run(m, names, dest)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "merged %d sources into %q (%d keys",
		len(res.Sources), dest, len(res.Vars))
	if res.Conflicts > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), ", %d conflict(s) resolved", res.Conflicts)
	}
	fmt.Fprintln(cmd.OutOrStdout(), ")")
	return nil
}
