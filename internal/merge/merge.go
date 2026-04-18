package merge

import "fmt"

// Manager can load and save snapshots by name.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Strategy controls how conflicts are resolved.
type Strategy int

const (
	StrategyBase  Strategy = iota // base wins on conflict
	StrategyOther                 // other wins on conflict
	StrategyError                 // error on conflict
)

// Result holds the merged vars and any conflicts detected.
type Result struct {
	Vars      map[string]string
	Conflicts []string
}

// Run merges snapshot `other` into snapshot `base`, writing the result to
// `dest`. Strategy controls conflict resolution.
func Run(m Manager, base, other, dest string, strategy Strategy) (*Result, error) {
	bVars, err := m.Load(base)
	if err != nil {
		return nil, fmt.Errorf("load base %q: %w", base, err)
	}
	oVars, err := m.Load(other)
	if err != nil {
		return nil, fmt.Errorf("load other %q: %w", other, err)
	}

	merged := make(map[string]string, len(bVars))
	for k, v := range bVars {
		merged[k] = v
	}

	var conflicts []string
	for k, v := range oVars {
		existing, exists := merged[k]
		if exists && existing != v {
			conflicts = append(conflicts, k)
			switch strategy {
			case StrategyError:
				return nil, fmt.Errorf("conflict on key %q", k)
			case StrategyOther:
				merged[k] = v
			// StrategyBase: keep existing
			}
		} else {
			merged[k] = v
		}
	}

	if err := m.Save(dest, merged); err != nil {
		return nil, fmt.Errorf("save dest %q: %w", dest, err)
	}
	return &Result{Vars: merged, Conflicts: conflicts}, nil
}
