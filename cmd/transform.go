package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envport/internal/profile"
	"github.com/yourorg/envport/internal/transform"
)

func init() {
	var (
		prefixAdd    string
		prefixRemove string
		suffixAdd    string
		suffixRemove string
		uppercase    bool
		lowercase    bool
		dryRun       bool
	)

	cmd := &cobra.Command{
		Use:   "transform <snapshot>",
		Short: "Apply key/value transformations to a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTransform(args[0], prefixAdd, prefixRemove, suffixAdd, suffixRemove, uppercase, lowercase, dryRun)
		},
	}

	cmd.Flags().StringVar(&prefixAdd, "prefix-add", "", "Add prefix to all keys")
	cmd.Flags().StringVar(&prefixRemove, "prefix-remove", "", "Remove prefix from all keys")
	cmd.Flags().StringVar(&suffixAdd, "suffix-add", "", "Add suffix to all keys")
	cmd.Flags().StringVar(&suffixRemove, "suffix-remove", "", "Remove suffix from all keys")
	cmd.Flags().BoolVar(&uppercase, "uppercase", false, "Convert keys and values to uppercase")
	cmd.Flags().BoolVar(&lowercase, "lowercase", false, "Convert keys and values to lowercase")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without saving")

	rootCmd.AddCommand(cmd)
}

func runTransform(name, prefixAdd, prefixRemove, suffixAdd, suffixRemove string, uppercase, lowercase, dryRun bool) error {
	var ops []transform.Op
	if prefixAdd != "" {
		ops = append(ops, transform.Op{Kind: "prefix-add", Value: prefixAdd})
	}
	if prefixRemove != "" {
		ops = append(ops, transform.Op{Kind: "prefix-remove", Value: prefixRemove})
	}
	if suffixAdd != "" {
		ops = append(ops, transform.Op{Kind: "suffix-add", Value: suffixAdd})
	}
	if suffixRemove != "" {
		ops = append(ops, transform.Op{Kind: "suffix-remove", Value: suffixRemove})
	}
	if uppercase {
		ops = append(ops, transform.Op{Kind: "uppercase"})
	}
	if lowercase {
		ops = append(ops, transform.Op{Kind: "lowercase"})
	}
	if len(ops) == 0 {
		return fmt.Errorf("no transformation flags provided")
	}

	m, err := profile.New(profile.DefaultPath())
	if err != nil {
		return err
	}

	r, _, err := transform.Run(m, name, ops, dryRun)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, transform.Format(r, dryRun))
	return nil
}
