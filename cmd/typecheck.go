package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"envport/internal/profile"
	"envport/internal/typecheck"
)

func init() {
	var types []string

	cmd := &cobra.Command{
		Use:   "typecheck <snapshot>",
		Short: "Validate environment variable types in a snapshot",
		Long: `Check that values in a snapshot conform to declared types.

Supported types: int, float, bool, url

Example:
  envport typecheck prod --type PORT=int --type BASE_URL=url`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTypecheck(cmd, args, types)
		},
	}

	cmd.Flags().StringArrayVar(&types, "type", nil, "key=type pair (repeatable)")
	_ = cmd.MarkFlagRequired("type")

	rootCmd.AddCommand(cmd)
}

func runTypecheck(cmd *cobra.Command, args []string, rawTypes []string) error {
	name := args[0]

	typeMap, err := parseTypeFlags(rawTypes)
	if err != nil {
		return err
	}

	storeDir := filepath.Join(os.Getenv("HOME"), ".envport")
	mgr, err := profile.New(storeDir)
	if err != nil {
		return fmt.Errorf("open profile store: %w", err)
	}

	res, err := typecheck.Run(mgr, name, typeMap)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), typecheck.Format(res))

	if len(res.Violations) > 0 {
		return fmt.Errorf("typecheck failed with %d violation(s)", len(res.Violations))
	}
	return nil
}

func parseTypeFlags(raw []string) (map[string]typecheck.Type, error) {
	out := make(map[string]typecheck.Type, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid --type flag %q: expected key=type", r)
		}
		t := typecheck.Type(strings.ToLower(parts[1]))
		switch t {
		case typecheck.TypeInt, typecheck.TypeFloat, typecheck.TypeBool, typecheck.TypeURL:
		default:
			return nil, fmt.Errorf("unknown type %q for key %q", parts[1], parts[0])
		}
		out[parts[0]] = t
	}
	return out, nil
}
