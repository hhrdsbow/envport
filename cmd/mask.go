package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envport/internal/mask"
	"envport/internal/profile"
)

func init() {
	var (
		visible int
		maskLen int
		keys    []string
	)

	cmd := &cobra.Command{
		Use:   "mask <profile> <snapshot>",
		Short: "Print snapshot variables with sensitive values masked",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMask(args[0], args[1], visible, maskLen, keys)
		},
	}

	cmd.Flags().IntVar(&visible, "visible", mask.DefaultVisibleChars,
		"number of leading characters to leave unmasked")
	cmd.Flags().IntVar(&maskLen, "mask-len", 8,
		"fixed length of the mask suffix (0 = match actual value length)")
	cmd.Flags().StringSliceVar(&keys, "keys", nil,
		"only mask these keys (comma-separated); masks all if omitted")

	rootCmd.AddCommand(cmd)
}

func runMask(profileName, snapshotName string, visible, maskLen int, keys []string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	snap, err := mgr.Load(profileName, snapshotName)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	opts := mask.Options{
		VisibleChars: visible,
		MaskChar:     '*',
		MaskLen:      maskLen,
	}

	var masked map[string]string
	if len(keys) > 0 {
		masked = mask.ApplyKeys(snap.Vars, keys, opts)
	} else {
		masked = mask.Apply(snap.Vars, opts)
	}

	sorted := make([]string, 0, len(masked))
	for k := range masked {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	w := os.Stdout
	for _, k := range sorted {
		fmt.Fprintf(w, "%s=%s\n", k, masked[k])
	}
	return nil
}
