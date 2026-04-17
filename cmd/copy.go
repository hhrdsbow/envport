package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	icoppy "envport/internal/copy"
	"envport/internal/profile"
)

func init() {
	var overwrite bool

	copyCmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy a snapshot to a new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCopy(args[0], args[1], overwrite)
		},
	}

	copyCmd.Flags().BoolVarP(&overwrite, "overwrite", "f", false, "overwrite destination if it exists")
	rootCmd.AddCommand(copyCmd)
}

func runCopy(src, dst string, overwrite bool) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	if err := icoppy.Run(m, src, dst, overwrite); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	fmt.Printf("copied snapshot %q → %q\n", src, dst)
	return nil
}
