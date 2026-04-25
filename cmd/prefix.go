package cmd

import (
	"fmt"
	"os"

	"envport/internal/prefix"
	"envport/internal/profile"

	"github.com/spf13/cobra"
)

func init() {
	prefixCmd := &cobra.Command{
		Use:   "prefix",
		Short: "Add or remove a key prefix across a snapshot",
	}

	addCmd := &cobra.Command{
		Use:   "add <snapshot> <prefix>",
		Short: "Add prefix to all keys in a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dst, _ := cmd.Flags().GetString("out")
			if dst == "" {
				dst = args[0]
			}
			return runPrefixAdd(args[0], dst, args[1])
		},
	}
	addCmd.Flags().String("out", "", "destination snapshot name (default: overwrite source)")

	removeCmd := &cobra.Command{
		Use:   "remove <snapshot> <prefix>",
		Short: "Remove prefix from matching keys in a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dst, _ := cmd.Flags().GetString("out")
			if dst == "" {
				dst = args[0]
			}
			return runPrefixRemove(args[0], dst, args[1])
		},
	}
	removeCmd.Flags().String("out", "", "destination snapshot name (default: overwrite source)")

	prefixCmd.AddCommand(addCmd, removeCmd)
	rootCmd.AddCommand(prefixCmd)
}

func runPrefixAdd(src, dst, p string) error {
	mgr, err := profile.NewManager()
	if err != nil {
		return err
	}
	res, err := prefix.RunAdd(mgr, src, dst, p)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, prefix.Format(res))
	return nil
}

func runPrefixRemove(src, dst, p string) error {
	mgr, err := profile.NewManager()
	if err != nil {
		return err
	}
	res, err := prefix.RunRemove(mgr, src, dst, p)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, prefix.Format(res))
	return nil
}
