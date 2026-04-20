package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envport/internal/profile"
	"envport/internal/validate"
)

var (
	validateRequiredKeys []string
	validateAllowEmpty   bool
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate <profile>",
		Short: "Validate that a snapshot contains required keys",
		Long: `Check a saved snapshot against a set of required keys.

Exits with a non-zero status if any required keys are missing or,
when --no-empty is set, if any values are empty strings.`,
		Args: cobra.ExactArgs(1),
		RunE: runValidate,
	}

	validateCmd.Flags().StringSliceVarP(
		&validateRequiredKeys,
		"keys", "k", nil,
		"Comma-separated list of required keys (e.g. DATABASE_URL,API_KEY)",
	)
	validateCmd.Flags().BoolVar(
		&validateAllowEmpty,
		"no-empty", false,
		"Treat keys with empty values as missing",
	)

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	if len(validateRequiredKeys) == 0 {
		return fmt.Errorf("at least one required key must be specified via --keys")
	}

	// Trim any accidental whitespace from key names.
	trimmed := make([]string, 0, len(validateRequiredKeys))
	for _, k := range validateRequiredKeys {
		if t := strings.TrimSpace(k); t != "" {
			trimmed = append(trimmed, t)
		}
	}

	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	mgr, err := profile.NewManager(dir)
	if err != nil {
		return fmt.Errorf("opening profile store: %w", err)
	}

	result, err := validate.Run(mgr, profileName, trimmed, validate.Options{
		RejectEmpty: validateAllowEmpty,
	})
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), validate.Format(result))

	if !result.OK {
		os.Exit(1)
	}

	return nil
}
