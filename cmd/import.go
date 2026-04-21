package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	importpkg "envport/internal/import"
	"envport/internal/profile"

	"github.com/spf13/cobra"
)

func init() {
	var format string
	var profileName string

	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import environment variables from a file into a named snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImport(args[0], profileName, format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "", "File format: shell, dotenv, json (auto-detected if omitted)")
	cmd.Flags().StringVarP(&profileName, "name", "n", "", "Snapshot name (defaults to filename without extension)")

	rootCmd.AddCommand(cmd)
}

func runImport(filePath, profileName, format string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if format == "" {
		format = detectFormat(filePath)
	}

	vars, err := importpkg.Parse(string(data), format)
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	if profileName == "" {
		base := filepath.Base(filePath)
		ext := filepath.Ext(base)
		profileName = base[:len(base)-len(ext)]
	}

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	if err := mgr.Save(profileName, vars); err != nil {
		return fmt.Errorf("saving snapshot %q: %w", profileName, err)
	}

	fmt.Printf("Imported %d variable(s) into snapshot %q\n", len(vars), profileName)
	return nil
}

func detectFormat(path string) string {
	switch filepath.Ext(path) {
	case ".json":
		return "json"
	case ".env":
		return "dotenv"
	default:
		return "shell"
	}
}
