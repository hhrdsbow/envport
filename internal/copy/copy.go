package copy

import "fmt"

// Manager is the interface required to copy snapshots.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, env map[string]string) error
	List() ([]string, error)
}

// ErrNotFound is returned when the source snapshot does not exist.
type ErrNotFound struct {
	Name string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("snapshot %q not found", e.Name)
}

// ErrAlreadyExists is returned when the destination snapshot already exists.
type ErrAlreadyExists struct {
	Name string
}

func (e *ErrAlreadyExists) Error() string {
	return fmt.Sprintf("snapshot %q already exists", e.Name)
}

// Run copies a snapshot from src to dst.
// If overwrite is true, an existing dst snapshot will be replaced.
func Run(m Manager, src, dst string, overwrite bool) error {
	env, err := m.Load(src)
	if err != nil {
		return &ErrNotFound{Name: src}
	}

	if !overwrite {
		names, err := m.List()
		if err != nil {
			return fmt.Errorf("listing snapshots: %w", err)
		}
		for _, n := range names {
			if n == dst {
				return &ErrAlreadyExists{Name: dst}
			}
		}
	}

	if err := m.Save(dst, env); err != nil {
		return fmt.Errorf("saving snapshot %q: %w", dst, err)
	}
	return nil
}
