package copy

import (
	"fmt"

	"github.com/envport/envport/internal/profile"
)

// Manager defines the interface for profile storage needed by copy.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
	List() ([]string, error)
}

// Run copies a snapshot from src to dst, optionally overwriting.
func Run(mgr Manager, src, dst string, overwrite bool) error {
	if src == "" {
		return fmt.Errorf("source snapshot name must not be empty")
	}
	if dst == "" {
		return fmt.Errorf("destination snapshot name must not be empty")
	}
	if src == dst {
		return fmt.Errorf("source and destination must differ")
	}

	vars, err := mgr.Load(src)
	if err != nil {
		return fmt.Errorf("load %q: %w", src, err)
	}

	if !overwrite {
		names, err := mgr.List()
		if err != nil {
			return fmt.Errorf("list snapshots: %w", err)
		}
		for _, n := range names {
			if n == dst {
				return fmt.Errorf("snapshot %q already exists; use --overwrite to replace", dst)
			}
		}
	}

	cloned := make(map[string]string, len(vars))
	for k, v := range vars {
		cloned[k] = v
	}

	if err := mgr.Save(dst, cloned); err != nil {
		return fmt.Errorf("save %q: %w", dst, err)
	}
	return nil
}

// Ensure *profile.Manager satisfies Manager at compile time.
var _ Manager = (*profile.Manager)(nil)
