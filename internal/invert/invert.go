// Package invert swaps keys and values in a snapshot.
package invert

import "fmt"

// Manager describes the snapshot store operations required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of an invert operation.
type Result struct {
	Source      string
	Destination string
	Count       int
	Conflicts   []string
}

// Run loads the snapshot named src, swaps every key↔value pair, and saves
// the result as dst. If two source values are identical the last key wins
// and the conflict is recorded in Result.Conflicts.
func Run(m Manager, src, dst string) (Result, error) {
	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("invert: load %q: %w", src, err)
	}

	inverted := make(map[string]string, len(vars))
	seen := make(map[string]string, len(vars)) // value → original key
	var conflicts []string

	for k, v := range vars {
		if prev, exists := seen[v]; exists {
			conflicts = append(conflicts, fmt.Sprintf("%s and %s share value %q", prev, k, v))
		}
		seen[v] = k
		inverted[v] = k
	}

	if err := m.Save(dst, inverted); err != nil {
		return Result{}, fmt.Errorf("invert: save %q: %w", dst, err)
	}

	return Result{
		Source:      src,
		Destination: dst,
		Count:       len(inverted),
		Conflicts:   conflicts,
	}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	base := fmt.Sprintf("inverted %d key(s) from %q into %q", r.Count, r.Source, r.Destination)
	if len(r.Conflicts) == 0 {
		return base
	}
	s := base + fmt.Sprintf(" (%d conflict(s)):\n", len(r.Conflicts))
	for _, c := range r.Conflicts {
		s += "  conflict: " + c + "\n"
	}
	return s
}
