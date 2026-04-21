// Package cmd provides the CLI commands for envport.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envport",
	Short: "Snapshot and restore environment variable sets across projects",
	Long: `envport is a CLI tool for capturing, storing, and restoring
environment variable sets. It supports named profiles, tagging,
diffing, exporting, encrypting, and more.

Examples:
  envport snapshot myproject
  envport restore myproject
  envport diff myproject
  envport export myproject --format dotenv
  envport list`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
