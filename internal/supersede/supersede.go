// Package supersede provides functionality to override specific keys in a
// snapshot with values from another snapshot or an explicit map.
package supersede

import (
	"errors"
	"fmt"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Options controls how supersede behaves.
type Options struct {
	// Keys is an explicit set of key=value pairs to apply on top of Base.
	// If empty, OverrideProfile is used instead.
	Keys map[string]string

	// OverrideProfile is the name of a snapshot whose values win over Base.
	// Ignored when Keys is non-empty.
	OverrideProfile string

	// Dest is the snapshot name to write the result into.
	// If empty, Base is overwritten.
	Dest string
}

// Result describes what Run produced.
type Result struct {
	Dest    string
	Applied map[string]string // keys that were overridden and their new values
}

// Run loads Base, applies overrides according to opts, and saves the result.
func Run(m Manager, base string, opts Options) (Result, error) {
	if base == "" {
		return Result{}, errors.New("supersede: base profile name is required")
	}

	vars, err := m.Load(base)
	if err != nil {
		return Result{}, fmt.Errorf("supersede: load base %q: %w", base, err)
	}

	overrides := opts.Keys
	if len(overrides) == 0 && opts.OverrideProfile != "" {
		overrides, err = m.Load(opts.OverrideProfile)
		if err != nil {
			return Result{}, fmt.Errorf("supersede: load override profile %q: %w", opts.OverrideProfile, err)
		}
	}

	applied := make(map[string]string, len(overrides))
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
		applied[k] = v
	}

	dest := opts.Dest
	if dest == "" {
		dest = base
	}

	if err := m.Save(dest, out); err != nil {
		return Result{}, fmt.Errorf("supersede: save %q: %w", dest, err)
	}

	return Result{Dest: dest, Applied: applied}, nil
}
