// Package clone provides functionality to duplicate an existing snapshot
// under a new name, automatically resolving naming conflicts by appending
// a configurable suffix.
//
// Basic usage:
//
//	// Clone a snapshot named "production" to "production-copy"
//	err := clone.Snapshot(ctx, store, "production", "production-copy", clone.Options{})
//
// If the destination name already exists, the package can automatically
// append a numeric suffix (e.g. "production-copy-2") to avoid conflicts.
package clone
