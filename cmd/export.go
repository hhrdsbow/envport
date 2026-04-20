package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envport/internal/export"
	"envport/internal/profile"
)

func init() {
	var format string

	cmd := &cobra.Command{
		Use:   "export <profile>",
		Short: "Export a saved profile in shell, dotenv, or json format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(args[0], export.Format(format))
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "shell", "output format: shell, dotenv, json")
	rootCmd.AddCommand(cmd)
}

func runExport(name string, format export.Format) error {
	if !export.IsValidFormat(format) {
		return fmt.Errorf("unknown format %q: must be one of shell, dotenv, json", format)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home dir: %w", err)
	}

	mgr, err := profile.New(home)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	snap, err := mgr.Load(name)
	if err != nil {
		return fmt.Errorf("load profile %q: %w", name, err)
	}

	if err := export.Write(os.Stdout, snap.Vars, format); err != nil {
		return fmt.Errorf("write export: %w", err)
	}
	return nil
}
