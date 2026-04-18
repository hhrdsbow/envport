package cmd

import (
	"fmt"
	"os"

	"envport/internal/clone"
	"envport/internal/profile"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

var cloneSuffix string

func init() {
	cloneCmd := &cobra.Command{
		Use:   "clone <src> [dest]",
		Short: "Clone an existing snapshot into a new name",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runClone,
	}
	cloneCmd.Flags().StringVar(&cloneSuffix, "suffix", "-copy", "suffix to append when dest already exists")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	src := args[0]
	dest := ""
	if len(args) == 2 {
		dest = args[1]
	}

	st, err := store.Open(storeDir())
	if err != nil {
		return err
	}
	m := profile.New(st)

	opts := clone.Options{Suffix: cloneSuffix}
	final, err := clone.Run(m, src, dest, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	fmt.Printf("cloned %q → %q\n", src, final)
	return nil
}
