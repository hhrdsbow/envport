package cmd

import (
	"fmt"
	"os"

	"envport/internal/audit"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Show audit log of snapshot operations",
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().StringP("profile", "p", "", "filter by profile name")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	dir, err := defaultStoreDir()
	if err != nil {
		return err
	}

	st, err := store.Open(dir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	auditStore := newFileAuditStore(st)
	mgr := audit.New(auditStore)

	profileFilter, _ := cmd.Flags().GetString("profile")

	entries, err := mgr.List()
	if err != nil {
		return fmt.Errorf("list audit: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no audit entries found")
		return nil
	}

	for _, e := range entries {
		if profileFilter != "" && e.Profile != profileFilter {
			continue
		}
		fmt.Fprintln(os.Stdout, audit.Format(e))
	}
	return nil
}
