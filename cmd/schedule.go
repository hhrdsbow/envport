package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"envport/internal/schedule"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func init() {
	schedCmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage snapshot schedules",
	}

	addCmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add or update a snapshot schedule",
		Args:  cobra.ExactArgs(1),
		RunE:  runScheduleAdd,
	}
	addCmd.Flags().String("profile", "", "Profile to snapshot (required)")
	addCmd.Flags().Duration("interval", 24*time.Hour, "Snapshot interval (e.g. 1h, 24h)")
	_ = addCmd.MarkFlagRequired("profile")

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a snapshot schedule",
		Args:  cobra.ExactArgs(1),
		RunE:  runScheduleRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all snapshot schedules",
		RunE:  runScheduleList,
	}

	schedCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(schedCmd)
}

func openScheduleManager() (*schedule.Manager, error) {
	s, err := store.Open("schedules")
	if err != nil {
		return nil, err
	}
	return schedule.New(s), nil
}

func runScheduleAdd(cmd *cobra.Command, args []string) error {
	profile, _ := cmd.Flags().GetString("profile")
	interval, _ := cmd.Flags().GetDuration("interval")
	m, err := openScheduleManager()
	if err != nil {
		return err
	}
	e := schedule.Entry{Name: args[0], Profile: profile, Interval: interval}
	if err := m.Add(e); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "schedule %q added (profile=%s interval=%s)\n", args[0], profile, interval)
	return nil
}

func runScheduleRemove(cmd *cobra.Command, args []string) error {
	m, err := openScheduleManager()
	if err != nil {
		return err
	}
	if err := m.Remove(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "schedule %q removed\n", args[0])
	return nil
}

func runScheduleList(cmd *cobra.Command, args []string) error {
	m, err := openScheduleManager()
	if err != nil {
		return err
	}
	entries, err := m.List()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no schedules defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tPROFILE\tINTERVAL\tLAST RUN\tENABLED")
	for _, e := range entries {
		last := "never"
		if !e.LastRun.IsZero() {
			last = e.LastRun.Format(time.RFC3339)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", e.Name, e.Profile, e.Interval, last, e.Enabled)
	}
	return w.Flush()
}
