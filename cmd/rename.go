package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/profile"
	"github.com/user/envport/internal/rename"
)

func init() {
	renameCmd := &cobra.Command{
		Use:   "rename <src> <dst>",
		Short: "Rename a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE:  runRename,
	}
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	src, dst := args[0], args[1]

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	if err := rename.Run(mgr, src, dst); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Renamed snapshot %q to %q\n", src, dst)
	return nil
}
