// Package sample provides functionality to randomly sample a subset of
// environment variable keys from a snapshot.
package sample

import (
	"fmt"
	"math/rand"
	"sort"
)

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a sample operation.
type Result struct {
	Source string
	Dest   string
	Keys   []string
}

// Run samples n keys from src and saves them as dest.
// If seed is non-zero it is used to seed the random source for reproducibility.
func Run(m Manager, src, dest string, n int, seed int64) (Result, error) {
	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("sample: load %q: %w", src, err)
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if n <= 0 || n > len(keys) {
		n = len(keys)
	}

	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	sampled := keys[:n]
	sort.Strings(sampled)

	out := make(map[string]string, len(sampled))
	for _, k := range sampled {
		out[k] = vars[k]
	}

	if err := m.Save(dest, out); err != nil {
		return Result{}, fmt.Errorf("sample: save %q: %w", dest, err)
	}

	return Result{Source: src, Dest: dest, Keys: sampled}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	return fmt.Sprintf("sampled %d key(s) from %q into %q", len(r.Keys), r.Source, r.Dest)
}
