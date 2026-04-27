// Package split partitions an existing snapshot into multiple child snapshots
// based on explicit key assignments, with an optional remainder snapshot for
// any keys not assigned to a named target.
package split
