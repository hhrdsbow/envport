// Package select provides functionality to select a subset of keys
// from a snapshot by explicit name or glob pattern.
package select

import (
	"fmt"
	"path"
	"sort"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a Select operation.
type Result struct {
	Kept    []string
	Dropped []string
}

// Run loads src, keeps only keys matching any of the given patterns
// (exact name or glob), and saves the result to dst.
func Run(m Manager, src, dst string, patterns []string) (Result, error) {
	if len(patterns) == 0 {
		return Result{}, fmt.Errorf("select: at least one key pattern is required")
	}

	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("select: load %q: %w", src, err)
	}

	selected := make(map[string]string)
	var kept, dropped []string

	for k, v := range vars {
		if matchesAny(k, patterns) {
			selected[k] = v
			kept = append(kept, k)
		} else {
			dropped = append(dropped, k)
		}
	}

	if err := m.Save(dst, selected); err != nil {
		return Result{}, fmt.Errorf("select: save %q: %w", dst, err)
	}

	sort.Strings(kept)
	sort.Strings(dropped)
	return Result{Kept: kept, Dropped: dropped}, nil
}

// Format returns a human-readable summary of the result.
func Format(r Result) string {
	out := fmt.Sprintf("kept %d key(s), dropped %d key(s)\n", len(r.Kept), len(r.Dropped))
	for _, k := range r.Kept {
		out += fmt.Sprintf("  + %s\n", k)
	}
	for _, k := range r.Dropped {
		out += fmt.Sprintf("  - %s\n", k)
	}
	return out
}

func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		if ok, _ := path.Match(p, key); ok {
			return true
		}
		if key == p {
			return true
		}
	}
	return false
}
