package rename

import "fmt"

// Manager defines the interface for snapshot storage needed by rename.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, env map[string]string) error
	Delete(name string) error
}

// ErrNotFound is returned when the source snapshot does not exist.
type ErrNotFound struct {
	Name string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("snapshot %q not found", e.Name)
}

// Run renames a snapshot from src to dst.
// It loads the src snapshot, saves it under dst, then deletes src.
func Run(m Manager, src, dst string) error {
	if src == dst {
		return fmt.Errorf("source and destination names are identical: %q", src)
	}

	env, err := m.Load(src)
	if err != nil {
		return &ErrNotFound{Name: src}
	}

	if err := m.Save(dst, env); err != nil {
		return fmt.Errorf("saving snapshot %q: %w", dst, err)
	}

	if err := m.Delete(src); err != nil {
		return fmt.Errorf("deleting old snapshot %q: %w", src, err)
	}

	return nil
}
