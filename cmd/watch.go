package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/profile"
	"envport/internal/watch"

	"github.com/spf13/cobra"
)

func init() {
	var keys string
	var quiet bool

	cmd := &cobra.Command{
		Use:   "watch <snapshot>",
		Short: "Check if the current environment has drifted from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(args[0], keys, quiet)
		},
	}

	cmd.Flags().StringVarP(&keys, "keys", "k", "", "Comma-separated list of keys to watch")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Exit non-zero silently if drift detected")

	rootCmd.AddCommand(cmd)
}

func runWatch(name, keys string, quiet bool) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	m, err := profile.New(dir)
	if err != nil {
		return err
	}

	var keyList []string
	if keys != "" {
		for _, k := range strings.Split(keys, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keyList = append(keyList, k)
			}
		}
	}

	res, err := watch.Check(m, name, keyList)
	if err != nil {
		return err
	}

	if !res.Changed {
		if !quiet {
			fmt.Printf("snapshot %q: no drift detected\n", name)
		}
		return nil
	}

	if !quiet {
		fmt.Printf("snapshot %q: drift detected\n%s\n", name, res.Report)
	}
	os.Exit(1)
	return nil
}
