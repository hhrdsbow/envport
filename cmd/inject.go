package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/inject"
	"github.com/user/envport/internal/profile"
)

func init() {
	var keys []string
	var override bool

	cmd := &cobra.Command{
		Use:   "inject <snapshot> -- <command> [args...]",
		Short: "Run a command with a snapshot's variables injected",
		Long: `inject loads the named snapshot and executes the given command with
those variables present in the child process environment.

By default the snapshot variables are appended to the inherited environment.
Use --override to replace any conflicting keys from the parent process.`,
		Args:               cobra.MinimumNArgs(2),
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInject(args[0], args[1:], keys, override)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "comma-separated list of keys to inject (default: all)")
	cmd.Flags().BoolVar(&override, "override", false, "override parent env vars with snapshot values")

	rootCmd.AddCommand(cmd)
}

func runInject(snapName string, cmdArgs []string, keys []string, override bool) error {
	if len(cmdArgs) == 0 {
		return fmt.Errorf("no command specified after snapshot name")
	}

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	snap, err := mgr.Load(snapName)
	if err != nil {
		return fmt.Errorf("snapshot %q not found: %w", snapName, err)
	}

	vars := snap.Vars
	if len(keys) > 0 {
		filtered := make(map[string]string, len(keys))
		for _, k := range keys {
			if v, ok := vars[k]; ok {
				filtered[k] = v
			}
		}
		vars = filtered
	}

	env := inject.BuildEnv(vars, override)

	c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	c.Env = env
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		// Preserve exit code when possible.
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "inject: command exited with status %d\n", exitErr.ExitCode())
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("inject: %w", err)
	}

	_ = strings.Join // keep import tidy
	return nil
}
