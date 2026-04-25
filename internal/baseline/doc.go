// Package baseline tracks which snapshot is considered the reference point
// (baseline) for a given profile, supporting drift detection workflows.
//
// A baseline is a named snapshot that serves as the "known good" state for
// a profile. When new snapshots are taken, they can be compared against the
// baseline to detect configuration drift.
//
// Typical usage:
//
//	// Set the current snapshot as the baseline for a profile
//	 err := baseline.Set(ctx, store, profileID, snapshotID)
//
//	// Retrieve the baseline snapshot for comparison
//	 snap, err := baseline.Get(ctx, store, profileID)
package baseline
