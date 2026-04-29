// Package sortkeys provides utilities for sorting snapshot keys
// by various strategies: alphabetical, by value length, or by key length.
package sortkeys

import (
	"fmt"
	"sort"
	"strings"
)

// Strategy defines how keys should be sorted.
type Strategy string

const (
	StrategyAlpha     Strategy = "alpha"
	StrategyValueLen  Strategy = "valuelen"
	StrategyKeyLen    Strategy = "keylen"
	StrategyReverse   Strategy = "reverse"
)

// Result holds the sorted key order and the original snapshot.
type Result struct {
	Keys     []string
	Vars     map[string]string
	Strategy Strategy
}

// Manager defines the interface for loading snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Run loads a snapshot and returns keys sorted by the given strategy.
func Run(mgr Manager, name string, strategy Strategy) (*Result, error) {
	vars, err := mgr.Load(name)
	if err != nil {
		return nil, fmt.Errorf("sortkeys: load %q: %w", name, err)
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	switch strategy {
	case StrategyAlpha, "":
		sort.Strings(keys)
	case StrategyKeyLen:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) != len(keys[j]) {
				return len(keys[i]) < len(keys[j])
			}
			return keys[i] < keys[j]
		})
	case StrategyValueLen:
		sort.Slice(keys, func(i, j int) bool {
			vi, vj := vars[keys[i]], vars[keys[j]]
			if len(vi) != len(vj) {
				return len(vi) < len(vj)
			}
			return keys[i] < keys[j]
		})
	case StrategyReverse:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	default:
		return nil, fmt.Errorf("sortkeys: unknown strategy %q", strategy)
	}

	return &Result{Keys: keys, Vars: vars, Strategy: strategy}, nil
}

// Format renders the sorted key=value pairs as a human-readable string.
func Format(r *Result) string {
	var sb strings.Builder
	for _, k := range r.Keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, r.Vars[k])
	}
	return sb.String()
}
