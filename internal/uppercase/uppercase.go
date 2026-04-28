// Package uppercase provides functionality to convert environment variable
// values to uppercase, optionally filtered to a specific set of keys.
package uppercase

import (
	"fmt"
	"strings"
)

// Manager describes the snapshot operations required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of an uppercase operation.
type Result struct {
	Name    string
	Changed []string
}

// Run loads the named snapshot, converts values to uppercase for the given
// keys (all keys when keys is empty), saves the result, and returns a
// summary of which keys were changed.
func Run(m Manager, name string, keys []string) (Result, error) {
	vars, err := m.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("uppercase: load %q: %w", name, err)
	}

	keySet := toSet(keys)
	updated := make(map[string]string, len(vars))
	var changed []string

	for k, v := range vars {
		if len(keySet) == 0 || keySet[k] {
			upper := strings.ToUpper(v)
			if upper != v {
				changed = append(changed, k)
			}
			updated[k] = upper
		} else {
			updated[k] = v
		}
	}

	if err := m.Save(name, updated); err != nil {
		return Result{}, fmt.Errorf("uppercase: save %q: %w", name, err)
	}

	return Result{Name: name, Changed: changed}, nil
}

// Format returns a human-readable summary of the Result.
func Format(r Result) string {
	if len(r.Changed) == 0 {
		return fmt.Sprintf("uppercase: %s — no values changed", r.Name)
	}
	return fmt.Sprintf("uppercase: %s — %d value(s) uppercased: %s",
		r.Name, len(r.Changed), strings.Join(r.Changed, ", "))
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
