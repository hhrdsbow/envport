package cmd

import (
	"fmt"
	"os"
	"strings"

	"envport/internal/profile"
	"envport/internal/template"

	"github.com/spf13/cobra"
)

func init() {
	var src, dst string
	var pairs []string

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Apply placeholder substitution to a snapshot",
		Long: `Load a snapshot, replace {{KEY}} tokens with supplied values, and save as a new snapshot.

Example:
  envport template --src base --dst prod --set HOST=db-prod --set ENV=production`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplate(src, dst, pairs)
		},
	}

	cmd.Flags().StringVar(&src, "src", "", "source snapshot name (required)")
	cmd.Flags().StringVar(&dst, "dst", "", "destination snapshot name (required)")
	cmd.Flags().StringArrayVar(&pairs, "set", nil, "placeholder=value pairs")
	_ = cmd.MarkFlagRequired("src")
	_ = cmd.MarkFlagRequired("dst")

	rootCmd.AddCommand(cmd)
}

func runTemplate(src, dst string, pairs []string) error {
	dir, err := storageDir()
	if err != nil {
		return err
	}

	mgr, err := profile.New(dir)
	if err != nil {
		return err
	}

	placeholders := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --set value %q: expected KEY=VALUE", p)
		}
		placeholders[parts[0]] = parts[1]
	}

	res, err := template.Run(mgr, src, dst, placeholders)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, template.Format(res))
	return nil
}
