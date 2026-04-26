package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"envport/internal/chain"
	"envport/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var skipMissing bool

	cmd := &cobra.Command{
		Use:   "chain <snapshot> [snapshot...]",
		Short: "Merge multiple snapshots in sequence and export the result",
		Long: `Apply a list of snapshots in order, with later snapshots overriding
earlier ones for duplicate keys. Prints shell export statements to stdout.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChain(cmd, args, skipMissing)
		},
	}

	cmd.Flags().BoolVar(&skipMissing, "skip-missing", false, "skip snapshots that do not exist instead of erroring")
	rootCmd.AddCommand(cmd)
}

func runChain(cmd *cobra.Command, names []string, skipMissing bool) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.New(filepath.Join(dir, "profiles"))
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	r, err := chain.Run(m, names, skipMissing)
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	for _, k := range sortedStringKeys(r.Vars) {
		fmt.Fprintf(w, "export %s=%q\n", k, r.Vars[k])
	}

	fmt.Fprintln(os.Stderr, chain.Format(r))
	return nil
}
