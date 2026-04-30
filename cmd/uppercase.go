package cmd

import (
	"fmt"
	"os"

	"github.com/user/envport/internal/profile"
	"github.com/user/envport/internal/uppercase"
	"github.com/spf13/cobra"
)

func init() {
	var keys []string

	cmd := &cobra.Command{
		Use:   "uppercase <profile>",
		Short: "Uppercase all or selected values in a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUppercase(cmd, args, keys)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "limit to specific keys")
	rootCmd.AddCommand(cmd)
}

func runUppercase(cmd *cobra.Command, args []string, keys []string) error {
	name := args[0]

	storeDir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.NewManager(storeDir)
	if err != nil {
		return err
	}

	result, err := uppercase.Run(mgr, name, keys)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, uppercase.Format(result))
	return nil
}
