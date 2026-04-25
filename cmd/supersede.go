package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/profile"
	"github.com/user/envport/internal/supersede"
)

func init() {
	var dest string
	var overrideProfile string
	var rawKeys []string

	cmd := &cobra.Command{
		Use:   "supersede <base>",
		Short: "Override specific keys in a snapshot",
		Long: `Load a snapshot and apply key overrides on top of it.

Overrides can come from explicit KEY=VALUE pairs (--set) or from another
snapshot (--from). The result is saved to --dest (defaults to <base>).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSupersede(cmd, args, dest, overrideProfile, rawKeys)
		},
	}

	cmd.Flags().StringVar(&dest, "dest", "", "destination snapshot name (default: overwrite base)")
	cmd.Flags().StringVar(&overrideProfile, "from", "", "snapshot to pull override values from")
	cmd.Flags().StringArrayVar(&rawKeys, "set", nil, "KEY=VALUE pairs to override (repeatable)")

	rootCmd.AddCommand(cmd)
}

func runSupersede(cmd *cobra.Command, args []string, dest, overrideProfile string, rawKeys []string) error {
	base := args[0]

	m, err := profile.NewManager()
	if err != nil {
		return fmt.Errorf("supersede: open profile manager: %w", err)
	}

	keys := make(map[string]string, len(rawKeys))
	for _, pair := range rawKeys {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("supersede: invalid KEY=VALUE pair: %q", pair)
		}
		keys[parts[0]] = parts[1]
	}

	res, err := supersede.Run(m, base, supersede.Options{
		Keys:            keys,
		OverrideProfile: overrideProfile,
		Dest:            dest,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Saved to %q with %d key(s) overridden.\n", res.Dest, len(res.Applied))
	for k, v := range res.Applied {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, v)
	}
	return nil
}
