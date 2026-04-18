package cmd

import (
	"fmt"
	"os"

	"github.com/envport/envport/internal/compare"
	"github.com/envport/envport/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var compareCmd = &cobra.Command{
		Use:   "compare <base> <other>",
		Short: "Compare two snapshots and show differences",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompare,
	}
	rootCmd.AddCommand(compareCmd)
}

type profileSnap struct {
	name string
	vars map[string]string
}

func (p *profileSnap) Name() string            { return p.name }
func (p *profileSnap) Vars() map[string]string { return p.vars }

func runCompare(cmd *cobra.Command, args []string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}
	mgr, err := profile.New(dir)
	if err != nil {
		return err
	}

	baseVars, err := mgr.Load(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", args[0], err)
		return err
	}
	otherVars, err := mgr.Load(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", args[1], err)
		return err
	}

	base := &profileSnap{name: args[0], vars: baseVars}
	other := &profileSnap{name: args[1], vars: otherVars}

	result := compare.Run(base, other)
	fmt.Print(compare.Format(result))
	return nil
}
