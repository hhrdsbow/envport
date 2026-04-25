// Package subset extracts a named subset of keys from a snapshot into a new snapshot.
package subset

import "fmt"

// Manager defines the storage operations required by subset.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a subset operation.
type Result struct {
	Source      string
	Destination string
	Keys        []string
	Extracted   map[string]string
	Missing     []string
}

// Run loads src, picks the requested keys, and saves them as dst.
// Keys that do not exist in src are recorded in Result.Missing but do not
// cause an error unless all requested keys are absent.
func Run(m Manager, src, dst string, keys []string) (Result, error) {
	if len(keys) == 0 {
		return Result{}, fmt.Errorf("subset: at least one key must be specified")
	}

	vars, err := m.Load(src)
	if err != nil {
		return Result{}, fmt.Errorf("subset: load %q: %w", src, err)
	}

	extracted := make(map[string]string, len(keys))
	var missing []string
	for _, k := range keys {
		v, ok := vars[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		extracted[k] = v
	}

	if len(extracted) == 0 {
		return Result{}, fmt.Errorf("subset: none of the requested keys found in %q", src)
	}

	if err := m.Save(dst, extracted); err != nil {
		return Result{}, fmt.Errorf("subset: save %q: %w", dst, err)
	}

	return Result{
		Source:      src,
		Destination: dst,
		Keys:        keys,
		Extracted:   extracted,
		Missing:     missing,
	}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	out := fmt.Sprintf("subset: %d key(s) copied from %q to %q", len(r.Extracted), r.Source, r.Destination)
	if len(r.Missing) > 0 {
		out += fmt.Sprintf(" (%d missing: %v)", len(r.Missing), r.Missing)
	}
	return out
}
