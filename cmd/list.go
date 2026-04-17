package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/profile"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved snapshots",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return fmt.Errorf("resolving store dir: %w", err)
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	names, err := mgr.List()
	if err != nil {
		return fmt.Errorf("listing snapshots: %w", err)
	}

	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "No snapshots saved.")
		return nil
	}

	fmt.Fprintln(os.Stdout, "Saved snapshots:")
	for _, name := range names {
		fmt.Fprintf(os.Stdout, "  - %s\n", name)
	}
	return nil
}
