// Package promote copies a snapshot from one profile to another,
// optionally overwriting the destination.
package promote

import (
	"errors"
	"fmt"

	"github.com/nicholasgasior/envport/internal/snapshot"
)

// ErrDestinationExists is returned when the destination profile already has a
// snapshot with the same name and --force was not supplied.
var ErrDestinationExists = errors.New("snapshot already exists in destination profile")

// Manager describes the persistence operations required by Run.
type Manager interface {
	Load(profile, name string) (*snapshot.Snapshot, error)
	Save(profile string, snap *snapshot.Snapshot) error
	Exists(profile, name string) bool
}

// Options controls the behaviour of Run.
type Options struct {
	SrcProfile  string
	DstProfile  string
	Name        string
	Force       bool
}

// Run promotes a named snapshot from the source profile to the destination
// profile. If the snapshot already exists in the destination and Force is
// false, ErrDestinationExists is returned.
func Run(m Manager, opts Options) error {
	if opts.SrcProfile == opts.DstProfile {
		return fmt.Errorf("source and destination profiles must differ")
	}

	snap, err := m.Load(opts.SrcProfile, opts.Name)
	if err != nil {
		return fmt.Errorf("load source snapshot: %w", err)
	}

	if !opts.Force && m.Exists(opts.DstProfile, opts.Name) {
		return ErrDestinationExists
	}

	// Preserve name; update profile association via a shallow copy.
	dst := &snapshot.Snapshot{
		Name: snap.Name,
		Vars: snap.Vars,
		CreatedAt: snap.CreatedAt,
	}

	if err := m.Save(opts.DstProfile, dst); err != nil {
		return fmt.Errorf("save destination snapshot: %w", err)
	}
	return nil
}
