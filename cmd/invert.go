package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/invert"
	"github.com/user/envport/internal/profile"
)

func init() {
	var storeDir string

	cmd := &cobra.Command{
		Use:   "invert <src> <dst>",
		Short: "Swap keys and values of a snapshot into a new snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInvert(storeDir, args[0], args[1])
		},
	}

	cmd.Flags().StringVar(&storeDir, "store", defaultStoreDir(), "path to snapshot store")
	rootCmd.AddCommand(cmd)
}

func runInvert(storeDir, src, dst string) error {
	m, err := profile.New(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	r, err := invert.Run(m, src, dst)
	if err != nil {
		return err
	}

	fmt.Println(invert.Format(r))
	return nil
}

// defaultStoreDir returns the default directory used to persist snapshots.
// It is defined here to avoid duplication across cmd files that already
// declare it; if another cmd file owns the canonical definition this file
// should be updated to reference that symbol instead.
func defaultStoreDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".envport"
	}
	return filepath.Join(home, ".envport")
}
