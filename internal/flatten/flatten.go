package flatten

import "fmt"

// Manager can load named snapshots by name.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Result holds the merged vars and metadata about the operation.
type Result struct {
	Vars     map[string]string
	Sources  []string
	Conflicts int
}

// Run merges the snapshots identified by names into dest.
// Later entries in names win on key conflicts.
func Run(m Manager, names []string, dest string) (Result, error) {
	if len(names) == 0 {
		return Result{}, fmt.Errorf("flatten: at least one source name required")
	}

	merged := make(map[string]string)
	seen := make(map[string]string) // key -> first source that set it
	conflicts := 0

	for _, name := range names {
		vars, err := m.Load(name)
		if err != nil {
			return Result{}, fmt.Errorf("flatten: load %q: %w", name, err)
		}
		for k, v := range vars {
			if _, exists := seen[k]; exists {
				conflicts++
			}
			seen[k] = name
			merged[k] = v
		}
	}

	if err := m.Save(dest, merged); err != nil {
		return Result{}, fmt.Errorf("flatten: save %q: %w", dest, err)
	}

	return Result{
		Vars:      merged,
		Sources:   names,
		Conflicts: conflicts,
	}, nil
}
