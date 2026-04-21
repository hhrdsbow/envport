package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nicholasgasior/envport/internal/profile"
	"github.com/nicholasgasior/envport/internal/promote"
	"github.com/spf13/cobra"
)

func init() {
	var srcProfile string
	var force bool

	cmd := &cobra.Command{
		Use:   "promote <name> --from <src-profile> <dst-profile>",
		Short: "Promote a snapshot from one profile to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			dstProfile := args[1]
			return runPromote(name, srcProfile, dstProfile, force)
		},
	}

	cmd.Flags().StringVar(&srcProfile, "from", "", "Source profile (required)")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite if snapshot already exists in destination")
	_ = cmd.MarkFlagRequired("from")

	rootCmd.AddCommand(cmd)
}

func runPromote(name, srcProfile, dstProfile string, force bool) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("open profile manager: %w", err)
	}

	opts := promote.Options{
		SrcProfile: srcProfile,
		DstProfile: dstProfile,
		Name:       name,
		Force:      force,
	}

	if err := promote.Run(m, opts); err != nil {
		if errors.Is(err, promote.ErrDestinationExists) {
			fmt.Fprintf(os.Stderr, "error: snapshot %q already exists in profile %q (use --force to overwrite)\n", name, dstProfile)
			return err
		}
		return err
	}

	fmt.Printf("promoted %q from %q to %q\n", name, srcProfile, dstProfile)
	return nil
}
