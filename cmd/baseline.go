package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"envport/internal/baseline"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func init() {
	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Manage baseline snapshots for profiles",
	}

	setCmd := &cobra.Command{
		Use:   "set <profile> <snapshot>",
		Short: "Mark a snapshot as the baseline for a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBaselineSet(args[0], args[1])
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <profile>",
		Short: "Show the current baseline for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBaselineGet(args[0])
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <profile>",
		Short: "Remove the baseline for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBaselineClear(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all baselines",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBaselineList()
		},
	}

	baselineCmd.AddCommand(setCmd, getCmd, clearCmd, listCmd)
	rootCmd.AddCommand(baselineCmd)
}

func openBaselineManager() (*baseline.Manager, error) {
	s, err := store.Open("baseline")
	if err != nil {
		return nil, err
	}
	return baseline.New(s), nil
}

func runBaselineSet(profile, snapshot string) error {
	m, err := openBaselineManager()
	if err != nil {
		return err
	}
	if err := m.Set(profile, snapshot); err != nil {
		return err
	}
	fmt.Printf("baseline set: %s -> %s\n", profile, snapshot)
	return nil
}

func runBaselineGet(profile string) error {
	m, err := openBaselineManager()
	if err != nil {
		return err
	}
	e, err := m.Get(profile)
	if errors.Is(err, baseline.ErrNotFound) {
		fmt.Printf("no baseline set for profile %q\n", profile)
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Printf("profile: %s\nsnapshot: %s\nset at: %s\n", e.Profile, e.Snapshot, e.SetAt.Format("2006-01-02 15:04:05"))
	return nil
}

func runBaselineClear(profile string) error {
	m, err := openBaselineManager()
	if err != nil {
		return err
	}
	if err := m.Clear(profile); err != nil {
		return err
	}
	fmt.Printf("baseline cleared for profile %q\n", profile)
	return nil
}

func runBaselineList() error {
	m, err := openBaselineManager()
	if err != nil {
		return err
	}
	entries, err := m.List()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("no baselines set")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROFILE\tSNAPSHOT\tSET AT")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Profile, e.Snapshot, e.SetAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
