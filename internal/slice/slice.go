// Package slice extracts a contiguous range of key-value pairs from a snapshot
// after sorting keys alphabetically, similar to array slicing.
package slice

import (
	"fmt"
	"sort"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the output of a slice operation.
type Result struct {
	Keys   []string
	Vars   map[string]string
	Src    string
	Dst    string
}

// Run extracts keys at positions [from, to) (0-based) from the src snapshot
// sorted alphabetically and saves them into dst.
func Run(m Manager, src, dst string, from, to int) (Result, error) {
	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("slice: load %q: %w", src, err)
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if from < 0 {
		from = 0
	}
	if to < 0 || to > len(keys) {
		to = len(keys)
	}
	if from > len(keys) {
		from = len(keys)
	}
	if from > to {
		return Result{}, fmt.Errorf("slice: from (%d) > to (%d)", from, to)
	}

	selected := keys[from:to]
	out := make(map[string]string, len(selected))
	for _, k := range selected {
		out[k] = vars[k]
	}

	if dst != "" {
		if err := m.Save(dst, out); err != nil {
			return Result{}, fmt.Errorf("slice: save %q: %w", dst, err)
		}
	}

	return Result{Keys: selected, Vars: out, Src: src, Dst: dst}, nil
}

// Format returns a human-readable summary of the slice result.
func Format(r Result) string {
	if len(r.Keys) == 0 {
		return "slice: no keys selected"
	}
	s := fmt.Sprintf("slice: %d key(s) from %q", len(r.Keys), r.Src)
	if r.Dst != "" {
		s += fmt.Sprintf(" → saved to %q", r.Dst)
	}
	return s
}
