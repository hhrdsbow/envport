package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"envport/internal/history"
	"envport/internal/store"

	"github.com/spf13/cobra"
)

func init() {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show operation history",
		RunE:  runHistory,
	}
	historyCmd.Flags().Bool("clear", false, "Clear all history")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	s, err := store.Open(defaultStorePath())
	if err != nil {
		return err
	}
	mgr := history.New(s)

	clear, _ := cmd.Flags().GetBool("clear")
	if clear {
		if err := mgr.Clear(); err != nil {
			return err
		}
		fmt.Println("history cleared")
		return nil
	}

	entries, err := mgr.List()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("no history recorded")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tOPERATION\tSNAPSHOT\tDETAIL")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format(time.RFC3339), e.Operation, e.Snapshot, e.Detail)
	}
	return w.Flush()
}
