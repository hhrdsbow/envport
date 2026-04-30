package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envport/internal/profile"
	"envport/internal/slice"

	"github.com/spf13/cobra"
)

func init() {
	var from, to int
	var dst string

	cmd := &cobra.Command{
		Use:   "slice <snapshot> --from N --to N [--dst <name>]",
		Short: "Extract a range of keys from a snapshot by sorted position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSlice(args[0], dst, from, to)
		},
	}

	cmd.Flags().IntVar(&from, "from", 0, "start index (inclusive, 0-based)")
	cmd.Flags().IntVar(&to, "to", -1, "end index (exclusive); -1 means end of list")
	cmd.Flags().StringVar(&dst, "dst", "", "destination snapshot name (optional)")

	rootCmd.AddCommand(cmd)
}

func runSlice(src, dst string, from, to int) error {
	dir := filepath.Join(os.Getenv("HOME"), ".envport")
	m, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("slice: open store: %w", err)
	}

	r, err := slice.Run(m, src, dst, from, to)
	if err != nil {
		return err
	}

	fmt.Println(slice.Format(r))
	for _, k := range r.Keys {
		fmt.Printf("  %s=%s\n", k, r.Vars[k])
	}
	return nil
}
