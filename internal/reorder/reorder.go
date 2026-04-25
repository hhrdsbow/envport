// Package reorder provides functionality to reorder environment variable
// keys within a snapshot according to a specified ordering strategy.
package reorder

import (
	"errors"
	"sort"

	"github.com/envport/envport/internal/snapshot"
)

// Strategy defines how keys should be reordered.
type Strategy string

const (
	StrategyAlpha  Strategy = "alpha"
	StrategyReverse Strategy = "reverse"
	StrategyCustom  Strategy = "custom"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (*snapshot.Snapshot, error)
	Save(name string, s *snapshot.Snapshot) error
}

// Options controls the behaviour of Run.
type Options struct {
	Strategy   Strategy
	CustomOrder []string // used when Strategy == StrategyCustom
}

// Run reorders the keys of the named snapshot in-place according to opts.
// The snapshot is saved back via the manager after reordering.
func Run(m Manager, name string, opts Options) ([]string, error) {
	if name == "" {
		return nil, errors.New("reorder: snapshot name must not be empty")
	}

	s, err := m.Load(name)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}

	switch opts.Strategy {
	case StrategyAlpha:
		sort.Strings(keys)
	case StrategyReverse:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case StrategyCustom:
		keys = applyCustomOrder(keys, opts.CustomOrder)
	default:
		return nil, errors.New("reorder: unknown strategy " + string(opts.Strategy))
	}

	if err := m.Save(name, s); err != nil {
		return nil, err
	}

	return keys, nil
}

// applyCustomOrder returns keys sorted so that those present in order come
// first (in the given order), followed by any remaining keys in alpha order.
func applyCustomOrder(keys []string, order []string) []string {
	index := make(map[string]int, len(order))
	for i, k := range order {
		index[k] = i
	}

	var pinned, rest []string
	for _, k := range keys {
		if _, ok := index[k]; ok {
			pinned = append(pinned, k)
		} else {
			rest = append(rest, k)
		}
	}

	sort.Slice(pinned, func(i, j int) bool {
		return index[pinned[i]] < index[pinned[j]]
	})
	sort.Strings(rest)

	return append(pinned, rest...)
}
