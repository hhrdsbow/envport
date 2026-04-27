package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/envport/internal/profile"
	"github.com/yourorg/envport/internal/split"
)

func init() {
	var remainderName string

	cmd := &cobra.Command{
		Use:   "split <snapshot> <dest1=KEY1,KEY2> [dest2=KEY3,...]",
		Short: "Split a snapshot into multiple smaller snapshots",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSplit(cmd, args, remainderName)
		},
	}

	cmd.Flags().StringVar(&remainderName, "remainder", "",
		"name for snapshot holding unassigned keys (omit to discard)")

	rootCmd.AddCommand(cmd)
}

func runSplit(cmd *cobra.Command, args []string, remainder string) error {
	src := args[0]

	targets := make(map[string][]string)
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid target %q: expected name=KEY1,KEY2", arg)
		}
		dest := parts[0]
		keys := strings.Split(parts[1], ",")
		targets[dest] = keys
	}

	storeDir := filepath.Join(os.Getenv("HOME"), ".envport")
	m, err := profile.New(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	results, err := split.Run(m, src, targets, remainder)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), split.Format(results))
	return nil
}
