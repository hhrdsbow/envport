package lowercase

import (
	"fmt"
	"strings"
)

// Manager is the interface for loading and saving snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the outcome of a lowercase operation.
type Result struct {
	Name    string
	Changed []string
}

// Run lowercases values for the given snapshot. If keys is non-empty, only
// those keys are affected. The modified snapshot is saved back via mgr.
func Run(mgr Manager, name string, keys []string) (Result, error) {
	vars, err := mgr.Load(name)
	if err != nil {
		return Result{}, fmt.Errorf("lowercase: load %q: %w", name, err)
	}

	filter := toSet(keys)
	var changed []string

	for k, v := range vars {
		if len(filter) > 0 && !filter[k] {
			continue
		}
		lc := strings.ToLower(v)
		if lc != v {
			vars[k] = lc
			changed = append(changed, k)
		}
	}

	if err := mgr.Save(name, vars); err != nil {
		return Result{}, fmt.Errorf("lowercase: save %q: %w", name, err)
	}

	return Result{Name: name, Changed: changed}, nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Changed) == 0 {
		return fmt.Sprintf("snapshot %q: no values changed", r.Name)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "snapshot %q: lowercased %d value(s)\n", r.Name, len(r.Changed))
	for _, k := range r.Changed {
		fmt.Fprintf(&sb, "  %s\n", k)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
