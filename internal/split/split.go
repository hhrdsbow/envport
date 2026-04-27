// Package split provides functionality to split a snapshot into multiple
// smaller snapshots based on key prefixes or explicit key lists.
package split

import "fmt"

// Manager defines the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a split operation.
type Result struct {
	Name string
	Keys []string
}

// Run splits the snapshot identified by src into one new snapshot per entry
// in targets. Each entry maps a destination snapshot name to the list of keys
// that should be moved into it. Keys not listed in any target remain in a
// remainder snapshot named src+suffix (empty string disables remainder).
func Run(m Manager, src string, targets map[string][]string, remainder string) ([]Result, error) {
	vars, err := m.Load(src)
	if err != nil {
		return nil, fmt.Errorf("split: load %q: %w", src, err)
	}

	used := make(map[string]bool)
	var results []Result

	for dest, keys := range targets {
		slice := make(map[string]string, len(keys))
		for _, k := range keys {
			if v, ok := vars[k]; ok {
				slice[k] = v
				used[k] = true
			}
		}
		if err := m.Save(dest, slice); err != nil {
			return nil, fmt.Errorf("split: save %q: %w", dest, err)
		}
		results = append(results, Result{Name: dest, Keys: keys})
	}

	if remainder != "" {
		rem := make(map[string]string)
		for k, v := range vars {
			if !used[k] {
				rem[k] = v
			}
		}
		if err := m.Save(remainder, rem); err != nil {
			return nil, fmt.Errorf("split: save remainder %q: %w", remainder, err)
		}
		var remKeys []string
		for k := range rem {
			remKeys = append(remKeys, k)
		}
		results = append(results, Result{Name: remainder, Keys: remKeys})
	}

	return results, nil
}

// Format returns a human-readable summary of split results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no snapshots created\n"
	}
	out := ""
	for _, r := range results {
		out += fmt.Sprintf("%s: %d key(s)\n", r.Name, len(r.Keys))
	}
	return out
}
