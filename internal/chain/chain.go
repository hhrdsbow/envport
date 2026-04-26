// Package chain provides functionality to apply multiple snapshots in sequence,
// merging their variables with later snapshots taking precedence.
package chain

import (
	"errors"
	"fmt"
)

// Manager defines the interface for loading snapshots by name.
type Manager interface {
	Load(name string) (map[string]string, error)
}

// Result holds the merged environment and metadata about the chain.
type Result struct {
	Vars     map[string]string
	Applied  []string
	Skipped  []string
}

// Run applies a sequence of named snapshots in order. Later snapshots override
// earlier ones for conflicting keys. If a snapshot is missing and skipMissing
// is true, it is recorded in Result.Skipped rather than returning an error.
func Run(m Manager, names []string, skipMissing bool) (Result, error) {
	if len(names) == 0 {
		return Result{}, errors.New("chain: at least one snapshot name required")
	}

	merged := make(map[string]string)
	result := Result{Vars: merged}

	for _, name := range names {
		vars, err := m.Load(name)
		if err != nil {
			if skipMissing {
				result.Skipped = append(result.Skipped, name)
				continue
			}
			return Result{}, fmt.Errorf("chain: loading %q: %w", name, err)
		}
		for k, v := range vars {
			merged[k] = v
		}
		result.Applied = append(result.Applied, name)
	}

	return result, nil
}

// Format returns a human-readable summary of the chain result.
func Format(r Result) string {
	out := fmt.Sprintf("applied: %d snapshot(s), %d key(s) merged", len(r.Applied), len(r.Vars))
	if len(r.Skipped) > 0 {
		out += fmt.Sprintf(", skipped: %v", r.Skipped)
	}
	return out
}
