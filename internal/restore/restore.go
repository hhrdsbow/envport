package restore

import (
	"fmt"

	"envport/internal/env"
	"envport/internal/profile"
	"envport/internal/snapshot"
)

// Options configures a restore operation.
type Options struct {
	ProfileName  string
	SnapshotName string
	DryRun       bool
	FilterKeys   []string
}

// Result holds the outcome of a restore operation.
type Result struct {
	Applied  map[string]string
	Skipped  []string
	DryRun   bool
}

// Run restores environment variables from a saved snapshot via a profile manager.
func Run(mgr *profile.Manager, opts Options) (*Result, error) {
	snap, err := mgr.Load(opts.ProfileName, opts.SnapshotName)
	if err != nil {
		return nil, fmt.Errorf("restore: load snapshot %q/%q: %w", opts.ProfileName, opts.SnapshotName, err)
	}

	vars := snap.Vars
	if len(opts.FilterKeys) > 0 {
		vars = env.FilterKeys(vars, opts.FilterKeys)
	}

	skipped := []string{}
	for k := range snap.Vars {
		if _, ok := vars[k]; !ok {
			skipped = append(skipped, k)
		}
	}

	if !opts.DryRun {
		if err := env.Apply(vars); err != nil {
			return nil, fmt.Errorf("restore: apply env: %w", err)
		}
	}

	return &Result{
		Applied: vars,
		Skipped: skipped,
		DryRun:  opts.DryRun,
	}, nil
}

// ExportScript returns a shell export script for the snapshot variables.
func ExportScript(snap *snapshot.Snapshot, filter []string) string {
	vars := snap.Vars
	if len(filter) > 0 {
		vars = env.FilterKeys(vars, filter)
	}
	return snapshot.ToExports(vars)
}
