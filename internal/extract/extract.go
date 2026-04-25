// Package extract provides functionality to extract a subset of keys
// from a snapshot into a new snapshot.
package extract

import (
	"errors"
	"fmt"
)

// Manager defines the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Options controls the behaviour of Run.
type Options struct {
	// Src is the source snapshot name.
	Src string
	// Dst is the destination snapshot name.
	Dst string
	// Keys is the list of keys to extract. If empty, ErrNoKeys is returned.
	Keys []string
	// Overwrite allows an existing destination snapshot to be replaced.
	Overwrite bool
}

var (
	ErrNoKeys      = errors.New("extract: no keys specified")
	ErrKeyNotFound = errors.New("extract: key not found in source snapshot")
)

// Run loads Src, picks out the requested Keys, and saves them as Dst.
// It returns the extracted map so callers can display or inspect it.
func Run(m Manager, opts Options) (map[string]string, error) {
	if len(opts.Keys) == 0 {
		return nil, ErrNoKeys
	}

	src, err := m.Load(opts.Src)
	if err != nil {
		return nil, fmt.Errorf("extract: load %q: %w", opts.Src, err)
	}

	out := make(map[string]string, len(opts.Keys))
	for _, k := range opts.Keys {
		v, ok := src[k]
		if !ok {
			return nil, fmt.Errorf("%w: %q", ErrKeyNotFound, k)
		}
		out[k] = v
	}

	if err := m.Save(opts.Dst, out); err != nil {
		return nil, fmt.Errorf("extract: save %q: %w", opts.Dst, err)
	}

	return out, nil
}
