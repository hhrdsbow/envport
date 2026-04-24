package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/envport/envport/internal/quota"
	"github.com/envport/envport/internal/store"
	"github.com/spf13/cobra"
)

func openQuotaManager() (*quota.Manager, error) {
	s, err := store.Open("quota")
	if err != nil {
		return nil, err
	}
	return quota.New(s), nil
}

func init() {
	quotaCmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage per-profile variable count quotas",
	}

	setCmd := &cobra.Command{
		Use:   "set <profile> <limit>",
		Short: "Set the maximum number of variables allowed in a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuotaSet(args[0], args[1])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <profile>",
		Short: "Remove the quota for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuotaRemove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profile quotas",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuotaList()
		},
	}

	quotaCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(quotaCmd)
}

func runQuotaSet(profile, limitStr string) error {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return fmt.Errorf("limit must be a positive integer")
	}
	mgr, err := openQuotaManager()
	if err != nil {
		return err
	}
	if err := mgr.Set(profile, limit); err != nil {
		return err
	}
	fmt.Printf("quota set: %s → %d variables\n", profile, limit)
	return nil
}

func runQuotaRemove(profile string) error {
	mgr, err := openQuotaManager()
	if err != nil {
		return err
	}
	if err := mgr.Delete(profile); err != nil {
		return err
	}
	fmt.Printf("quota removed for profile: %s\n", profile)
	return nil
}

func runQuotaList() error {
	mgr, err := openQuotaManager()
	if err != nil {
		return err
	}
	all, err := mgr.List()
	if err != nil {
		return err
	}
	if len(all) == 0 {
		fmt.Println("no quotas configured")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROFILE\tLIMIT")
	for profile, limit := range all {
		fmt.Fprintf(w, "%s\t%d\n", profile, limit)
	}
	return w.Flush()
}
